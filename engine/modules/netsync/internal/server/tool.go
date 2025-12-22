package server

import (
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

	netsync.World
	netsync.NetSyncTool
	logger logger.Logger
}

type Tool struct {
	config.Config
	*toolState
}

func NewTool(
	config config.Config,
	netSyncToolFactory netsync.ToolFactory,
	logger logger.Logger,
	world netsync.World,
) Tool {
	t := Tool{
		config,
		&toolState{
			nil,
			0,

			&sync.Mutex{},

			ecs.NewDirtySet(),
			nil,
			nil,
			datastructures.NewSparseSet[ecs.EntityID](),

			world,
			netSyncToolFactory.Build(world),
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
	events.Listen(t.EventsBuilder(), func(frames.TickEvent) {
		t.loadConnections()
		for len(t.toRemove) != 0 {
			entity := t.toRemove[0]
			t.mutex.Lock()
			t.toRemove = t.toRemove[1:]
			t.mutex.Unlock()

			t.listeners.Remove(entity)
			t.RemoveEntity(entity)
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
	t.NetSync().Client().AddDirtySet(t.dirtySet)
	t.Connection().Component().AddDirtySet(t.dirtySet)

	return t
}

// public methods

func (t Tool) BeforeEvent(event any) {
	t.loadConnections()
	if len(t.NetSync().Client().GetEntities()) == 0 {
		return
	}

	if t.recordedEventUUID == nil {
		uuid := t.UUID().NewUUID()
		t.recordedEventUUID = &uuid
	}
	t.recordingID = t.Record().UUID().StartRecording(t.RecordConfig)
}

func (t Tool) AfterEvent(event any) {
	t.loadConnections()
	if len(t.NetSync().Client().GetEntities()) == 0 {
		return
	}

	if recording, ok := t.Record().UUID().Stop(t.recordingID); ok && t.recordedEventUUID != nil {
		t.emitChanges(*t.recordedEventUUID, recording)
	}
	t.recordingID = 0
}

func (t Tool) OnTransparentEvent(event any) {
	if len(t.NetSync().Client().GetEntities()) == 0 {
		return
	}

	for _, client := range t.NetSync().Client().GetEntities() {
		connComp, ok := t.Connection().Component().Get(client)
		if !ok {
			return
		}
		t.logger.Warn(connComp.Conn().Send(servertypes.TransparentEventDTO{Event: event}))
	}
}

func (t Tool) ListenFetchState(entity ecs.EntityID, dto clienttypes.FetchStateDTO) {
	state := t.Record().UUID().GetState(t.RecordConfig)
	t.sendVisible(entity, nil, state)
}

func (t Tool) ListenEmitEvent(entity ecs.EntityID, dto clienttypes.EmitEventDTO) {
	conn, ok := t.Connection().Component().Get(entity)
	if !ok {
		return
	}
	event, err := t.Config.Auth(entity, dto.Event)
	if err != nil {
		conn.Conn().Send(servertypes.SendChangeDTO{Error: err})
		t.logger.Warn(err)
		return
	}
	t.recordedEventUUID = &dto.ID
	events.EmitAny(t.Events(), event)
}

func (t Tool) ListenTransparentEvent(entity ecs.EntityID, dto clienttypes.TransparentEventDTO) {
	conn, ok := t.Connection().Component().Get(entity)
	if !ok {
		return
	}
	event, err := t.Config.Auth(entity, dto.Event)
	if err != nil {
		conn.Conn().Send(servertypes.TransparentEventDTO{Error: err})
		t.logger.Warn(err)
		return
	}
	events.EmitAny(t.Events(), event)
}

// private methods

func (t Tool) loadConnections() {
	for _, entity := range t.dirtySet.Get() {
		if ok := t.listeners.Get(entity); ok {
			continue
		}
		t.listeners.Add(entity)
		if _, ok := t.NetSync().Client().Get(entity); !ok {
			continue
		}

		comp, ok := t.Connection().Component().Get(entity)
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

func (t Tool) sendVisible(client ecs.EntityID, eventUUID *uuid.UUID, changes record.UUIDRecording) {
	connComp, ok := t.Connection().Component().Get(client)
	if !ok {
		return
	}

	// TODO manage visibility
	sentChanges := changes
	// for uuid, _ := range changes.Entities {
	// 	// if cannot use remove it
	// 	delete(changes.Entities, uuid)
	// }

	if len(sentChanges.UUIDEntities) == 0 {
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

func (t Tool) emitChanges(eventUUID uuid.UUID, changes record.UUIDRecording) {
	for _, client := range t.NetSync().Client().GetEntities() {
		t.sendVisible(client, &eventUUID, changes)
	}
}
