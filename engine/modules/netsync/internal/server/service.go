package server

import (
	"engine/modules/connection"
	"engine/modules/netsync"
	"engine/modules/netsync/internal/clienttypes"
	"engine/modules/netsync/internal/config"
	"engine/modules/netsync/internal/servertypes"
	"engine/modules/record"
	"engine/modules/uuid"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"
	"fmt"
	"reflect"
	"sync"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type clientMessage struct {
	Client  ecs.EntityID
	Message any
}

type Service struct {
	config.Config

	EventsBuilder events.Builder     `inject:"1"`
	Events        events.Events      `inject:"1"`
	World         ecs.World          `inject:"1"`
	NetSync       netsync.Service    `inject:"1"`
	Connection    connection.Service `inject:"1"`
	Record        record.Service     `inject:"1"`
	UUID          uuid.Service       `inject:"1"`

	Logger logger.Logger `inject:"1"`

	recordedEventUUID *uuid.UUID
	recordingID       record.UUIDRecordingID

	mutex *sync.Mutex

	dirtySet               ecs.DirtySet
	messagesSentFromClient []clientMessage
	toRemove               []ecs.EntityID
	listeners              datastructures.SparseSet[ecs.EntityID]
}

func NewService(c ioc.Dic, config config.Config) *Service {
	t := ioc.GetServices[*Service](c)
	t.Config = config
	t.recordedEventUUID = nil
	t.recordingID = 0

	t.mutex = &sync.Mutex{}

	t.dirtySet = ecs.NewDirtySet()
	t.messagesSentFromClient = nil
	t.toRemove = nil
	t.listeners = datastructures.NewSparseSet[ecs.EntityID]()

	// listen to server messages
	listeners := map[reflect.Type]func(ecs.EntityID, any){
		reflect.TypeFor[clienttypes.FetchStateDTO](): func(entity ecs.EntityID, a any) {
			t.ListenFetchState(entity, a.(clienttypes.FetchStateDTO))
		},
		reflect.TypeFor[clienttypes.EmitEventDTO](): func(entity ecs.EntityID, a any) {
			t.ListenEmitEvent(entity, a.(clienttypes.EmitEventDTO))
		},
		reflect.TypeFor[clienttypes.TransparentEventDTO](): func(entity ecs.EntityID, a any) {
			t.ListenTransparentEvent(entity, a.(clienttypes.TransparentEventDTO))
		},
	}
	events.Listen(t.EventsBuilder, func(frames.FrameEvent) {
		t.loadConnections()
		for len(t.toRemove) != 0 {
			entity := t.toRemove[0]
			t.mutex.Lock()
			t.toRemove = t.toRemove[1:]
			t.mutex.Unlock()

			t.listeners.Remove(entity)
			t.World.RemoveEntity(entity)
		}

		for len(t.messagesSentFromClient) != 0 {
			message := t.messagesSentFromClient[0]
			t.mutex.Lock()
			t.messagesSentFromClient = t.messagesSentFromClient[1:]
			t.mutex.Unlock()

			messageType := reflect.TypeOf(message.Message)
			listener, ok := listeners[messageType]
			if !ok {
				t.Logger.Warn(fmt.Errorf("invalid listener called there is no %v type", messageType.String()))
				continue
			}
			listener(message.Client, message.Message)
		}
	})
	t.NetSync.Client().AddDirtySet(t.dirtySet)
	t.Connection.Component().AddDirtySet(t.dirtySet)

	return t
}

// public methods

func (t *Service) BeforeEvent(event any) {
	t.loadConnections()
	if len(t.NetSync.Client().GetEntities()) == 0 {
		return
	}

	if t.recordedEventUUID == nil {
		uuid := t.UUID.NewUUID()
		t.recordedEventUUID = &uuid
	}
	t.recordingID = t.Record.UUID().StartRecording(t.RecordConfig)
}

func (t *Service) AfterEvent(event any) {
	t.loadConnections()
	if len(t.NetSync.Client().GetEntities()) == 0 {
		return
	}

	if recording, ok := t.Record.UUID().Stop(t.recordingID); ok && t.recordedEventUUID != nil {
		t.emitChanges(*t.recordedEventUUID, recording)
	}
	t.recordingID = 0
}

func (t *Service) OnTransparentEvent(event any) {
	if len(t.NetSync.Client().GetEntities()) == 0 {
		return
	}

	for _, client := range t.NetSync.Client().GetEntities() {
		connComp, ok := t.Connection.Component().Get(client)
		if !ok {
			return
		}
		t.Logger.Warn(connComp.Conn().Send(servertypes.TransparentEventDTO{Event: event}))
	}
}

func (t *Service) ListenFetchState(entity ecs.EntityID, dto clienttypes.FetchStateDTO) {
	state := t.Record.UUID().GetState(t.RecordConfig)
	t.sendVisible(entity, nil, state)
}

func (t *Service) ListenEmitEvent(entity ecs.EntityID, dto clienttypes.EmitEventDTO) {
	conn, ok := t.Connection.Component().Get(entity)
	if !ok {
		return
	}
	event, err := t.Auth(entity, dto.Event)
	if err != nil {
		err := conn.Conn().Send(servertypes.SendChangeDTO{Error: err})
		t.Logger.Warn(err)
		return
	}
	t.recordedEventUUID = &dto.ID
	events.EmitAny(t.Events, event)
}

func (t *Service) ListenTransparentEvent(entity ecs.EntityID, dto clienttypes.TransparentEventDTO) {
	conn, ok := t.Connection.Component().Get(entity)
	if !ok {
		return
	}
	event, err := t.Auth(entity, dto.Event)
	if err != nil {
		err := conn.Conn().Send(servertypes.TransparentEventDTO{Error: err})
		t.Logger.Warn(err)
		return
	}
	events.EmitAny(t.Events, event)
}

// private methods

func (t *Service) loadConnections() {
	for _, entity := range t.dirtySet.Get() {
		if ok := t.listeners.Get(entity); ok {
			continue
		}
		t.listeners.Add(entity)
		if _, ok := t.NetSync.Client().Get(entity); !ok {
			continue
		}

		comp, ok := t.Connection.Component().Get(entity)
		if !ok {
			continue
		}
		messages := comp.Conn().Messages()
		go func(entity ecs.EntityID) {
			for {
				message, ok := <-messages
				if !ok {
					break
				}
				t.mutex.Lock()
				t.messagesSentFromClient = append(t.messagesSentFromClient, clientMessage{
					Client:  entity,
					Message: message,
				})
				t.mutex.Unlock()
			}
			t.mutex.Lock()
			t.toRemove = append(t.toRemove, entity)
			t.mutex.Unlock()
		}(entity)
	}
}

func (t *Service) sendVisible(client ecs.EntityID, eventUUID *uuid.UUID, changes record.UUIDRecording) {
	connComp, ok := t.Connection.Component().Get(client)
	if !ok {
		return
	}

	// TODO manage visibility
	sentChanges := changes
	// for uuid, _ := range changes.Entities {
	// 	// if cannot use remove it
	// 	delete(changes.Entities, uuid)
	// }

	if len(sentChanges.Entities) == 0 {
		return
	}

	go func() {
		if eventUUID != nil {
			err := connComp.Conn().Send(servertypes.SendChangeDTO{
				EventID: *eventUUID,
				Changes: sentChanges,
			})
			if err != nil {
				t.mutex.Lock()
				t.toRemove = append(t.toRemove, client)
				t.mutex.Unlock()
			}
			// t.logger.Warn(err)
		} else {
			err := connComp.Conn().Send(servertypes.SendStateDTO{
				State: sentChanges,
			})
			if err != nil {
				t.mutex.Lock()
				t.toRemove = append(t.toRemove, client)
				t.mutex.Unlock()
			}
			// t.logger.Warn(err)
		}
	}()
}

func (t *Service) emitChanges(eventUUID uuid.UUID, changes record.UUIDRecording) {
	for _, client := range t.NetSync.Client().GetEntities() {
		t.sendVisible(client, &eventUUID, changes)
	}
}
