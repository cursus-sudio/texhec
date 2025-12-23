package netsync

import (
	"engine/modules/connection"
	"engine/modules/record"
	"engine/modules/uuid"
	"engine/services/ecs"
)

type ToolFactory ecs.ToolFactory[World, NetSyncTool]
type NetSyncTool interface {
	NetSync() Interface
}
type World interface {
	ecs.World
	uuid.UUIDTool
	connection.ConnectionTool
	record.RecordTool
}
type Interface interface {
	Server() ecs.ComponentsArray[ServerComponent]
	Client() ecs.ComponentsArray[ClientComponent]
}

type StartSystem ecs.SystemRegister[World]
type StopSystem ecs.SystemRegister[World]

// entity with this component and with connection component will be one with which we'll synchronize
type ServerComponent struct{}

// entity with this component and connection will get notifications about changes
type ClientComponent struct {
	// TODO permissions
}

// system stores:
// - versions changes (event id and may loop)
// - predicted events
// on any change:
// - store in system all adds, changes and removes to a system with note of version
// on any event (local):
// - change reconciliation version
// - store all changes
// on any event (server):
// - if matches with predicted event we're good
// - if doesn't match we revert latest changes and push event before predicted event
// - if we have to many predicted events than we remove them all

// event pointer should implement it
type AuthorizedEvent interface {
	SetConnection(ecs.EntityID)
}
