package collisions

import (
	"frontend/modules/collider"
	"frontend/modules/transform"
	"math"
	"shared/services/datastructures"
	"shared/services/ecs"
	"shared/services/logger"

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
	logger               logger.Logger
	world                ecs.World
	transformTransaction transform.TransformTransaction
	chunkSize            float32
	chunks               map[mgl32.Vec2]datastructures.Set[ecs.EntityID]
	entitiesPositions    map[ecs.EntityID][]mgl32.Vec2
}

func newWorldCollider(
	logger logger.Logger,
	world ecs.World,
	transformTransaction transform.TransformTransaction,
	chunkSize float32,
) worldCollider {
	return &worldColliderImpl{
		logger:               logger,
		world:                world,
		transformTransaction: transformTransaction,
		chunkSize:            chunkSize,
		chunks:               make(map[mgl32.Vec2]datastructures.Set[ecs.EntityID]),
		entitiesPositions:    make(map[ecs.EntityID][]mgl32.Vec2),
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
		transform := c.transformTransaction.GetEntity(entity)
		aabb, err := collider.TransformAABB(transform)
		if err != nil {
			c.logger.Warn(err)
			continue
		}
		positions := c.getPositions(aabb)
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
