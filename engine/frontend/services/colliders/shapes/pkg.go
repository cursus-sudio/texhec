package shapes

import (
	"frontend/services/colliders"
	"frontend/services/console"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b colliders.ColliderServiceBuilder) colliders.ColliderServiceBuilder {
		cs := newCollidersService(
			ioc.Get[console.Console](c),
		)
		colliders.AddHandler(b, cs.ellipsoid2DEllipsoid2DHandler)
		colliders.AddHandler(b, cs.rect2DEllipse2DHandler)
		colliders.AddHandler(b, cs.rect2DRect2DHandler)
		colliders.AddHandler(b, cs.rect2DRayHandler)
		colliders.AddHandler(b, cs.ellipse2DRayHandler)
		// colliders.AddHandler(b, rayRayHandler) // this is not implemented and there should be no reason to use it
		// implement ray - ellipsoid2d and rect2d handlers
		return b
	})

	// ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s colliders.ColliderService) colliders.ColliderService {
	// 	// return s
	// 	c1 := colliders.NewCollider([]colliders.Shape{NewCircle2D(mgl32.Vec3{0, 0, 0}, 5)})
	// 	c2 := colliders.NewCollider([]colliders.Shape{NewRay(mgl32.Vec3{0, -100, -100}, mgl32.QuatIdent())})
	// 	// c2 := colliders.NewCollider([]colliders.Shape{NewCircle2D(mgl32.Vec3{11, 0, 0}, 5)})
	// 	// c2 := colliders.NewCollider([]colliders.Shape{NewCircle2D(mgl32.Vec3{10, 0, 0}, 5)})
	// 	// c2 := colliders.NewCollider([]colliders.Shape{NewCircle2D(mgl32.Vec3{3, 0, 0}, 5)})
	// 	// c2 := colliders.NewCollider([]colliders.Shape{NewCircle2D(mgl32.Vec3{0, 0, 0}, 5)})
	// 	// c1 := colliders.NewCollider([]colliders.Shape{NewRect2D(transform.NewTransform().
	// 	// 	SetRotation(mgl32.QuatRotate(mgl32.DegToRad(45), mgl32.Vec3{0, 0, 1})).
	// 	// 	SetSize(mgl32.Vec3{10, 10, 0}))})
	// 	// c2 := colliders.NewCollider([]colliders.Shape{NewRect2D(transform.NewTransform().
	// 	// 	SetPos(mgl32.Vec3{14.142134, 0, 0}).
	// 	// 	SetRotation(mgl32.QuatRotate(mgl32.DegToRad(45), mgl32.Vec3{0, 0, 1})).
	// 	// 	SetSize(mgl32.Vec3{10, 10, 0}))})
	//
	// 	collision, err := s.Collides(c1, c2)
	// 	// collision, err := s.Collides(c2, c1)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if collision == nil {
	// 		panic("deeply funny")
	// 	}
	// 	for _, intersection := range collision.Intersections() {
	// 		fmt.Printf("intersection %v\n", intersection)
	// 	}
	// 	panic("hihi")
	//
	// 	return s
	// })
}
