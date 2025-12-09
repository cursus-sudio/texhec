package tilerenderer

import (
	_ "embed"
	"engine/modules/camera"
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
	layers       []*layer

	tileSize  int32
	gridDepth float32

	world       ecs.World
	cameraQuery ecs.LiveQuery
	groupsArray ecs.ComponentsArray[groups.GroupsComponent]
	gridGroups  groups.GroupsComponent
	cameraCtors camera.Tool
}

type locations struct {
	Camera    int32 `uniform:"camera"`    // mat4
	TileSize  int32 `uniform:"tileSize"`  // int
	GridDepth int32 `uniform:"gridDepth"` // float32
}

func (s *system) Listen(render.RenderEvent) {
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

		for _, cameraEntity := range s.cameraQuery.Entities() {
			camera, err := s.cameraCtors.GetObject(cameraEntity)
			if err != nil {
				continue
			}

			cameraGroups, err := s.groupsArray.GetComponent(cameraEntity)
			if err != nil {
				cameraGroups = groups.DefaultGroups()
			}

			if !cameraGroups.SharesAnyGroup(s.gridGroups) {
				continue
			}

			cameraMatrix := camera.Mat4()
			gl.UniformMatrix4fv(s.locations.Camera, 1, false, &cameraMatrix[0])

			gl.Viewport(camera.Viewport())
			gl.DrawArrays(gl.POINTS, 0, layer.verticesCount)
		}
	}
}
