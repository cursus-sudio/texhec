package backendapi

import "backend/src/modules/tacticalmap"

type Backend interface {
	TacticalMap() tacticalmap.TacticalMap
}

type backend struct {
	TacticalMapService tacticalmap.TacticalMap `inject:"1"`
}

func (backend backend) TacticalMap() tacticalmap.TacticalMap {
	return backend.TacticalMapService
}
