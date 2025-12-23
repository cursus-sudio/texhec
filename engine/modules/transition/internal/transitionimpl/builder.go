package transitionimpl

import (
	"engine/modules/transition"
	"engine/services/datastructures"
	"engine/services/ecs"
)

type Builder interface {
	Register(transition.System)
	Build() transition.System
}

type builder struct {
	systems        []transition.System
	easingFunction datastructures.SparseArray[transition.EasingID, transition.EasingFunction]
}

func NewBuilder() Builder {
	return &builder{
		systems:        nil,
		easingFunction: datastructures.NewSparseArray[transition.EasingID, transition.EasingFunction](),
	}
}

func (b *builder) Register(system transition.System) {
	b.systems = append(b.systems, system)
}
func (b *builder) RegisterEasingFunction(id transition.EasingID, fn transition.EasingFunction) {
	b.easingFunction.Set(id, fn)
}

//

func (b *builder) Build() transition.System {
	systems := b.systems
	return ecs.NewSystemRegister(func(w transition.World) error {
		for _, system := range systems {
			if err := system.Register(w); err != nil {
				return err
			}
		}
		return nil
	})
}
