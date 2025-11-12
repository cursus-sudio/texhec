package anchorsys

import (
	"frontend/modules/anchor"
	"frontend/modules/transform"
	"shared/services/datastructures"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
)

func applyChildTransform(
	parent transform.TransformComponent,
	child transform.TransformComponent,
	anchor anchor.ParentAnchorComponent,
) transform.TransformComponent {
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

		transformArray := ecs.GetComponentsArray[transform.TransformComponent](w.Components())
		parentAnchorArray := ecs.GetComponentsArray[anchor.ParentAnchorComponent](w.Components())
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
				logger.Warn(transformTransaction.Flush())
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

			query := w.Query().
				Require(ecs.GetComponentType(transform.TransformComponent{})).
				Build()
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

			query := w.Query().
				Require(ecs.GetComponentType(anchor.ParentAnchorComponent{})).
				Track(ecs.GetComponentType(transform.TransformComponent{})).
				Build()

			query.OnAdd(onAdd)
			query.OnChange(func(ei []ecs.EntityID) {
				onRemove(ei)
				onAdd(ei)
			})
			query.OnRemove(onRemove)
		}

		return nil
	})
}
