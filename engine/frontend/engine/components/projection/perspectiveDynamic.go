package projection

type DynamicPerspective struct {
	FovY      float32
	Near, Far float32
}

func NewDynamicPerspective(fovY float32, near, far float32) DynamicPerspective {
	return DynamicPerspective{
		FovY: fovY,
		Near: near,
		Far:  far,
	}
}
