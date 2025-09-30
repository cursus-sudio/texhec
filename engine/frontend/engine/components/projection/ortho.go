package projection

type Ortho struct {
	Width, Height float32
	Near, Far     float32
}

func NewOrtho(w, h, near, far float32) Ortho {
	return Ortho{Width: w, Height: h, Near: min(near, far), Far: max(near, far)}
}
