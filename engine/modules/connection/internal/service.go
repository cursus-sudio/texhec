package internal

import (
	"engine/modules/connection"
	"engine/services/codec"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"
	"net"

	"github.com/ogiusek/events"
)

type service struct {
	*factory

	listenersDirtySet ecs.DirtySet
	listeners         datastructures.Set[net.Listener]
	listenersArray    ecs.ComponentsArray[connection.ListenerComponent]

	connectionDirtySet ecs.DirtySet
	connections        datastructures.Set[connection.Conn]
	connectionArray    ecs.ComponentsArray[connection.ConnectionComponent]
}

func NewToolFactory(
	codec codec.Codec,
	logger logger.Logger,
	world ecs.World,
	eventsBuilder events.Builder,
) connection.Service {
	t := &service{
		NewFactory(codec, logger),

		ecs.NewDirtySet(),
		datastructures.NewSet[net.Listener](),
		ecs.GetComponentsArray[connection.ListenerComponent](world),

		ecs.NewDirtySet(),
		datastructures.NewSet[connection.Conn](),
		ecs.GetComponentsArray[connection.ConnectionComponent](world),
	}
	events.Listen(eventsBuilder, func(frames.FrameEvent) {
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
