package rotation

import "math"

// rotation is represented in quaternions
type Rotation struct {
	X, Y, Z, W float64
}

func NewRotation(x, y, z, w float64) Rotation {
	return Rotation{X: x, Y: y, Z: z, W: w}
}

// IdentityRotation returns an identity quaternion, representing no rotation.
func IdentityRotation() Rotation {
	return Rotation{W: 1.0, X: 0.0, Y: 0.0, Z: 0.0}
}

// NewRotationAxisAngle creates a Rotation (quaternion) from an axis and an angle (in radians).
// The axis should be a unit vector.
func NewRotationAxisAngle(angle float64, axisX, axisY, axisZ float64) Rotation {
	halfAngle := angle / 2.0
	s := math.Sin(halfAngle) // Sine of half the angle
	return Rotation{
		W: math.Cos(halfAngle), // Cosine of half the angle
		X: axisX * s,
		Y: axisY * s,
		Z: axisZ * s,
	}
}

// Multiply performs quaternion multiplication (q1 * q2).
// This operation combines two rotations.
func (q1 Rotation) Multiply(q2 Rotation) Rotation {
	return Rotation{
		W: q1.W*q2.W - q1.X*q2.X - q1.Y*q2.Y - q1.Z*q2.Z,
		X: q1.W*q2.X + q1.X*q2.W + q1.Y*q2.Z - q1.Z*q2.Y,
		Y: q1.W*q2.Y - q1.X*q2.Z + q1.Y*q2.W + q1.Z*q2.X,
		Z: q1.W*q2.Z + q1.X*q2.Y - q1.Y*q2.X + q1.Z*q2.W,
	}
}

// Normalize normalizes the quaternion, ensuring it's a unit quaternion.
// Unit quaternions are essential for representing rotations correctly.
func (q *Rotation) Normalize() {
	magnitude := math.Sqrt(q.W*q.W + q.X*q.X + q.Y*q.Y + q.Z*q.Z)
	if magnitude > 0 {
		q.W /= magnitude
		q.X /= magnitude
		q.Y /= magnitude
		q.Z /= magnitude
	}
}

// ToMatrix4 converts the quaternion into a 4x4 rotation matrix.
// This matrix can then be used by OpenGL for transformations.
// The matrix is returned in column-major order, as expected by OpenGL.
func (q Rotation) ToMatrix4() []float32 {
	x2 := q.X * q.X
	y2 := q.Y * q.Y
	z2 := q.Z * q.Z
	xy := q.X * q.Y
	xz := q.X * q.Z
	yz := q.Y * q.Z
	wx := q.W * q.X
	wy := q.W * q.Y
	wz := q.W * q.Z

	// OpenGL expects column-major order for gl.MultMatrixf
	return []float32{
		float32(1.0 - 2.0*(y2+z2)), float32(2.0 * (xy + wz)), float32(2.0 * (xz - wy)), 0.0, // Column 0
		float32(2.0 * (xy - wz)), float32(1.0 - 2.0*(x2+z2)), float32(2.0 * (yz + wx)), 0.0, // Column 1
		float32(2.0 * (xz + wy)), float32(2.0 * (yz - wx)), float32(1.0 - 2.0*(x2+y2)), 0.0, // Column 2
		0.0, 0.0, 0.0, 1.0, // Column 3 (Translation components, identity for rotation matrix)
	}
}
