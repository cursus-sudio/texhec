package collisions

import (
	"frontend/engine/collider"
	"frontend/engine/transform"
	"math"
	"shared/services/datastructures"
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type CollidersTrackingService interface {
	Add(entities ...ecs.EntityID)
	Update(entities ...ecs.EntityID)
	Remove(entities ...ecs.EntityID)
}

type worldCollider interface {
	ChunkSize() float32
	Chunks() map[mgl32.Vec2]datastructures.Set[ecs.EntityID]

	CollidersTrackingService
}

type worldColliderImpl struct {
	world             ecs.World
	transformArray    ecs.ComponentsArray[transform.Transform]
	chunkSize         float32
	chunks            map[mgl32.Vec2]datastructures.Set[ecs.EntityID]
	entitiesPositions map[ecs.EntityID][]mgl32.Vec2
}

func newWorldCollider(world ecs.World, chunkSize float32) worldCollider {
	return &worldColliderImpl{
		world:             world,
		transformArray:    ecs.GetComponentsArray[transform.Transform](world.Components()),
		chunkSize:         chunkSize,
		chunks:            make(map[mgl32.Vec2]datastructures.Set[ecs.EntityID]),
		entitiesPositions: make(map[ecs.EntityID][]mgl32.Vec2),
	}
}

func floorF32ToInt(num float32) int {
	return int(math.Floor(float64(num)))
}

func (c *worldColliderImpl) getPositions(aabb collider.AABB) []mgl32.Vec2 {
	minPos, maxPos := aabb.Min.Vec2(), aabb.Max.Vec2()
	minGridX := floorF32ToInt(minPos.X() / c.chunkSize)
	minGridY := floorF32ToInt(minPos.Y() / c.chunkSize)
	maxGridX := floorF32ToInt(maxPos.X() / c.chunkSize)
	maxGridY := floorF32ToInt(maxPos.Y() / c.chunkSize)

	var positions []mgl32.Vec2
	for x := minGridX; x <= maxGridX; x++ {
		for y := minGridY; y <= maxGridY; y++ {
			tileX := float32(x) * c.chunkSize
			tileY := float32(y) * c.chunkSize
			tileCenter := mgl32.Vec2{tileX, tileY}

			positions = append(positions, tileCenter)
		}
	}

	return positions
}

func (c *worldColliderImpl) ChunkSize() float32                                      { return c.chunkSize }
func (c *worldColliderImpl) Chunks() map[mgl32.Vec2]datastructures.Set[ecs.EntityID] { return c.chunks }

func (c *worldColliderImpl) Add(entities ...ecs.EntityID) {
	for _, entity := range entities {
		transformComponent, err := c.transformArray.GetComponent(entity)
		if err != nil {
			transformComponent = transform.NewTransform()
		}
		positions := c.getPositions(collider.TransformAABB(transformComponent))
		c.entitiesPositions[entity] = positions
		for _, position := range positions {
			arr, ok := c.chunks[position]
			if !ok {
				arr = datastructures.NewSet[ecs.EntityID]()
			}
			arr.Add(entity)
			c.chunks[position] = arr
		}
	}
}
func (c *worldColliderImpl) Update(entities ...ecs.EntityID) {
	c.Remove(entities...)
	c.Add(entities...)
}
func (c *worldColliderImpl) Remove(entities ...ecs.EntityID) {
	for _, entity := range entities {
		positions, ok := c.entitiesPositions[entity]
		if !ok {
			continue
		}
		delete(c.entitiesPositions, entity)
		for _, position := range positions {
			arr, ok := c.chunks[position]
			if !ok {
				continue
			}
			arr.RemoveElements(entity)
			if len(arr.Get()) == 0 {
				delete(c.chunks, position)
				continue
			}
		}
	}
}
