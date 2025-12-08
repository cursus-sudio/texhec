package state

import (
	"engine/modules/uuid"
)

// this is a component of any type
// nil means that component don't exists.
// otherwise it should be of fixed type
type ComponentState any

type EntitySnapshot struct {
	// nil means that entity do not exists.
	// otherwise size of array is fixed and bound to component types synchronized.
	Components []ComponentState
}

type State struct {
	Entities map[uuid.UUID]EntitySnapshot
}

func (c1 State) MergeC1OverC2(c2 State) {
	for uuid, snapshot := range c2.Entities {
		if _, ok := c1.Entities[uuid]; !ok {
			c1.Entities[uuid] = snapshot
		}
	}
}

func (c1 State) MergeC2OverC1(c2 State) {
	for uuid, snapshot := range c2.Entities {
		c1.Entities[uuid] = snapshot
	}
}
