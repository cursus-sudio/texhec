package collisions

import (
	"frontend/modules/collider"
	"frontend/modules/transform"

	"github.com/go-gl/mathgl/mgl32"
)

func TransformAABB(t transform.EntityTransform) (collider.AABB, error) {
	pos, err := t.AbsolutePos().Get()
	if err != nil {
		return collider.AABB{}, err
	}
	rot, err := t.AbsoluteRotation().Get()
	if err != nil {
		return collider.AABB{}, err
	}
	size, err := t.AbsoluteSize().Get()
	if err != nil {
		return collider.AABB{}, err
	}
	halfSize := size.Size.Mul(0.5)

	corners := [8]mgl32.Vec3{
		{-1, -1, -1}, {1, -1, -1}, {-1, 1, -1}, {1, 1, -1},
		{-1, -1, 1}, {1, -1, 1}, {-1, 1, 1}, {1, 1, 1},
	}

	var minCorner, maxCorner mgl32.Vec3

	for i, corner := range corners {
		transformedCorner := rot.Rotation.
			Rotate(mgl32.Vec3{corner[0] * halfSize[0], corner[1] * halfSize[1], corner[2] * halfSize[2]}).
			Add(pos.Pos)

		if i == 0 {
			minCorner, maxCorner = transformedCorner, transformedCorner
			continue
		}
		minCorner = mgl32.Vec3{
			min(minCorner[0], transformedCorner[0]),
			min(minCorner[1], transformedCorner[1]),
			min(minCorner[2], transformedCorner[2]),
		}
		maxCorner = mgl32.Vec3{
			max(maxCorner[0], transformedCorner[0]),
			max(maxCorner[1], transformedCorner[1]),
			max(maxCorner[2], transformedCorner[2]),
		}
	}

	return collider.NewAABB(minCorner, maxCorner), nil
}
