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
