package media

import (
	"frontend/services/media/audio"
	"frontend/services/media/inputs"
	"frontend/services/media/window"
)

type Api interface {
	Inputs() inputs.Api
	Window() window.Api
	Audio() audio.Api
}

type api struct {
	inputs inputs.Api
	window window.Api
	audio  audio.Api
}

func newApi(
	inputs inputs.Api,
	window window.Api,
	audio audio.Api,
) Api {
	return api{
		inputs: inputs,
		window: window,
		audio:  audio,
	}
}

func (api api) Inputs() inputs.Api { return api.inputs }
func (api api) Window() window.Api { return api.window }
func (api api) Audio() audio.Api   { return api.audio }
