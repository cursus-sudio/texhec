package mergedsystems

func GetSystem[Event any](system func(Event)) func(Event) error {
	return func(e Event) error { system(e); return nil }
}

type MergedSystems[Event any] interface {
	AddSystems(system ...func(Event) error) MergedSystems[Event]
	Listen(event Event)
}

type mergedSystems[Event any] struct {
	errorHandler func(error)
	systems      []func(Event) error
}

func NewMergedSystems[Event any](errorHandler func(error)) MergedSystems[Event] {
	return &mergedSystems[Event]{
		errorHandler: errorHandler,
		systems:      nil,
	}
}

func (s *mergedSystems[Event]) AddSystems(system ...func(Event) error) MergedSystems[Event] {
	s.systems = append(s.systems, system...)
	return s
}

func (s *mergedSystems[Event]) Listen(event Event) {
	errors := []error{}
	for _, system := range s.systems {
		if err := system(event); err != nil {
			errors = append(errors, err)
		}
	}
	for _, err := range errors {
		s.errorHandler(err)
	}
}
