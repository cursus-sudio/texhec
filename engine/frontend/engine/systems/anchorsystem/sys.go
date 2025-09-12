package anchorsystem

import (
	"frontend/engine/components/anchor"
	"frontend/engine/components/transform"
	"frontend/services/datastructures"
	"frontend/services/ecs"

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
		Add(anchor.RelativeTransform.Pos),
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

func NewAnchorSystem(world ecs.World) {
	parentsChildren := map[ecs.EntityID]datastructures.Set[ecs.EntityID]{}
	childParent := map[ecs.EntityID]ecs.EntityID{}
	{
		query := world.QueryEntitiesWithComponents(
			ecs.GetComponentType(transform.Transform{}),
		)

		transformArray := ecs.GetComponentsArray[transform.Transform](world.Components())

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
					transformArray.SaveComponent(child, childTransform)
				}
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

		query.OnAdd(onChange)
		query.OnChange(onChange)
		query.OnRemove(onRemove)
	}

	{
		query := world.QueryEntitiesWithComponents(
			ecs.GetComponentType(transform.Transform{}),
			ecs.GetComponentType(anchor.ParentAnchor{}),
		)

		transformArray := ecs.GetComponentsArray[transform.Transform](world.Components())
		parentAnchorArray := ecs.GetComponentsArray[anchor.ParentAnchor](world.Components())

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
					continue
				}
				parentTransform, err := transformArray.GetComponent(anchor.Parent)
				if err != nil {
					continue
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

		query.OnAdd(onAdd)
		query.OnChange(func(ei []ecs.EntityID) {
			onRemove(ei)
			onAdd(ei)
		})
		query.OnRemove(onRemove)
	}
}
