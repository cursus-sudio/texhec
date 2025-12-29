package tilerenderer

import (
	"core/modules/tile"
	_ "embed"
	"engine/modules/groups"
	"engine/modules/render"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/graphics/program"
	"engine/services/graphics/texturearray"
	"engine/services/graphics/vao"
	"engine/services/graphics/vao/vbo"
	"engine/services/logger"
	"engine/services/media/window"
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
)

//go:embed shader.vert
var vertSource string

//go:embed shader.geom
var geomSource string

//go:embed shader.frag
var fragSource string

type TileType struct {
	Texture image.Image
}

type layer struct {
	vao           vao.VAO
	vertices      vbo.VBOSetter[TileData]
	verticesCount int32
	changed       bool
	tiles         datastructures.SparseArray[ecs.EntityID, TileData]
}

type system struct {
	program   program.Program
	locations locations
	window    window.Api

	logger logger.Logger

	textureArray texturearray.TextureArray
	rendered     datastructures.SparseArray[ecs.EntityID, tile.PosComponent]
	layers       []*layer

	tileSize  int32
	gridDepth float32

	dirtySet     ecs.DirtySet
	world        tile.World
	tilePosArray ecs.ComponentsArray[tile.PosComponent]
	gridGroups   groups.GroupsComponent
}

type locations struct {
	Camera    int32 `uniform:"camera"`    // mat4
	TileSize  int32 `uniform:"tileSize"`  // int
	GridDepth int32 `uniform:"gridDepth"` // float32
}

func (s *system) Listen(render.RenderEvent) {
	dirtyEntities := s.dirtySet.Get()
	for _, entity := range dirtyEntities {
		if tilePos, ok := s.rendered.Get(entity); ok {
			layer := s.layers[tilePos.Layer]
			layer.changed = true
			layer.tiles.Remove(entity)
			s.rendered.Remove(entity)
		}
		tileType, ok := s.world.Definition().Link().Get(entity)
		if !ok {
			continue
		}
		tilePos, ok := s.tilePosArray.Get(entity)
		if !ok {
			continue
		}
		layer := s.layers[tilePos.Layer]
		layer.changed = true
		tile := TileData{tilePos.X, tilePos.Y, tileType.DefinitionID}
		layer.tiles.Set(entity, tile)
		s.rendered.Set(entity, tilePos)
	}

	w, h := s.window.Window().GetSize()
	defer func() { gl.Viewport(0, 0, w, h) }()
	for _, layer := range s.layers {
		if layer.changed {
			layer.vertices.SetVertices(layer.tiles.GetValues())
			layer.verticesCount = int32(len(layer.tiles.GetValues()))
			layer.changed = false
		}
	}

	s.program.Use()
	s.textureArray.Use()
	for _, layer := range s.layers {
		layer.vao.Use()

		gl.Uniform1i(s.locations.TileSize, s.tileSize)
		gl.Uniform1f(s.locations.GridDepth, s.gridDepth)

		for _, cameraEntity := range s.world.Camera().Component().GetEntities() {
			cameraGroups, ok := s.world.Groups().Component().Get(cameraEntity)
			if !ok {
				cameraGroups = groups.DefaultGroups()
			}

			if !cameraGroups.SharesAnyGroup(s.gridGroups) {
				continue
			}

			cameraMatrix := s.world.Camera().Mat4(cameraEntity)
			gl.UniformMatrix4fv(s.locations.Camera, 1, false, &cameraMatrix[0])

			gl.Viewport(s.world.Camera().GetViewport(cameraEntity))
			gl.DrawArrays(gl.POINTS, 0, layer.verticesCount)
		}
	}
}
