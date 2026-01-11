package connection

import (
	"engine/services/ecs"
	"net"
)

// system

type System ecs.SystemRegister

// types

// singular connection interface
type Conn interface {
	// send has block behavior
	Send(message any) error

	// closed channel can be returned if connection is closed
	Messages() chan any
	Close() error
}

// components

type ListenerComponent struct {
	listener net.Listener
}

func NewListener(listener net.Listener) ListenerComponent {
	return ListenerComponent{listener}
}

func (comp *ListenerComponent) Listener() net.Listener {
	return comp.listener
}

//

type ConnectionComponent struct {
	conn Conn
}

func NewConnection(conn Conn) ConnectionComponent {
	return ConnectionComponent{conn}
}

func (comp *ConnectionComponent) Conn() Conn {
	return comp.conn
}

// tool

type Service interface {
	Component() ecs.ComponentsArray[ConnectionComponent]
	Listener() ecs.ComponentsArray[ListenerComponent]

	Host(addr string, conn func(ConnectionComponent)) (ListenerComponent, error)
	Connect(addr string) (ConnectionComponent, error)
	MockConnectionPair() (c1, c2 ConnectionComponent)

	TransferConnection(fromEntity, toEntity ecs.EntityID) error
}
