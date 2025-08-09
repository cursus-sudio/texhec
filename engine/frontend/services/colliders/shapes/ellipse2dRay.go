package shapes

import (
	"frontend/services/colliders"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func (s *collidersService) ellipse2DRayHandler(s1 Ellipsoid2D, s2 Ray) colliders.Collision {
	circle := s1
	ray := s2

	circlePos := circle.Pos
	circleRadius := circle.R()

	rayPos := ray.Pos
	rayDir := ray.Direction

	// Remove the Z-coordinate checks
	// The original code checked for Z-plane alignment, which is not desired for 2D collision with Z as an index.
	// if mgl32.Abs(rayDir.Z()) < 1e-6 {
	//     if mgl32.Abs(rayPos.Z()-circlePos.Z()) > 1e-6 {
	//         return nil
	//     }
	// } else {
	//     tPlane := (circlePos.Z() - rayPos.Z()) / rayDir.Z()
	//     if tPlane < 0 {
	//         return nil
	//     }
	// }

	circlePos2D := mgl32.Vec2{circlePos.X(), circlePos.Y()}
	rayPos2D := mgl32.Vec2{rayPos.X(), rayPos.Y()}
	rayDir2D := mgl32.Vec2{rayDir.X(), rayDir.Y()}

	m2D := rayPos2D.Sub(circlePos2D)

	if rayDir2D.LenSqr() < 1e-6 {
		if m2D.LenSqr() <= circleRadius*circleRadius {
			normal2D := m2D.Normalize()
			contactPoint3D := mgl32.Vec3{rayPos2D.X(), rayPos2D.Y(), circlePos.Z()}
			normal3D := mgl32.Vec3{normal2D.X(), normal2D.Y(), 0}

			intersection := colliders.NewIntersection(
				contactPoint3D,
				contactPoint3D,
				normal3D,
				0,
			)
			return colliders.NewCollision(intersection)
		}
		return nil
	}
	a := rayDir2D.Dot(rayDir2D)
	b := 2 * m2D.Dot(rayDir2D)
	c := m2D.Dot(m2D) - circleRadius*circleRadius

	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return nil
	}

	sqrtDiscriminant := float32(math.Sqrt(float64(discriminant)))
	t1 := (-b - sqrtDiscriminant) / (2 * a)
	t2 := (-b + sqrtDiscriminant) / (2 * a)

	var finalT float32 = -1.0

	if t1 >= 0 {
		finalT = t1
	}

	if t2 >= 0 {
		if finalT < 0 || t2 < finalT {
			finalT = t2
		}
	}

	if finalT < 0 {
		return nil
	}

	intersectionPoint2D := rayPos2D.Add(rayDir2D.Mul(finalT))
	intersectionPoint3D := mgl32.Vec3{intersectionPoint2D.X(), intersectionPoint2D.Y(), circlePos.Z()}

	normal2D := intersectionPoint2D.Sub(circlePos2D).Normalize()
	normal3D := mgl32.Vec3{normal2D.X(), normal2D.Y(), 0}

	intersection := colliders.NewIntersection(
		intersectionPoint3D,
		intersectionPoint3D,
		normal3D,
		0,
	)

	return colliders.NewCollision(intersection)
}
