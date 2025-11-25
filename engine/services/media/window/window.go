package window

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

type MousePos struct{ X, Y int32 }

func NewMousePos(x, y int32) MousePos  { return MousePos{x, y} }
func (p *MousePos) Elem() (x, y int32) { return p.X, p.Y }

type Api interface {
	NormalizeMousePos(MousePos) mgl32.Vec2
	GetMousePos() MousePos
	Window() *sdl.Window
	Ctx() sdl.GLContext
}

type api struct {
	window  *sdl.Window
	context sdl.GLContext
}

func newApi(
	window *sdl.Window,
	context sdl.GLContext,
) Api {
	return api{
		window:  window,
		context: context,
	}
}

func (api api) NormalizeMousePos(mousePos MousePos) mgl32.Vec2 {
	x, y := mousePos.Elem()
	w, h := api.Window().GetSize()
	return mgl32.Vec2{
		(2*float32(x)/float32(w) - 1),
		-(2*float32(y)/float32(h) - 1),
	}
}
func (api api) GetMousePos() MousePos {
	x, y, _ := sdl.GetMouseState()
	return NewMousePos(x, y)
}
func (api api) Window() *sdl.Window { return api.window }
func (api api) Ctx() sdl.GLContext  { return nil }
