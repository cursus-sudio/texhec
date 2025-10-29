package anchorsys

import (
	"frontend/engine/components/anchor"
	"frontend/engine/components/transform"
	"shared/services/datastructures"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
)

func applyChildTransform(
	parent transform.Transform,
	child transform.Transform,
	anchor anchor.ParentAnchor,
) transform.Transform {
	child.SetPos(parent.Pos.
		Add(mgl32.Vec3{
			parent.Size[0] * (anchor.ParentPivot.Point[0] - .5),
			parent.Size[1] * (anchor.ParentPivot.Point[1] - .5),
			parent.Size[2] * (anchor.ParentPivot.Point[2] - .5),
		}).
		Add(anchor.RelativeTransform.Pos).
		Add(anchor.Offset),
	)

	if anchor.RelativeTransform.Rotation != mgl32.QuatIdent() {
		child.SetRotation(parent.Rotation.Mul(anchor.RelativeTransform.Rotation))
	}

	if (anchor.RelativeTransform.Size != mgl32.Vec3{}) {
		child.SetSize(mgl32.Vec3{
			parent.Size[0] * anchor.RelativeTransform.Size[0],
			parent.Size[1] * anchor.RelativeTransform.Size[1],
			parent.Size[2] * anchor.RelativeTransform.Size[2],
		})
	}

	return child
}

type system struct{}

func NewAnchorSystem(logger logger.Logger) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		parentsChildren := map[ecs.EntityID]datastructures.Set[ecs.EntityID]{}
		childParent := map[ecs.EntityID]ecs.EntityID{}

		transformArray := ecs.GetComponentsArray[transform.Transform](w.Components())
		parentAnchorArray := ecs.GetComponentsArray[anchor.ParentAnchor](w.Components())
		{
			transformTransaction := transformArray.Transaction()

			onChange := func(ei []ecs.EntityID) {
				for _, parent := range ei {
					children, ok := parentsChildren[parent]
					if !ok {
						continue
					}
					for _, child := range children.Get() {
						childTransform, err := transformArray.GetComponent(child)
						if err != nil {
							childTransform = transform.NewTransform()
						}
						transformTransaction.SaveComponent(child, childTransform)
					}
				}
				if err := transformTransaction.Flush(); err != nil {
					logger.Error(err)
				}
			}

			onRemove := func(ei []ecs.EntityID) {
				for _, parent := range ei {
					children, ok := parentsChildren[parent]
					if !ok {
						continue
					}
					delete(parentsChildren, parent)
					for _, child := range children.Get() {
						delete(childParent, child)
					}
				}
			}

			query := w.QueryEntitiesWithComponents(
				ecs.GetComponentType(transform.Transform{}),
			)
			query.OnAdd(onChange)
			query.OnChange(onChange)
			query.OnRemove(onRemove)
		}

		{

			onAdd := func(ei []ecs.EntityID) {
				for _, child := range ei {
					anchor, err := parentAnchorArray.GetComponent(child)
					if err != nil {
						continue
					}

					set, ok := parentsChildren[anchor.Parent]
					if !ok {
						set = datastructures.NewSet[ecs.EntityID]()
						parentsChildren[anchor.Parent] = set
					}
					set.Add(child)

					childTransform, err := transformArray.GetComponent(child)
					if err != nil {
						childTransform = transform.NewTransform()
					}
					parentTransform, err := transformArray.GetComponent(anchor.Parent)
					if err != nil {
						parentTransform = transform.NewTransform()
					}

					childTransform = applyChildTransform(parentTransform, childTransform, anchor)
					transformArray.DirtySaveComponent(child, childTransform)
				}
			}

			onRemove := func(ei []ecs.EntityID) {
				for _, child := range ei {
					parent, ok := childParent[child]
					if !ok {
						continue
					}
					delete(childParent, child)
					children, ok := parentsChildren[parent]
					children.RemoveElements(child)
					if len(children.Get()) == 0 {
						delete(parentsChildren, parent)
					}
				}
			}

			query := w.QueryEntitiesWithComponents(
				ecs.GetComponentType(transform.Transform{}),
				ecs.GetComponentType(anchor.ParentAnchor{}),
			)

			query.OnAdd(onAdd)
			query.OnChange(func(ei []ecs.EntityID) {
				onRemove(ei)
				onAdd(ei)
			})
			query.OnRemove(onRemove)
		}

		{
			query := w.QueryEntitiesWithComponents(
				ecs.GetComponentType(anchor.ParentAnchor{}),
			)
			listener := func(ei []ecs.EntityID) {
				for _, entity := range ei {
					if _, err := transformArray.GetComponent(entity); err == nil {
						continue
					}
					transformArray.SaveComponent(entity, transform.NewTransform())
				}
			}
			query.OnAdd(listener)
			query.OnChange(listener)
		}
		return nil
	})
}
