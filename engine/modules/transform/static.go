package transform

// this says that his object doesn't change
type StaticComponent struct{}

func NewStatic() StaticComponent {
	return StaticComponent{}
}
