package broadcollision

import (
	"frontend/engine/components/collider"
	"github.com/go-gl/mathgl/mgl32"
)

func rayAABBIntersect(r collider.Ray, box collider.AABB, maxDist float32) (bool, float32) {
	var tMin, tMax float32 = 0.0, maxDist

	for i := 0; i < 3; i++ {
		invDir := 1.0 / r.Direction[i]
		t0 := (box.Min[i] - r.Pos[i]) * invDir
		t1 := (box.Max[i] - r.Pos[i]) * invDir

		if invDir < 0.0 {
			t0, t1 = t1, t0
		}

		tMin = max(t0, tMin)
		tMax = min(t1, tMax)

		if tMax <= tMin {
			return false, 0.0
		}
	}
	if tMin > 0 && tMin < maxDist {
		return true, tMin
	}
	return false, 0.0
}

//

func rayTriangleIntersect(r collider.Ray, tri collider.Polygon, maxDist float32) (bool, float32) {
	edge1 := tri.B.Sub(tri.A)
	edge2 := tri.C.Sub(tri.A)

	pvec := r.Direction.Cross(edge2)
	det := edge1.Dot(pvec)

	if det < mgl32.Epsilon {
		return false, 0.0
	}

	invDet := 1.0 / det
	tvec := r.Pos.Sub(tri.A)

	u := tvec.Dot(pvec) * invDet
	if u < 0 || u > 1 {
		return false, 0.0
	}

	qvec := tvec.Cross(edge1)
	v := r.Direction.Dot(qvec) * invDet
	if v < 0 || u+v > 1 {
		return false, 0.0
	}

	t := edge2.Dot(qvec) * invDet
	if t > mgl32.Epsilon && t < maxDist {
		return true, t
	}

	return false, 0.0
}
