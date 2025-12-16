package server

import (
	"engine/modules/connection"
	"engine/modules/netsync"
	"engine/modules/netsync/internal/clienttypes"
	"engine/modules/netsync/internal/config"
	"engine/modules/netsync/internal/servertypes"
	"engine/modules/netsync/internal/state"
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

	mutex *sync.Mutex

	dirtySet               ecs.DirtySet
	messagesSentFromClient []clientMessage
	toRemove               []ecs.EntityID
	listeners              datastructures.SparseSet[ecs.EntityID]

	world           ecs.World
	clientArray     ecs.ComponentsArray[netsync.ClientComponent]
	connectionArray ecs.ComponentsArray[connection.ConnectionComponent]
	uuidArray       ecs.ComponentsArray[uuid.Component]
	stateTool       state.Tool
	uniqueTool      uuid.Tool
	logger          logger.Logger
}

type Tool struct {
	config.Config
	*toolState
}

func NewTool(
	config config.Config,
	uniqueToolFactory ecs.ToolFactory[uuid.Tool],
	stateToolFactory ecs.ToolFactory[state.Tool],
	logger logger.Logger,
	world ecs.World,
) Tool {
	t := Tool{
		config,
		&toolState{
			nil,

			&sync.Mutex{},

			ecs.NewDirtySet(),
			nil,
			nil,
			datastructures.NewSparseSet[ecs.EntityID](),

			world,
			ecs.GetComponentsArray[netsync.ClientComponent](world),
			ecs.GetComponentsArray[connection.ConnectionComponent](world),
			ecs.GetComponentsArray[uuid.Component](world),
			stateToolFactory.Build(world),
			uniqueToolFactory.Build(world),
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
	events.Listen(t.world.EventsBuilder(), func(frames.FrameEvent) {
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
	t.clientArray.AddDirtySet(t.dirtySet)
	t.connectionArray.AddDirtySet(t.dirtySet)

	// listen to entities changes

	for _, arrayCtor := range config.ArraysOfComponents {
		array := arrayCtor(world)
		// before changes record
		array.BeforeGet(t.stateTool.RecordEntitiesChange)
	}

	return t
}

// public methods

func (t Tool) BeforeEvent(event any) {
	t.loadConnections()
	if len(t.clientArray.GetEntities()) == 0 {
		return
	}

	if t.recordedEventUUID == nil {
		uuid := t.uniqueTool.NewUUID()
		t.recordedEventUUID = &uuid
	}
	t.stateTool.StartRecording()
}

func (t Tool) AfterEvent(event any) {
	t.loadConnections()
	if len(t.clientArray.GetEntities()) == 0 {
		return
	}

	if changes := t.stateTool.FinishRecording(); changes != nil && t.recordedEventUUID != nil {
		t.emitChanges(*t.recordedEventUUID, *changes)
	} else {
		t.logger.Warn(ErrRecordingDidntStartProperly)
	}
}

func (t Tool) OnTransparentEvent(event any) {
	if len(t.clientArray.GetEntities()) == 0 {
		return
	}

	for _, client := range t.clientArray.GetEntities() {
		connComp, ok := t.connectionArray.GetComponent(client)
		if !ok {
			return
		}
		t.logger.Warn(connComp.Conn().Send(servertypes.TransparentEventDTO{Event: event}))
	}
}

func (t Tool) ListenFetchState(entity ecs.EntityID, dto clienttypes.FetchStateDTO) {
	state := t.stateTool.GetState()
	t.sendVisible(entity, nil, state)
}

func (t Tool) ListenEmitEvent(entity ecs.EntityID, dto clienttypes.EmitEventDTO) {
	conn, ok := t.connectionArray.GetComponent(entity)
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
	events.EmitAny(t.world.Events(), event)
}

func (t Tool) ListenTransparentEvent(entity ecs.EntityID, dto clienttypes.TransparentEventDTO) {
	conn, ok := t.connectionArray.GetComponent(entity)
	if !ok {
		return
	}
	event, err := t.Config.Auth(entity, dto.Event)
	if err != nil {
		conn.Conn().Send(servertypes.TransparentEventDTO{Error: err})
		t.logger.Warn(err)
		return
	}
	events.EmitAny(t.world.Events(), event)
}

// private methods

func (t Tool) loadConnections() {
	for _, entity := range t.dirtySet.Get() {
		if ok := t.listeners.Get(entity); ok {
			continue
		}
		t.listeners.Add(entity)
		if _, ok := t.clientArray.GetComponent(entity); !ok {
			continue
		}

		comp, ok := t.connectionArray.GetComponent(entity)
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

func (t Tool) sendVisible(client ecs.EntityID, eventUUID *uuid.UUID, changes state.State) {
	connComp, ok := t.connectionArray.GetComponent(client)
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

func (t Tool) emitChanges(eventUUID uuid.UUID, changes state.State) {
	for _, client := range t.clientArray.GetEntities() {
		t.sendVisible(client, &eventUUID, changes)
	}
}
