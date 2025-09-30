package tile

import (
	"frontend/engine/systems/render"
	"frontend/services/graphics"
	"frontend/services/graphics/buffers"
	"image"
	"shared/services/datastructures"
	"shared/services/ecs"
)

type TileType struct {
	Texture image.Image
}

type TileRenderSystem struct {
	TileTypes datastructures.SparseArray[TileTypeID, TileType]
	Tiles     datastructures.SparseArray[ecs.EntityID, TileComponent]

	Entities             datastructures.Set[ecs.EntityID]
	TilePosBuffer        buffers.Buffer[TilePos]
	TileTypeBuffer       buffers.Buffer[TileTypeID]
	TypesTextureIDBuffer buffers.Buffer[int32]
	CommandBuffer        buffers.Buffer[graphics.DrawElementsIndirectCommand]
}

func (TileRenderSystem) Listen(render.RenderEvent) {
}
