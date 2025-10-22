package projection

type Ortho struct {
	Width, Height float32
	Near, Far     float32
	Zoom          float32
}

func NewOrtho(w, h, near, far float32, zoom float32) Ortho {
	return Ortho{
		Width:  w / zoom,
		Height: h / zoom,
		Near:   min(near, far),
		Far:    max(near, far),
		Zoom:   zoom,
	}
}
