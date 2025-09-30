package projection

type Perspective struct {
	FovY        float32
	AspectRatio float32
	Near, Far   float32
}

func NewPerspective(fovY float32, aspectRatio float32, near, far float32) Perspective {
	return Perspective{FovY: fovY, AspectRatio: aspectRatio, Near: near, Far: far}
}
