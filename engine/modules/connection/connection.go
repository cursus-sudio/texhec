package connection

import "engine/services/ecs"

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

type ConnectionComponent struct {
	conn Conn
}

func NewConnection(conn Conn) ConnectionComponent {
	return ConnectionComponent{conn}
}

func (comp ConnectionComponent) Conn() Conn {
	return comp.conn
}

// tool

type ConnectionTool interface {
	Connection() Interface
}

type Interface interface {
	Host(addr string, conn func(ConnectionComponent)) error
	Connect(addr string) (ConnectionComponent, error)
	MockConnectionPair() (c1, c2 ConnectionComponent)

	TransferConnection(fromEntity, toEntity ecs.EntityID) error
}
