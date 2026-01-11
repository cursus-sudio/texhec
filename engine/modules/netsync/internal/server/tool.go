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
)

type clientMessage struct {
	Client  ecs.EntityID
	Message any
}

type toolState struct {
	recordedEventUUID *uuid.UUID
	recordingID       record.UUIDRecordingID

	mutex *sync.Mutex

	dirtySet               ecs.DirtySet
	messagesSentFromClient []clientMessage
	toRemove               []ecs.EntityID
	listeners              datastructures.SparseSet[ecs.EntityID]

	events     events.Events
	world      ecs.World
	netSync    netsync.Service
	connection connection.Service
	record     record.Service
	uuid       uuid.Service

	logger logger.Logger
}

type Service struct {
	config.Config
	*toolState
}

func NewService(
	config config.Config,
	logger logger.Logger,
	eventsBuilder events.Builder,
	world ecs.World,
	netSync netsync.Service,
	connection connection.Service,
	record record.Service,
	uuid uuid.Service,
) *Service {
	t := &Service{
		config,
		&toolState{
			nil,
			0,

			&sync.Mutex{},

			ecs.NewDirtySet(),
			nil,
			nil,
			datastructures.NewSparseSet[ecs.EntityID](),

			eventsBuilder.Events(),
			world,
			netSync,
			connection,
			record,
			uuid,

			logger,
		},
	}

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
	events.Listen(eventsBuilder, func(frames.FrameEvent) {
		t.loadConnections()
		for len(t.toRemove) != 0 {
			entity := t.toRemove[0]
			t.mutex.Lock()
			t.toRemove = t.toRemove[1:]
			t.mutex.Unlock()

			t.listeners.Remove(entity)
			t.world.RemoveEntity(entity)
		}

		for len(t.messagesSentFromClient) != 0 {
			message := t.messagesSentFromClient[0]
			t.mutex.Lock()
			t.messagesSentFromClient = t.messagesSentFromClient[1:]
			t.mutex.Unlock()

			messageType := reflect.TypeOf(message.Message)
			listener, ok := listeners[messageType]
			if !ok {
				t.logger.Warn(fmt.Errorf("invalid listener called there is no %v type", messageType.String()))
				continue
			}
			listener(message.Client, message.Message)
		}
	})
	t.netSync.Client().AddDirtySet(t.dirtySet)
	t.connection.Component().AddDirtySet(t.dirtySet)

	return t
}

// public methods

func (t *Service) BeforeEvent(event any) {
	t.loadConnections()
	if len(t.netSync.Client().GetEntities()) == 0 {
		return
	}

	if t.recordedEventUUID == nil {
		uuid := t.uuid.NewUUID()
		t.recordedEventUUID = &uuid
	}
	t.recordingID = t.record.UUID().StartRecording(t.RecordConfig)
}

func (t *Service) AfterEvent(event any) {
	t.loadConnections()
	if len(t.netSync.Client().GetEntities()) == 0 {
		return
	}

	if recording, ok := t.record.UUID().Stop(t.recordingID); ok && t.recordedEventUUID != nil {
		t.emitChanges(*t.recordedEventUUID, recording)
	}
	t.recordingID = 0
}

func (t *Service) OnTransparentEvent(event any) {
	if len(t.netSync.Client().GetEntities()) == 0 {
		return
	}

	for _, client := range t.netSync.Client().GetEntities() {
		connComp, ok := t.connection.Component().Get(client)
		if !ok {
			return
		}
		t.logger.Warn(connComp.Conn().Send(servertypes.TransparentEventDTO{Event: event}))
	}
}

func (t *Service) ListenFetchState(entity ecs.EntityID, dto clienttypes.FetchStateDTO) {
	state := t.record.UUID().GetState(t.RecordConfig)
	t.sendVisible(entity, nil, state)
}

func (t *Service) ListenEmitEvent(entity ecs.EntityID, dto clienttypes.EmitEventDTO) {
	conn, ok := t.connection.Component().Get(entity)
	if !ok {
		return
	}
	event, err := t.Auth(entity, dto.Event)
	if err != nil {
		err := conn.Conn().Send(servertypes.SendChangeDTO{Error: err})
		t.logger.Warn(err)
		return
	}
	t.recordedEventUUID = &dto.ID
	events.EmitAny(t.events, event)
}

func (t *Service) ListenTransparentEvent(entity ecs.EntityID, dto clienttypes.TransparentEventDTO) {
	conn, ok := t.connection.Component().Get(entity)
	if !ok {
		return
	}
	event, err := t.Auth(entity, dto.Event)
	if err != nil {
		err := conn.Conn().Send(servertypes.TransparentEventDTO{Error: err})
		t.logger.Warn(err)
		return
	}
	events.EmitAny(t.events, event)
}

// private methods

func (t *Service) loadConnections() {
	for _, entity := range t.dirtySet.Get() {
		if ok := t.listeners.Get(entity); ok {
			continue
		}
		t.listeners.Add(entity)
		if _, ok := t.netSync.Client().Get(entity); !ok {
			continue
		}

		comp, ok := t.connection.Component().Get(entity)
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
	connComp, ok := t.connection.Component().Get(client)
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
	for _, client := range t.netSync.Client().GetEntities() {
		t.sendVisible(client, &eventUUID, changes)
	}
}
