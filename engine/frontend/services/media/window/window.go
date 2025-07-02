package window

import "github.com/veandco/go-sdl2/sdl"

type Api interface {
	Window() *sdl.Window
	Renderer() *sdl.Renderer
}

type api struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

func newApi(
	window *sdl.Window,
	renderer *sdl.Renderer,
) Api {
	return api{
		window:   window,
		renderer: renderer,
	}
}

func (api api) Window() *sdl.Window     { return api.window }
func (api api) Renderer() *sdl.Renderer { return api.renderer }
