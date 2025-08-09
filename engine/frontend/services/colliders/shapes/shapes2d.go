package shapes

import (
	"frontend/engine/components/transform"
	"frontend/services/colliders"

	"github.com/go-gl/mathgl/mgl32"
)

// circle
// rectangle
// stadium (not needed for now)

type ellipsoid2D struct {
	transform.Transform
}

type Ellipsoid2D struct{ ellipsoid2D }

func (s ellipsoid2D) R() float32 { return max(s.Size[0], s.Size[1]) }

func (s ellipsoid2D) Apply(t transform.Transform) colliders.Shape {
	return Ellipsoid2D{ellipsoid2D{s.Transform.Merge(t)}}
}

func (s ellipsoid2D) Position() mgl32.Vec3 { return s.Pos }

// We currently implement only circles but store it as ellipse for further development
// but currently treat ellipse2D as circle
// func NewEllipsoid2D(t transform.Transform) colliders.Shape {
// 	return Ellipsoid2D{ellipsoid2D{t}}
// }

func NewCircle2D(pos mgl32.Vec3, r float32) colliders.Shape {
	return Ellipsoid2D{ellipsoid2D{
		transform.NewTransform().
			SetPos(pos).
			SetSize(mgl32.Vec3{r, r, r}),
	}}
}

//

type rect2D struct {
	transform.Transform
}

type Rect2D struct{ rect2D }

func (s rect2D) Apply(t transform.Transform) colliders.Shape {
	return Rect2D{rect2D{s.Transform.Merge(t)}}
}

func (s rect2D) Position() mgl32.Vec3 { return s.Pos }

func NewRect2D(t transform.Transform) Rect2D {
	return Rect2D{rect2D{t}}
}
