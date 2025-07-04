package transform2d

import "frontend/components/shared/rotation"

type Transform struct {
	Pos      Position
	Rotation rotation.Rotation
	Scale    Scale
}
