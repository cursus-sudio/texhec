package internal

import (
	"engine/modules/connection"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"
	"net"
	"sync"
)

type tool struct {
	*factory
	connectionArray ecs.ComponentsArray[connection.ConnectionComponent]
}

func NewToolFactory(
	codec codec.Codec,
	logger logger.Logger,
) ecs.ToolFactory[connection.Tool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ecs.World) connection.Tool {
		if t, err := ecs.GetGlobal[tool](w); err == nil {
			return t
		}
		mutex.Lock()
		defer mutex.Unlock()
		if t, err := ecs.GetGlobal[tool](w); err == nil {
			return t
		}
		t := tool{
			NewFactory(codec, logger),
			ecs.GetComponentsArray[connection.ConnectionComponent](w),
		}
		w.SaveGlobal(t)

		t.connectionArray.BeforeRemove(func(ei []ecs.EntityID) {
			for _, entity := range ei {
				comp, err := t.connectionArray.GetComponent(entity)
				if err != nil {
					continue
				}
				if comp.Conn() != nil {
					logger.Warn(comp.Conn().Close())
				}
			}
		})

		return t
	})
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
	connComp := connection.NewConnection(conn)
	return connComp, nil
}

func (t tool) TransferConnection(entityFrom, entityTo ecs.EntityID) error {
	comp, err := t.connectionArray.GetComponent(entityFrom)
	if err != nil {
		return nil
	}
	transaction := t.connectionArray.Transaction()
	transaction.SaveComponent(entityFrom, connection.ConnectionComponent{})
	if err := ecs.FlushMany(transaction); err != nil {
		return err
	}

	transaction.RemoveComponent(entityFrom)
	transaction.SaveComponent(entityTo, comp)
	if err := ecs.FlushMany(transaction); err != nil {
		return err
	}
	return nil
}
