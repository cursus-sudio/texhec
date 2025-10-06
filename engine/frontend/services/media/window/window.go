package window

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

type Api interface {
	NormalizeMousePos(x, y int) mgl32.Vec2
	GetMousePos() (x, y int)
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

func (api api) NormalizeMousePos(x, y int) mgl32.Vec2 {
	w, h := api.Window().GetSize()
	return mgl32.Vec2{
		(2*float32(x)/float32(w) - 1),
		-(2*float32(y)/float32(h) - 1),
	}
}
func (api api) GetMousePos() (int, int) {
	x, y, _ := sdl.GetMouseState()
	return int(x), int(y)
}
func (api api) Window() *sdl.Window { return api.window }
func (api api) Ctx() sdl.GLContext  { return nil }
