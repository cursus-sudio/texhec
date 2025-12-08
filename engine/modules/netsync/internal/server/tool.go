package server

import (
	"engine/modules/connection"
	"engine/modules/netsync"
	"engine/modules/netsync/internal/clienttypes"
	"engine/modules/netsync/internal/config"
	"engine/modules/netsync/internal/servertypes"
	"engine/modules/netsync/internal/state"
	"engine/modules/uuid"
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

	mutex                  *sync.Mutex
	messagesSentFromClient []clientMessage

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
			nil,

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
	t.world.Query().
		Require(netsync.ClientComponent{}).
		Require(connection.ConnectionComponent{}).
		Build().OnAdd(func(ei []ecs.EntityID) {
		for _, entity := range ei {
			comp, err := t.connectionArray.GetComponent(entity)
			if err != nil {
				t.logger.Warn(err)
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
				world.RemoveEntity(entity)
			}(entity)
		}
	})

	// listen to entities changes

	for _, arrayCtor := range config.ArraysOfComponents {
		array := arrayCtor(world)
		array.BeforeAdd(t.stateTool.RecordEntitiesChange)
		array.BeforeChange(t.stateTool.RecordEntitiesChange)
		array.BeforeRemove(t.stateTool.RecordEntitiesChange)
	}

	return t
}

// public methods

func (t Tool) BeforeEvent(event any) {
	// if there are no clients return

	if t.recordedEventUUID == nil {
		uuid := t.uniqueTool.NewUUID()
		t.recordedEventUUID = &uuid
	}
	t.stateTool.StartRecording()
}

func (t Tool) AfterEvent(event any) {
	// if there are no clients return

	if changes := t.stateTool.FinishRecording(); changes != nil && t.recordedEventUUID != nil {
		t.emitChanges(*t.recordedEventUUID, *changes)
	} else {
		t.logger.Warn(ErrRecordingDidntStartProperly)
	}
}

func (t Tool) OnTransparentEvent(event any) {
	// if there are no clients return

	for _, client := range t.clientArray.GetEntities() {
		connComp, err := t.connectionArray.GetComponent(client)
		if err != nil {
			t.logger.Warn(err)
			return
		}
		t.logger.Warn(connComp.Conn().Send(servertypes.TransparentEventDTO{Event: event}))
	}
}

func (t Tool) ListenFetchState(client ecs.EntityID, dto clienttypes.FetchStateDTO) {
	state := t.stateTool.GetState()
	t.sendVisible(client, nil, state)
}

func (t Tool) ListenEmitEvent(client ecs.EntityID, dto clienttypes.EmitEventDTO) {
	// TODO
	// is this event even present in config ?
	// can client do that ?
	// if yes than do that
	t.recordedEventUUID = &dto.ID
	events.EmitAny(t.world.Events(), dto.Event)
}

func (t Tool) ListenTransparentEvent(client ecs.EntityID, dto clienttypes.TransparentEventDTO) {
	events.EmitAny(t.world.Events(), dto.Event)
}

// private methods

func (t Tool) sendVisible(client ecs.EntityID, eventUUID *uuid.UUID, changes state.State) {
	connComp, err := t.connectionArray.GetComponent(client)
	if err != nil {
		t.logger.Warn(err)
		return
	}

	// TODO manage visibility
	sentChanges := changes
	// for uuid, _ := range changes.Entities {
	// 	// if cannot use remove it
	// 	delete(changes.Entities, uuid)
	// }

	if eventUUID != nil {
		// TODO make sending non-blocking
		err := connComp.Conn().Send(servertypes.SendChangeDTO{
			EventID: *eventUUID,
			Changes: sentChanges,
		})
		t.logger.Warn(err)
	} else {
		err := connComp.Conn().Send(servertypes.SendStateDTO{
			State: sentChanges,
		})
		t.logger.Warn(err)
	}
}

func (t Tool) emitChanges(eventUUID uuid.UUID, changes state.State) {
	for _, client := range t.clientArray.GetEntities() {
		t.sendVisible(client, &eventUUID, changes)
	}
}
