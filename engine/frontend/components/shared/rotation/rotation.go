package rotation

// rotation is in radians
type Rotation struct {
	X, Y, Z, W float64
}

func NewRotation(x, y, z, w float64) Rotation {
	return Rotation{X: x, Y: y, Z: z, W: w}
}
