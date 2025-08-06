package colliders

import (
	// "fmt"
	// "frontend/engine/components/transform"
	// "math"
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ColliderServiceBuilder { return NewBuilder() })

	ioc.RegisterSingleton(b, func(c ioc.Dic) ColliderService {
		s, errs := ioc.Get[ColliderServiceBuilder](c).Build()
		if len(errs) != 0 {
			logger := ioc.Get[logger.Logger](c)
			logger.Fatal(errs...)
		}
		return s
	})

	// pow := func(n, pow float32) float32 {
	// 	return float32(math.Pow(float64(n), float64(pow)))
	// }
	// pow2 := func(n float32) float32 {
	// 	return pow(n, 2)
	// }
	//
	// ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b ColliderServiceBuilder) ColliderServiceBuilder {
	// 	// 3d
	// 	// sphere
	// 	// box
	// 	// ray
	// 	AddHandler(b, func(s1 Sphere, s2 Sphere) Collision {
	// 		distSq := pow2(s1.Pos.X-s2.Pos.X) + pow2(s1.Pos.Y-s2.Pos.Y) + pow2(s1.Pos.Z-s2.Pos.Z)
	// 		radiusSumSq := pow2(s1.R + s2.R)
	// 		if distSq > radiusSumSq {
	// 			return nil
	// 		}
	// 		dist := math.Sqrt(float64(distSq))
	// 		if dist == 0 {
	// 			return NewCollision()
	// 		}
	// 		penetrationDepth := (s1.R + s2.R) - float32(dist)
	// 		normal := transform.NewPos(
	// 			(s2.Pos.X-s1.Pos.X)/float32(dist),
	// 			(s2.Pos.Y-s1.Pos.Y)/float32(dist),
	// 			(s2.Pos.Z-s1.Pos.Z)/float32(dist),
	// 		)
	// 		contactPoint1 := transform.NewPos(
	// 			s1.Pos.X+normal.X*s1.R,
	// 			s1.Pos.Y+normal.Y*s1.R,
	// 			s1.Pos.Z+normal.Z*s1.R,
	// 		)
	// 		contactPoint2 := transform.NewPos(
	// 			s2.Pos.X-normal.X*s2.R,
	// 			s2.Pos.Y-normal.Y*s2.R,
	// 			s2.Pos.Z-normal.Z*s2.R,
	// 		)
	// 		intersection := NewIntersection(
	// 			contactPoint1,
	// 			contactPoint2,
	// 			normal,
	// 			penetrationDepth,
	// 		)
	// 		return NewCollision(intersection)
	// 	})
	// 	// b.AddHandler(Sphere, Sphere, func(s1, s2 Shape) Intersection {
	// 	// 	t1, t2 := s1.Transform(), s2.Transform()
	// 	// 	distSq := pow2(t1.Pos.X-t2.Pos.X) + pow2(t1.Pos.Y-t2.Pos.Y) + pow2(t1.Pos.Z-t2.Pos.Z)
	// 	// 	radius1, radius2 := t1.Size.X/2, t2.Size.X/2
	// 	// 	radiusSumSq := pow2(radius1 + radius2)
	// 	// 	if distSq > radiusSumSq {
	// 	// 		return nil
	// 	// 	}
	// 	// 	dist := math.Sqrt(float64(distSq))
	// 	// 	if dist == 0 {
	// 	// 		return nil
	// 	// 	}
	// 	// 	// penetrationDepth := (radius1 + radius2) - float32(dist)
	// 	// 	normal := transform.NewPos(
	// 	// 		(t2.Pos.X-t1.Pos.X)/float32(dist),
	// 	// 		(t2.Pos.Y-t1.Pos.Y)/float32(dist),
	// 	// 		(t2.Pos.Z-t1.Pos.Z)/float32(dist),
	// 	// 	)
	// 	// 	contactPoint1 := transform.NewPos(
	// 	// 		t1.Pos.X+normal.X*radius1,
	// 	// 		t1.Pos.Y+normal.Y*radius1,
	// 	// 		t1.Pos.Z+normal.Z*radius1,
	// 	// 	)
	// 	// 	contactPoint2 := transform.NewPos(
	// 	// 		t2.Pos.X-normal.X*radius2,
	// 	// 		t2.Pos.Y-normal.Y*radius2,
	// 	// 		t2.Pos.Z-normal.Z*radius2,
	// 	// 	)
	// 	// 	return NewIntersection(
	// 	// 		contactPoint1,
	// 	// 		contactPoint2,
	// 	// 		// normal,
	// 	// 		// penetrationDepth,
	// 	// 	)
	// 	// })
	//
	// 	return b
	// })
	//
	// ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s ColliderService) ColliderService {
	// 	c1 := NewCollider([]Shape{NewSphere(transform.NewPos(0, 0, 0), 5)})
	// 	// c2 := NewCollider([]Shape{NewSphere(transform.NewPos(11, 0, 0), 5)})
	// 	// c2 := NewCollider([]Shape{NewSphere(transform.NewPos(10, 0, 0), 5)})
	// 	c2 := NewCollider([]Shape{NewSphere(transform.NewPos(3, 0, 0), 5)})
	// 	// c2 := NewCollider([]Shape{NewSphere(transform.NewPos(0, 0, 0), 5)})
	//
	// 	collision, err := s.Collides(c1, c2)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if collision == nil {
	// 		panic("deeply funny")
	// 	}
	// 	for _, intersection := range collision.Intersections() {
	// 		fmt.Printf("intersection %v\n", intersection)
	// 		// fmt.Printf(
	// 		// 	"point on a: %v; point on b: %v;normal %v; depth %v;\n",
	// 		// 	intersection.PointOnA(),
	// 		// 	intersection.PointOnB(),
	// 		// 	intersection.Normal(),
	// 		// 	intersection.Depth(),
	// 		// )
	// 	}
	// 	panic("hihi")
	//
	// 	return s
	// })

}
