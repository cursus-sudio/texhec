package tilerenderer

import (
	_ "embed"
	"frontend/modules/camera"
	"frontend/modules/groups"
	"frontend/modules/render"
	"frontend/services/graphics/program"
	"frontend/services/graphics/texturearray"
	"frontend/services/graphics/vao"
	"frontend/services/graphics/vao/vbo"
	"frontend/services/media/window"
	"image"
	"shared/services/datastructures"
	"shared/services/ecs"
	"shared/services/logger"
	"sync"

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

type system struct {
	program   program.Program
	locations locations
	window    window.Api

	logger logger.Logger

	textureArray  texturearray.TextureArray
	vao           vao.VAO
	vertices      vbo.VBOSetter[TileData]
	verticesCount int32

	tileSize  int32
	gridDepth float32

	world       ecs.World
	cameraQuery ecs.LiveQuery
	groupsArray ecs.ComponentsArray[groups.GroupsComponent]
	gridGroups  groups.GroupsComponent
	cameraCtors camera.CameraTool

	changed     bool
	changeMutex sync.Locker
	tiles       datastructures.SparseArray[ecs.EntityID, TileData]
}

type locations struct {
	Camera    int32 `uniform:"camera"`    // mat4
	TileSize  int32 `uniform:"tileSize"`  // int
	GridDepth int32 `uniform:"gridDepth"` // float32
}

func (s *system) Listen(render.RenderEvent) {
	w, h := s.window.Window().GetSize()
	// gl.Viewport(0, 0, w-100, h-100)
	defer func() { gl.Viewport(0, 0, w, h) }()
	if s.changed {
		s.changeMutex.Lock()
		s.vertices.SetVertices(s.tiles.GetValues())
		s.verticesCount = int32(len(s.tiles.GetValues()))
		s.changed = false
		s.changeMutex.Unlock()
	}

	s.program.Use()
	s.textureArray.Use()
	s.vao.Use()

	gl.Uniform1i(s.locations.TileSize, s.tileSize)
	gl.Uniform1f(s.locations.GridDepth, s.gridDepth)

	for _, cameraEntity := range s.cameraQuery.Entities() {
		camera, err := s.cameraCtors.Get(cameraEntity)
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

		gl.DrawArrays(gl.POINTS, 0, s.verticesCount)
	}
}
