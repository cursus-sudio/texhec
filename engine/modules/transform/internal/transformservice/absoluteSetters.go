package transformservice

import (
	"engine/modules/transform"
	"engine/services/ecs"
)

type absolutePosArray struct {
	t *service
	ecs.ComponentsArray[transform.AbsolutePosComponent]
}

func (t *absolutePosArray) Set(entity ecs.EntityID, absolutePos transform.AbsolutePosComponent) {
	size, _ := t.t.absoluteSizeArray.Get(entity)
	pos := transform.NewPos(absolutePos.Pos.
		Sub(t.t.GetRelativeParentPos(entity)).
		Sub(t.t.GetPivotPos(entity, size)).Elem())

	t.t.posArray.Set(entity, pos)
}

//

type absoluteSizeArray struct {
	t *service
	ecs.ComponentsArray[transform.AbsoluteSizeComponent]
}

func (t *absoluteSizeArray) Set(entity ecs.EntityID, absoluteSize transform.AbsoluteSizeComponent) {
	parentSize := t.t.GetRelativeParentSize(entity)
	size := transform.NewSize(
		absoluteSize.Size[0]/parentSize[0],
		absoluteSize.Size[1]/parentSize[1],
		absoluteSize.Size[2]/parentSize[2],
	)

	t.t.sizeArray.Set(entity, size)
}

//

type absoluteRotationArray struct {
	t *service
	ecs.ComponentsArray[transform.AbsoluteRotationComponent]
}

func (t *absoluteRotationArray) Set(entity ecs.EntityID, absoluteRot transform.AbsoluteRotationComponent) {
	rot := transform.NewRotation(absoluteRot.Rotation.
		Mul(t.t.GetRelativeParentRotation(entity).Inverse()))

	t.t.rotationArray.Set(entity, rot)
}
