package collisions

import (
	"engine/modules/collider"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func RayAABBIntersect(r collider.Ray, box collider.AABB) (bool, float32) {
	maxDist := r.MaxDistance
	if maxDist == 0 {
		maxDist = math.MaxFloat32
	}
	var tMin, tMax float32 = 0.0, maxDist

	for i := 0; i < 3; i++ {
		if r.Direction[i] == 0 {
			if r.Pos[i] < box.Min[i] || r.Pos[i] > box.Max[i] {
				return false, 0.0
			}
			continue
		}
		invDir := 1.0 / r.Direction[i]
		t0 := (box.Min[i] - r.Pos[i]) * invDir
		t1 := (box.Max[i] - r.Pos[i]) * invDir

		if invDir < 0.0 {
			t0, t1 = t1, t0
		}

		tMin = max(t0, tMin)
		tMax = min(t1, tMax)

		if tMax < tMin || r.MaxDistance == 0 {
			return false, 0.0
		}
	}
	if tMin > 0 && tMin < maxDist {
		return true, tMin
	}
	return false, 0.0
}

//

func RayTriangleIntersect(r collider.Ray, poly collider.Polygon) (bool, float32) {
	edge1 := poly.B.Sub(poly.A)
	edge2 := poly.C.Sub(poly.A)
	normal := edge1.Cross(edge2).Normalize()

	denom := normal.Dot(r.Direction)
	if mgl32.Abs(denom) < mgl32.Epsilon {
		return false, 0
	}

	t := poly.A.Sub(r.Pos).Dot(normal) / denom
	if t < mgl32.Epsilon || (t > r.MaxDistance && r.MaxDistance != 0) {
		return false, 0
	}

	p := r.Pos.Add(r.Direction.Mul(t))

	crossAB := edge1.Cross(p.Sub(poly.A))
	if normal.Dot(crossAB) < 0 {
		return false, 0
	}

	crossBC := poly.C.Sub(poly.B).Cross(p.Sub(poly.B))
	if normal.Dot(crossBC) < 0 {
		return false, 0
	}

	crossCA := poly.A.Sub(poly.C).Cross(p.Sub(poly.C))
	if normal.Dot(crossCA) < 0 {
		return false, 0
	}

	return true, t
}
