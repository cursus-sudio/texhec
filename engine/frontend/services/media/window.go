package media

//

type WindowTitle string

type OnWindowTitleChangeEvent struct{ Before, Now WindowTitle }

//

type Resolution struct{ X, Y int }

type OnResolutionChangeEvent struct{ Before, Now Resolution }

//

type Display int

const (
	DisplayFullScreen Display = iota
	DisplayWindowed
)

type WindowApi interface {
	DrawApi() DrawApi

	Title() WindowTitle
	SetTitle(WindowTitle) error

	// TODO add icon

	Resolution() Resolution
	SetResolution(Resolution) error

	Display() Display
	SetDisplay(Display) error
}
