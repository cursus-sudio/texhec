package projection

type DynamicOrtho struct {
	Near, Far float32
	Zoom      float32
}

func NewDynamicOrtho(near, far float32, zoom float32) DynamicOrtho {
	return DynamicOrtho{
		Near: near,
		Far:  far,
		Zoom: zoom,
	}
}
