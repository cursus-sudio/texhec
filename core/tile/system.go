package tile

import (
	_ "embed"
	"frontend/engine/components/groups"
	"frontend/engine/components/projection"
	"frontend/engine/systems/render"
	"frontend/engine/tools/cameras"
	"frontend/services/graphics/program"
	"frontend/services/graphics/texturearray"
	"frontend/services/graphics/vao"
	"frontend/services/graphics/vao/vbo"
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

type System struct {
	program   program.Program
	locations locations

	logger logger.Logger

	textureArray  texturearray.TextureArray
	vao           vao.VAO
	vertices      vbo.VBOSetter[TileComponent]
	verticesCount int32

	tileSize  int32
	gridDepth float32

	world       ecs.World
	cameraQuery ecs.LiveQuery
	groupsArray ecs.ComponentsArray[groups.Groups]
	gridGroups  groups.Groups
	cameraCtors cameras.CameraConstructors

	changed     bool
	changeMutex sync.Locker
	tiles       datastructures.SparseArray[ecs.EntityID, TileComponent]
}

type locations struct {
	Camera    int32 `uniform:"camera"`    // mat4
	TileSize  int32 `uniform:"tileSize"`  // int
	GridDepth int32 `uniform:"gridDepth"` // float32
}

func (s *System) Listen(render.RenderEvent) {
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
		camera, err := s.cameraCtors.Get(cameraEntity, ecs.GetComponentType(projection.Ortho{}))
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
