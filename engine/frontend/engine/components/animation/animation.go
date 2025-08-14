package animation

// discard this

type Step float32

type SlerpedComponent[Component any] struct {
	Before, After Component
}

type Slerp[Component any] interface {
	Slerp(SlerpedComponent[Component], Step) Component
}

type Animation struct {
	Step         Step
	StepFunction func(Step) Step
	// components which are slerpable (SlerpedComponent[any])
	Components []any
}
