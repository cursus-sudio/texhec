package internal

import (
	"engine/modules/connection"
	"engine/services/codec"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/logger"
	"net"
	"sync"
)

type tool struct {
	*factory
	dirtySet        ecs.DirtySet
	connections     datastructures.Set[connection.Conn]
	connectionArray ecs.ComponentsArray[connection.ConnectionComponent]
}

func NewToolFactory(
	codec codec.Codec,
	logger logger.Logger,
) ecs.ToolFactory[connection.Connection] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ecs.World) connection.Connection {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[tool](w); ok {
			return t
		}
		t := tool{
			NewFactory(codec, logger),
			ecs.NewDirtySet(),
			datastructures.NewSet[connection.Conn](),
			ecs.GetComponentsArray[connection.ConnectionComponent](w),
		}
		w.SaveGlobal(t)

		t.connectionArray.AddDirtySet(t.dirtySet)
		t.connectionArray.BeforeGet(t.BeforeGet)

		return t
	})
}

func (t tool) BeforeGet() {
	t.dirtySet.Get()
	present := datastructures.NewSet[connection.Conn]()
	for _, entity := range t.connectionArray.GetEntities() {
		comp, ok := t.connectionArray.GetComponent(entity)
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
		connection.Close()
	}
}

func (t tool) Host(addr string, onConn func(connection.ConnectionComponent)) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	go func() {
		defer listener.Close()

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
	}()
	return nil
}

func (t tool) Connect(addr string) (connection.ConnectionComponent, error) {
	rawConn, err := net.Dial("tcp", addr)
	if err != nil {
		return connection.ConnectionComponent{}, err
	}
	conn := t.NewConnection(rawConn)
	t.connections.Add(conn)
	connComp := connection.NewConnection(conn)
	return connComp, nil
}

func (t tool) MockConnectionPair() (connection.ConnectionComponent, connection.ConnectionComponent) {
	rawC1, rawC2 := net.Pipe()
	c1, c2 := t.NewConnection(rawC1), t.NewConnection(rawC2)
	t.connections.Add(c1)
	t.connections.Add(c2)
	comp1, comp2 := connection.NewConnection(c1), connection.NewConnection(c2)
	return comp1, comp2
}

func (t tool) TransferConnection(entityFrom, entityTo ecs.EntityID) error {
	comp, ok := t.connectionArray.GetComponent(entityFrom)
	if !ok {
		return nil
	}
	t.connectionArray.RemoveComponent(entityFrom)
	t.connectionArray.SaveComponent(entityTo, comp)
	return nil
}
