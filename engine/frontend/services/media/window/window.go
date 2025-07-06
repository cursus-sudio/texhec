package window

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Api interface {
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

func (api api) Window() *sdl.Window { return api.window }
func (api api) Ctx() sdl.GLContext  { return nil }
