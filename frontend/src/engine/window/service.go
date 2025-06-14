package window

/*
objects can have:
- position
- zindex (also sorting group)
- other objects
- collision box
- texture / model
- children
*/

// type Object2D struct {
// 	X, Y   int
// 	ZIndex uint
// }

type Window interface {
	WindowWidth() uint
	WindowHeight() uint

	// icons
	// title
	// etc.
	// type Canvas interface {
	// draw (text, images, figures)
	// clean
	// }
}
