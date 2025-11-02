package internal

import (
	"core/modules/unit"
	_ "embed"
	"frontend/modules/camera"
	"frontend/modules/groups"
	"frontend/modules/render"
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

type UnitType struct {
	Texture image.Image
}

type system struct {
	program   program.Program
	locations locations

	logger logger.Logger

	textureArray  texturearray.TextureArray
	vao           vao.VAO
	vertices      vbo.VBOSetter[unit.UnitComponent]
	verticesCount int32

	unitSize  int32
	gridDepth float32

	world       ecs.World
	cameraQuery ecs.LiveQuery
	groupsArray ecs.ComponentsArray[groups.GroupsComponent]
	gridGroups  groups.GroupsComponent
	cameraCtors camera.CameraTool

	changed     bool
	changeMutex sync.Locker
	units       datastructures.SparseArray[ecs.EntityID, unit.UnitComponent]
}

type locations struct {
	Camera    int32 `uniform:"camera"`    // mat4
	UnitSize  int32 `uniform:"unitSize"`  // int
	GridDepth int32 `uniform:"gridDepth"` // float32
}

func (s *system) Listen(render.RenderEvent) {
	if s.changed {
		s.changeMutex.Lock()
		s.vertices.SetVertices(s.units.GetValues())
		s.verticesCount = int32(len(s.units.GetValues()))
		s.changed = false
		s.changeMutex.Unlock()
	}

	s.program.Use()
	s.textureArray.Use()
	s.vao.Use()

	gl.Uniform1i(s.locations.UnitSize, s.unitSize)
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
