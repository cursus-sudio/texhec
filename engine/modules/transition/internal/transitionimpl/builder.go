package transitionimpl

import (
	"engine/modules/transition"
	"engine/services/ecs"
)

type Builder interface {
	Register(transition.System)
	Build() transition.System
}

type builder struct {
	systems []transition.System
}

func NewBuilder() Builder {
	return &builder{
		systems: nil,
	}
}

func (b *builder) Register(system transition.System) {
	b.systems = append(b.systems, system)
}

//

func (b *builder) Build() transition.System {
	systems := b.systems
	return ecs.NewSystemRegister(func() error {
		for _, system := range systems {
			if err := system.Register(); err != nil {
				return err
			}
		}
		return nil
	})
}
