package transform

// this says that his object doesn't change
type Static struct{}

func NewStatic() Static {
	return Static{}
}
