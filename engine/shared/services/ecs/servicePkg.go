package ecs

type SystemRegister interface {
	Register(World) error
}

// impl

type systemRegister struct{ register func(World) error }

func (s systemRegister) Register(w World) error            { return s.register(w) }
func NewSystemRegister(l func(World) error) SystemRegister { return systemRegister{l} }

// helpers

func RegisterSystems(w World, systems ...SystemRegister) []error {
	errors := []error{}
	for _, system := range systems {
		if system == nil {
			continue
		}
		if err := system.Register(w); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

//

type ToolFactory[Tool any] interface {
	Build(World) Tool
}
type toolFactory[Tool any] struct{ build func(World) Tool }

func (f toolFactory[Tool]) Build(w World) Tool { return f.build(w) }
func NewToolFactory[Tool any](l func(World) Tool) ToolFactory[Tool] {
	return &toolFactory[Tool]{build: l}
}
