package media

import "frontend/services/media/inputs"

type Api interface {
	Draw() DrawApi
	Inputs() inputs.InputsApi
	Window() WindowApi
	Audio() AudioApi
}
