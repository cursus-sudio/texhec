package worldprojections

import (
	"frontend/services/datastructures"
	"frontend/services/ecs"
)

type WorldProjectionsRegister struct {
	Projections datastructures.Set[ecs.ComponentType]
}

func (r WorldProjectionsRegister) Release() {}

func NewWorldProjectionsRegister(projections ...ecs.ComponentType) WorldProjectionsRegister {
	set := datastructures.NewSet[ecs.ComponentType]()
	set.Add(projections...)
	return WorldProjectionsRegister{set}
}
