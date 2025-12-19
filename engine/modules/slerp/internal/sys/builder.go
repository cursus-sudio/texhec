package sys

import (
	"engine/modules/slerp"
	"engine/services/ecs"
)

type Builder interface {
	Register(slerp.System)
	Build() slerp.System
}

type builder struct {
	systems []slerp.System
}

func NewBuilder() Builder {
	return &builder{}
}

func (b *builder) Register(system slerp.System) {
	b.systems = append(b.systems, system)
}

//

func (b *builder) Build() slerp.System {
	systems := b.systems
	return ecs.NewSystemRegister(func(w slerp.World) error {
		for _, system := range systems {
			system.Register(w)
		}
		return nil
	})
}
