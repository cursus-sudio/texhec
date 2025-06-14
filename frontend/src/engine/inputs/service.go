package inputs

type Key any

type MousePos struct{ X, Y int }

type Inputs interface {
	// keys
	OnPress(func(Key))
	OnUnPress(func(Key))

	// mouse (pos, move, click, unclick, maybe double click)
	OnClick(MousePos)
}
