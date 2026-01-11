package client

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
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type savedPrediction struct {
	PredictedEvent clienttypes.PredictedEvent
	Snapshot       record.UUIDRecording
}

type recordedPrediction struct {
	PredictedEvent clienttypes.PredictedEvent
}

// can:
// - apply server event
// - apply predicted event (starts and ends prediction)
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

	recordNextEvent    bool
	predictions        []savedPrediction
	recordedPrediction *recordedPrediction
	sentTransparentEvent,
	receivedTransparentEvent bool
	recordingID record.UUIDRecordingID

	mutex *sync.Mutex

	dirtySet           ecs.DirtySet
	messagesFromServer []any
	toRemove           []ecs.EntityID
	listeners          datastructures.SparseSet[ecs.EntityID]
}

func NewService(c ioc.Dic, config config.Config) *Service {
	t := ioc.GetServices[*Service](c)
	t.Config = config
	t.recordNextEvent = true
	t.predictions = make([]savedPrediction, 0)
	t.recordedPrediction = nil
	t.sentTransparentEvent = false
	t.receivedTransparentEvent = false
	t.recordingID = 0

	t.mutex = &sync.Mutex{}

	t.dirtySet = ecs.NewDirtySet()
	t.messagesFromServer = nil
	t.toRemove = nil
	t.listeners = datastructures.NewSparseSet[ecs.EntityID]()

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
	events.Listen(t.EventsBuilder, func(frames.FrameEvent) {
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
			t.World.RemoveEntity(entity)
			t.listeners.Remove(entity)
		}
		for len(t.messagesFromServer) != 0 {
			message := t.messagesFromServer[0]
			t.messagesFromServer = t.messagesFromServer[1:]

			messageType := reflect.TypeOf(message)
			listener, ok := listeners[messageType]
			if !ok {
				t.Logger.Warn(fmt.Errorf("invalid listener of type '%v' called", messageType.String()))
				_ = conn.Close()
				return
			}
			listener(message)
		}
	})

	t.NetSync.Server().AddDirtySet(t.dirtySet)
	t.Connection.Component().AddDirtySet(t.dirtySet)

	return t
}

// public methods

// doesn't send event to server
func (t *Service) BeforeEventRecord(event any) {
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
		t.Logger.Warn(ErrExceededPredictions)
		t.undoPredictions()
		// reconciliate
		if err := clientConn.Send(clienttypes.FetchStateDTO{}); err != nil {
			t.Logger.Warn(err)
		}
		return
	}

	t.recordingID = t.Record.UUID().StartBackwardsRecording(t.RecordConfig)
	t.recordedPrediction = &recordedPrediction{
		PredictedEvent: clienttypes.PredictedEvent{
			ID:    t.UUID.NewUUID(),
			Event: event,
		},
	}
}

func (t *Service) BeforeEvent(event any) {
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
		t.Logger.Warn(err)
	}
}

func (t *Service) AfterEvent(event any) {
	conn := t.getConnection()
	if conn == nil {
		return
	}

	if t.recordedPrediction == nil {
		return
	}

	recording, ok := t.Record.UUID().Stop(t.recordingID)
	t.recordingID = 0
	if !ok {
		return
	}
	newPrediction := savedPrediction{
		PredictedEvent: t.recordedPrediction.PredictedEvent,
		Snapshot:       recording,
	}
	t.recordedPrediction = nil

	t.predictions = append(t.predictions, newPrediction)
}

func (t *Service) OnTransparentEvent(event any) {
	if t.receivedTransparentEvent {
		t.receivedTransparentEvent = false
		return
	}
	conn := t.getConnection()
	if conn == nil {
		return
	}

	t.sentTransparentEvent = true
	err := conn.Send(clienttypes.TransparentEventDTO{Event: event})
	t.Logger.Warn(err)
}

func (t *Service) ListenSendChange(dto servertypes.SendChangeDTO) {
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
		t.Logger.Warn(dto.Error)
		return
	}
	// check is event predicted. if is then remove first event from queue
	// if isn't then undo predictions, emit server event(as not recordable), emit all predicted events again
	if len(t.predictions) == 0 {
		t.Record.UUID().Apply(t.RecordConfig, dto.Changes)
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
	t.Record.UUID().Apply(t.RecordConfig, dto.Changes)
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
func (t *Service) ListenSendState(dto servertypes.SendStateDTO) {
	conn := t.getConnection()
	if conn == nil {
		return
	}
	if dto.Error != nil {
		t.predictions = nil
		t.Logger.Warn(dto.Error)
		_ = conn.Close()
		return
	}
	t.predictions = nil
	t.Record.UUID().Apply(t.RecordConfig, dto.State)
}

func (t *Service) ListenTransparentEvent(dto servertypes.TransparentEventDTO) {
	if t.sentTransparentEvent {
		t.sentTransparentEvent = false
		return
	}
	if dto.Error != nil {
		t.Logger.Warn(dto.Error)
		return
	}
	t.receivedTransparentEvent = true
	events.EmitAny(t.Events, dto.Event)
}

// private methods

func (t *Service) loadConnections() {
	ei := t.dirtySet.Get()
	if len(ei) == 0 {
		return
	}
	if len(ei) != 1 {
		t.Logger.Warn(errors.New("has more than one server"))
		return
	}
	entity := ei[0]
	if ok := t.listeners.Get(entity); ok {
		return
	}
	if _, ok := t.NetSync.Server().Get(entity); !ok {
		return
	}
	t.listeners.Add(entity)
	comp, ok := t.Connection.Component().Get(entity)
	if !ok {
		return
	}
	conn := comp.Conn()
	messages := conn.Messages()
	if err := conn.Send(clienttypes.FetchStateDTO{}); err != nil {
		t.Logger.Warn(err)
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

func (t *Service) undoPredictions() []clienttypes.PredictedEvent {
	// add events to the list
	var unDoneEvents []clienttypes.PredictedEvent
	snapshots := make([]record.UUIDRecording, len(t.predictions))
	for _, prediction := range t.predictions {
		unDoneEvents = append(unDoneEvents, prediction.PredictedEvent)
		// snapshots = append([]record.UUIDRecording{prediction.Snapshot}, snapshots...)
		snapshots = append(snapshots, prediction.Snapshot)
	}
	t.Record.UUID().Apply(t.RecordConfig, snapshots...)
	t.predictions = nil
	return unDoneEvents
}

func (t *Service) applyPredictedEvents(predictedEvents []clienttypes.PredictedEvent) {
	for _, predictedEvent := range predictedEvents[1:] {
		t.recordNextEvent = false
		events.EmitAny(t.Events, predictedEvent.Event)
	}
}

func (t *Service) getConnection() connection.Conn {
	var conn connection.Conn
	if entities := t.NetSync.Server().GetEntities(); len(entities) == 1 {
		server := entities[0]
		comp, ok := t.Connection.Component().Get(server)
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
