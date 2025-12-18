package client

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
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/ogiusek/events"
)

type savedPrediction struct {
	PredictedEvent clienttypes.PredictedEvent
	Snapshot       state.State
}

type recordedPrediction struct {
	PredictedEvent clienttypes.PredictedEvent
}

type toolState struct {
	recordNextEvent    bool
	predictions        []savedPrediction
	recordedPrediction *recordedPrediction
	sentTransparentEvent,
	receivedTransparentEvent bool

	mutex *sync.Mutex

	dirtySet           ecs.DirtySet
	messagesFromServer []any
	toRemove           []ecs.EntityID
	listeners          datastructures.SparseSet[ecs.EntityID]

	world           netsync.World
	serverArray     ecs.ComponentsArray[netsync.ServerComponent]
	connectionArray ecs.ComponentsArray[connection.ConnectionComponent]
	uuidArray       ecs.ComponentsArray[uuid.Component]
	stateTool       state.Tool
	logger          logger.Logger
}

// can:
// - apply server event
// - apply predicted event (starts and ends prediction)
type Tool struct {
	config.Config
	*toolState
}

func NewTool(
	config config.Config,
	stateToolFactory ecs.ToolFactory[netsync.World, state.Tool],
	logger logger.Logger,
	world netsync.World,
) Tool {
	t := Tool{
		config,
		&toolState{
			true,
			make([]savedPrediction, 0),
			nil,
			false,
			false,

			&sync.Mutex{},

			ecs.NewDirtySet(),
			nil,
			nil,
			datastructures.NewSparseSet[ecs.EntityID](),

			world,
			ecs.GetComponentsArray[netsync.ServerComponent](world),
			ecs.GetComponentsArray[connection.ConnectionComponent](world),
			ecs.GetComponentsArray[uuid.Component](world),
			stateToolFactory.Build(world),
			logger,
		},
	}

	listeners := map[reflect.Type]func(any){
		reflect.TypeFor[servertypes.SendStateDTO](): func(a any) {
			t.ListenSendState(a.(servertypes.SendStateDTO))
		},
		reflect.TypeFor[servertypes.SendChangeDTO](): func(a any) {
			t.ListenSendChange(a.(servertypes.SendChangeDTO))
		},
		reflect.TypeFor[servertypes.TransparentEventDTO](): func(a any) {
			t.ListenTransparentEvent(a.(servertypes.TransparentEventDTO))
		},
	}
	events.Listen(t.world.EventsBuilder(), func(frames.FrameEvent) {
		t.loadConnections()
		conn := t.getConnection()
		if conn == nil {
			return
		}
		t.mutex.Lock()
		defer t.mutex.Unlock()

		for len(t.toRemove) != 0 {
			entity := t.toRemove[0]
			t.toRemove = t.toRemove[1:]
			t.world.RemoveEntity(entity)
			t.listeners.Remove(entity)
		}
		for len(t.messagesFromServer) != 0 {
			message := t.messagesFromServer[0]
			t.messagesFromServer = t.messagesFromServer[1:]

			messageType := reflect.TypeOf(message)
			listener, ok := listeners[messageType]
			if !ok {
				t.logger.Warn(fmt.Errorf("invalid listener of type '%v' called", messageType.String()))
				conn.Close()
				return
			}
			listener(message)
		}
	})

	t.serverArray.AddDirtySet(t.dirtySet)
	t.connectionArray.AddDirtySet(t.dirtySet)

	// listen to entities changes
	for _, arrayCtor := range config.ArraysOfComponents {
		array := arrayCtor(world)
		array.BeforeGet(t.stateTool.RecordEntitiesChange)
	}

	return t
}

// public methods

// doesn't send event to server
func (t Tool) BeforeEventRecord(event any) {
	t.loadConnections()
	clientConn := t.getConnection()
	if clientConn == nil {
		return
	}
	if !t.recordNextEvent {
		t.recordNextEvent = true
		return
	}

	if len(t.predictions) > t.MaxPredictions {
		t.logger.Warn(ErrExceededPredictions)
		t.undoPredictions()
		// reconciliate
		if err := clientConn.Send(clienttypes.FetchStateDTO{}); err != nil {
			t.logger.Warn(err)
		}
		return
	}

	t.stateTool.StartRecording()
	t.recordedPrediction = &recordedPrediction{
		PredictedEvent: clienttypes.PredictedEvent{
			ID:    t.world.UUID().NewUUID(),
			Event: event,
		},
	}
}

func (t Tool) BeforeEvent(event any) {
	t.loadConnections()
	clientConn := t.getConnection()
	if clientConn == nil {
		return
	}
	t.BeforeEventRecord(event)
	if t.recordedPrediction == nil {
		return
	}

	dto := clienttypes.EmitEventDTO(t.recordedPrediction.PredictedEvent)
	if err := clientConn.Send(dto); err != nil {
		t.logger.Warn(err)
	}
}

func (t Tool) AfterEvent(event any) {
	conn := t.getConnection()
	if conn == nil {
		return
	}

	if t.recordedPrediction == nil {
		return
	}

	changes := t.stateTool.FinishRecording()
	newPrediction := savedPrediction{
		PredictedEvent: t.recordedPrediction.PredictedEvent,
		Snapshot:       *changes,
	}
	t.recordedPrediction = nil

	t.predictions = append(t.predictions, newPrediction)
}

func (t Tool) OnTransparentEvent(event any) {
	if t.receivedTransparentEvent {
		t.receivedTransparentEvent = false
		return
	}
	conn := t.getConnection()
	if conn == nil {
		return
	}

	t.sentTransparentEvent = true
	conn.Send(clienttypes.TransparentEventDTO{Event: event})
}

func (t Tool) ListenSendChange(dto servertypes.SendChangeDTO) {
	conn := t.getConnection()
	if conn == nil {
		return
	}
	if dto.Error != nil {
		predictedEvents := t.undoPredictions()
		// reApplied events are events without applied event
		reEmitedEvents := make([]clienttypes.PredictedEvent, 0, len(predictedEvents))
		for _, predictedEvent := range predictedEvents {
			if predictedEvent.ID != dto.EventID {
				reEmitedEvents = append(reEmitedEvents, predictedEvent)
			}
		}
		t.applyPredictedEvents(reEmitedEvents)
		t.logger.Warn(dto.Error)
		return
	}
	// check is event predicted. if is then remove first event from queue
	// if isn't then undo predictions, emit server event(as not recordable), emit all predicted events again
	if len(t.predictions) == 0 {
		t.stateTool.ApplyState(dto.Changes)
		return
	}
	if t.predictions[0].PredictedEvent.ID == dto.EventID {
		t.predictions = t.predictions[1:]
		return
		// TODO later. add test is prediction correct
		// replace this with state comparer for first prediction
		// stateEqual := true
		// if !stateEqual {
		// 	t.logger.Warn(ErrInvalidPrediction)
		// 	predictedEvents := t.UndoPredictions()
		// 	t.ApplyState(dto.Changes)
		// 	t.ApplyPredictedEvents(predictedEvents[1:])
		// } else {
		// 	t.predictions = t.predictions[1:]
		// }
	}
	predictedEvents := t.undoPredictions()
	t.stateTool.ApplyState(dto.Changes)
	// reApplied events are events without applied event
	reEmitedEvents := make([]clienttypes.PredictedEvent, 0, len(predictedEvents))
	for _, predictedEvent := range predictedEvents {
		if predictedEvent.ID != dto.EventID {
			reEmitedEvents = append(reEmitedEvents, predictedEvent)
		}
	}
	t.applyPredictedEvents(reEmitedEvents)
}

// reconciliate
func (t Tool) ListenSendState(dto servertypes.SendStateDTO) {
	conn := t.getConnection()
	if conn == nil {
		return
	}
	if dto.Error != nil {
		t.predictions = nil
		t.logger.Warn(dto.Error)
		conn.Close()
		return
	}
	t.predictions = nil
	t.stateTool.ApplyState(dto.State)
}

func (t Tool) ListenTransparentEvent(dto servertypes.TransparentEventDTO) {
	if t.sentTransparentEvent {
		t.sentTransparentEvent = false
		return
	}
	if dto.Error != nil {
		t.logger.Warn(dto.Error)
		return
	}
	t.receivedTransparentEvent = true
	events.EmitAny(t.world.Events(), dto.Event)
}

// private methods

func (t Tool) loadConnections() {
	ei := t.dirtySet.Get()
	if len(ei) == 0 {
		return
	}
	if len(ei) != 1 {
		t.logger.Warn(errors.New("has more than one server"))
		return
	}
	entity := ei[0]
	if ok := t.listeners.Get(entity); ok {
		return
	}
	if _, ok := t.serverArray.Get(entity); !ok {
		return
	}
	t.listeners.Add(entity)
	comp, ok := t.connectionArray.Get(entity)
	if !ok {
		return
	}
	conn := comp.Conn()
	messages := conn.Messages()
	if err := conn.Send(clienttypes.FetchStateDTO{}); err != nil {
		t.logger.Warn(err)
		return
	}
	go func(entity ecs.EntityID) {
		for {
			message, ok := <-messages
			if !ok {
				break
			}
			t.mutex.Lock()
			t.messagesFromServer = append(t.messagesFromServer, message)
			t.mutex.Unlock()
		}
		t.mutex.Lock()
		t.toRemove = append(t.toRemove, entity)
		t.mutex.Unlock()
	}(entity)
}

func (t Tool) undoPredictions() []clienttypes.PredictedEvent {
	// add events to the list
	var unDoneEvents []clienttypes.PredictedEvent
	for _, prediction := range t.predictions {
		unDoneEvents = append(unDoneEvents, prediction.PredictedEvent)
	}
	original := state.State{
		Entities: make(map[uuid.UUID]state.EntitySnapshot),
	}
	for _, prediction := range t.predictions {
		original.MergeC1OverC2(prediction.Snapshot)
	}
	t.stateTool.ApplyState(original)
	t.predictions = nil
	return unDoneEvents
}

func (t Tool) applyPredictedEvents(predictedEvents []clienttypes.PredictedEvent) {
	for _, predictedEvent := range predictedEvents[1:] {
		t.recordNextEvent = false
		events.EmitAny(t.world.Events(), predictedEvent.Event)
	}
}

func (t Tool) getConnection() connection.Conn {
	var conn connection.Conn
	if entities := t.serverArray.GetEntities(); len(entities) == 1 {
		server := entities[0]
		comp, ok := t.connectionArray.Get(server)
		if ok {
			conn = comp.Conn()
		}
	}
	if conn == nil { // isn't client clear all client data
		t.recordedPrediction = nil
		t.predictions = nil
	}
	return conn
}
