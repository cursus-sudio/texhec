package projection

type DynamicOrtho struct{ Near, Far float32 }

func NewDynamicOrtho(near, far float32) DynamicOrtho {
	return DynamicOrtho{
		Near: near,
		Far:  far,
	}
}
