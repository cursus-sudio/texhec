package broadcollision

import (
	"frontend/engine/components/collider"
	"frontend/engine/components/transform"
	"frontend/services/datastructures"
	"frontend/services/ecs"
)

type CollidersTrackingService interface {
	Add(entities ...ecs.EntityID)
	Update(entities ...ecs.EntityID)
	Remove(entities ...ecs.EntityID)
}

// TODO 2
type worldCollider interface {
	AABBs() []collider.AABB
	Ranges() []collider.Range
	Entities() []ecs.EntityID

	CollidersTrackingService
}

type node struct {
	Parent           *node
	Collider         collider.AABB
	ChildrenNodes    []node
	ChildrenEntities datastructures.Set[ecs.EntityID]
}

type worldColliderImpl struct {
	world ecs.World

	leavesPerNode int

	node         *node
	entitiesNode map[ecs.EntityID]datastructures.Set[*node]

	// objects
	aabbs  []collider.AABB
	ranges []collider.Range

	// leaf targets
	entities []ecs.EntityID
}

func newWorldCollider(world ecs.World) worldCollider {
	return &worldColliderImpl{
		world:         world,
		leavesPerNode: 8, // for octo tree
	}
}

func (c *worldColliderImpl) Refresh() {
	type RangeToPopulate struct {
		NodeID     *node
		RangeIndex int
	}

	aabbs := []collider.AABB{}
	ranges := []collider.Range{}
	entities := []ecs.EntityID{}

	nodesToVisit := []*node{c.node}
	rangesToPopulate := []RangeToPopulate{}
	for len(nodesToVisit) > 0 {
		node := nodesToVisit[0]
		if len(rangesToPopulate) != 0 {
			if rangeToPopulate := rangesToPopulate[0]; rangeToPopulate.NodeID == node {
				ranges[rangeToPopulate.RangeIndex].First = uint32(len(aabbs))
				rangesToPopulate = rangesToPopulate[1:]
			}
		}
		nodesToVisit = nodesToVisit[1:]

		if node.ChildrenNodes != nil {
			aabbs = append(aabbs, node.Collider)
			count := len(node.ChildrenNodes)
			colliderRange := collider.NewRange(collider.Branch, 0, uint32(count))
			rangeIndex := len(ranges)
			ranges = append(ranges, colliderRange)

			rangesToPopulate = append(rangesToPopulate, RangeToPopulate{node, rangeIndex})
			for i := range node.ChildrenNodes {
				nodesToVisit = append(nodesToVisit, &node.ChildrenNodes[i])
			}
			continue
		}

		if node.ChildrenEntities != nil {
			aabbs = append(aabbs, node.Collider)
			index := len(entities)
			count := len(node.ChildrenEntities.Get())
			colliderRange := collider.NewRange(collider.Leaf, uint32(index), uint32(count))

			ranges = append(ranges, colliderRange)
			entities = append(entities, node.ChildrenEntities.Get()...)
			continue
		}
	}
}

func (c *worldColliderImpl) AABBs() []collider.AABB {
	return nil
	if c.aabbs == nil {
		c.Refresh()
	}
	return c.aabbs
}
func (c *worldColliderImpl) Ranges() []collider.Range {
	return nil
	if c.ranges == nil {
		c.Refresh()
	}
	return c.ranges
}
func (c *worldColliderImpl) Entities() []ecs.EntityID {
	return nil
	if c.entities == nil {
		c.Refresh()
	}
	return c.entities
}

func (c *worldColliderImpl) add(entity ecs.EntityID, aabb collider.AABB) {
	// check size.
	// add to leafest node.
	// if is on one tile: split if there are to many entities
	// if is on many tiles: do nothing or split if there are to many entities
}

func (c *worldColliderImpl) Add(entities ...ecs.EntityID) {
	c.aabbs = nil
	c.ranges = nil
	c.entities = nil
	for _, entity := range entities {
		if _, ok := c.entitiesNode[entity]; ok {
			continue
		}
		transformComponent, err := ecs.GetComponent[transform.Transform](c.world, entity)
		if err != nil {
			continue
		}
		aabb := collider.TransformAABB(transformComponent)
		c.add(entity, aabb)
	}
}
func (c *worldColliderImpl) Update(entities ...ecs.EntityID) {
	c.Remove(entities...)
	c.Add(entities...)
}
func (c *worldColliderImpl) Remove(entities ...ecs.EntityID) {
	// clear collision tree cache
	c.aabbs = nil
	c.ranges = nil
	c.entities = nil

	for _, entity := range entities {
		nodes, ok := c.entitiesNode[entity]
		if !ok {
			continue
		}
		parentized := []*node{}
		delete(c.entitiesNode, entity)
		for _, node := range nodes.Get() {
			node.ChildrenEntities.RemoveElements(entity)
			if node.Parent == nil {
				continue
			}

		checkParent:
			children := datastructures.NewSet[ecs.EntityID]()
			count := 0
			if node.Parent == nil {
				goto finishParentChecking
			}
			for _, child := range node.Parent.ChildrenNodes {
				if len(child.ChildrenNodes) != 0 {
					goto finishParentChecking
				}
				count += len(child.ChildrenNodes)
				if count >= c.leavesPerNode {
					goto finishParentChecking
				}
				for _, entity := range child.ChildrenEntities.Get() {
					children.Add(entity)
				}
			}

			parentized = append(parentized, node)
			node = node.Parent
			node.ChildrenNodes = nil
			node.ChildrenEntities = children
			goto checkParent
		}
	finishParentChecking:
		for _, node := range parentized {
			nodes.RemoveElements(node)
			nodes.Add(node.Parent)
		}
	}
}
