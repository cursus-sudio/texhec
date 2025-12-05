package connection

import "engine/services/ecs"

// system

type System ecs.SystemRegister

// types

type Connection interface {
	// send has block behavior
	Send(message any) error

	// closed channel can be returned if connection is closed
	Messages() chan any
	Close() error
}

// components

type ConnectionComponent struct {
	conn Connection
}

func NewConnection(conn Connection) ConnectionComponent {
	return ConnectionComponent{conn}
}

func (comp ConnectionComponent) Conn() Connection {
	return comp.conn
}

// tool

type Tool interface {
	Host(addr string, conn func(ConnectionComponent)) error
	Connect(addr string) (ConnectionComponent, error)

	TransferConnection(fromEntity, toEntity ecs.EntityID) error
}
