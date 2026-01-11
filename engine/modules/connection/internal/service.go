package internal

import (
	"encoding/binary"
	"engine/modules/connection"
	"engine/services/codec"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"
	"io"
	"net"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type service struct {
	EventsBuilder events.Builder `inject:"1"`

	World  ecs.World     `inject:"1"`
	Codec  codec.Codec   `inject:"1"`
	Logger logger.Logger `inject:"1"`

	listenersDirtySet ecs.DirtySet
	listeners         datastructures.Set[net.Listener]
	listenersArray    ecs.ComponentsArray[connection.ListenerComponent]

	connectionDirtySet ecs.DirtySet
	connections        datastructures.Set[connection.Conn]
	connectionArray    ecs.ComponentsArray[connection.ConnectionComponent]
}

func NewService(c ioc.Dic) connection.Service {
	t := ioc.GetServices[*service](c)
	t.listenersDirtySet = ecs.NewDirtySet()
	t.listeners = datastructures.NewSet[net.Listener]()
	t.listenersArray = ecs.GetComponentsArray[connection.ListenerComponent](t.World)

	t.connectionDirtySet = ecs.NewDirtySet()
	t.connections = datastructures.NewSet[connection.Conn]()
	t.connectionArray = ecs.GetComponentsArray[connection.ConnectionComponent](t.World)

	events.Listen(t.EventsBuilder, func(frames.FrameEvent) {
		t.BeforeConnectionGet()
	})

	t.listenersArray.AddDirtySet(t.listenersDirtySet)
	t.listenersArray.BeforeGet(t.BeforeListenerGet)

	t.connectionArray.AddDirtySet(t.connectionDirtySet)
	t.connectionArray.BeforeGet(t.BeforeConnectionGet)

	return t
}

func (t *service) BeforeListenerGet() {
	if entities := t.connectionDirtySet.Get(); len(entities) == 0 {
		return
	}
	present := datastructures.NewSet[net.Listener]()
	for _, entity := range t.listenersArray.GetEntities() {
		comp, ok := t.listenersArray.Get(entity)
		if !ok {
			continue
		}
		conn := comp.Listener()
		if conn == nil {
			continue
		}
		present.Add(conn)
	}

	for _, listener := range t.listeners.Get() {
		_, ok := present.GetIndex(listener)
		if ok {
			continue
		}
		t.listeners.RemoveElements(listener)
		_ = listener.Close()
	}
}

func (t *service) BeforeConnectionGet() {
	if entities := t.connectionDirtySet.Get(); len(entities) == 0 {
		return
	}
	present := datastructures.NewSet[connection.Conn]()
	for _, entity := range t.connectionArray.GetEntities() {
		comp, ok := t.connectionArray.Get(entity)
		if !ok {
			continue
		}
		conn := comp.Conn()
		if conn == nil {
			continue
		}
		present.Add(conn)
	}

	for _, connection := range t.connections.Get() {
		_, ok := present.GetIndex(connection)
		if ok {
			continue
		}
		t.connections.RemoveElements(connection)
		_ = connection.Close()
	}
}

func (t *service) Component() ecs.ComponentsArray[connection.ConnectionComponent] {
	return t.connectionArray
}
func (t *service) Listener() ecs.ComponentsArray[connection.ListenerComponent] {
	return t.listenersArray
}

func (t *service) Host(addr string, onConn func(connection.ConnectionComponent)) (connection.ListenerComponent, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return connection.ListenerComponent{}, err
	}
	t.listeners.Add(listener)
	go func() {
		for {
			rawConn, err := listener.Accept()
			if err != nil {
				break
			}
			conn := t.NewConnection(rawConn)
			t.connections.Add(conn)
			connComp := connection.NewConnection(conn)
			onConn(connComp)
		}
		_ = listener.Close()
	}()
	return connection.NewListener(listener), nil
}

func (t *service) Connect(addr string) (connection.ConnectionComponent, error) {
	rawConn, err := net.Dial("tcp", addr)
	if err != nil {
		return connection.ConnectionComponent{}, err
	}
	conn := t.NewConnection(rawConn)
	t.connections.Add(conn)
	connComp := connection.NewConnection(conn)
	return connComp, nil
}

func (t *service) MockConnectionPair() (connection.ConnectionComponent, connection.ConnectionComponent) {
	rawC1, rawC2 := net.Pipe()
	c1, c2 := t.NewConnection(rawC1), t.NewConnection(rawC2)
	t.connections.Add(c1)
	t.connections.Add(c2)
	comp1, comp2 := connection.NewConnection(c1), connection.NewConnection(c2)
	return comp1, comp2
}

func (t *service) TransferConnection(entityFrom, entityTo ecs.EntityID) error {
	comp, ok := t.connectionArray.Get(entityFrom)
	if !ok {
		return nil
	}
	t.connectionArray.Remove(entityFrom)
	t.connectionArray.Set(entityTo, comp)
	return nil
}

func (s *service) NewConnection(rawConn net.Conn) connection.Conn {
	messages := make(chan any)
	go func() {
		defer close(messages)

		for {
			messageLengthInBytes := make([]byte, 4)
			if _, err := io.ReadFull(rawConn, messageLengthInBytes); err != nil {
				break
			}
			messageLength := binary.BigEndian.Uint32(messageLengthInBytes)
			messageBytes := make([]byte, messageLength)
			if _, err := io.ReadFull(rawConn, messageBytes); err != nil {
				break
			}

			message, err := s.Codec.Decode(messageBytes)
			if err != nil {
				s.Logger.Warn(err)
				continue
			}
			// f.logger.Info(fmt.Sprintf("received '***' type '%v'", reflect.TypeOf(message).String()))

			messages <- message
		}

		_ = rawConn.Close()
	}()
	return &conn{
		service:  s,
		conn:     rawConn,
		messages: messages,
	}
}
