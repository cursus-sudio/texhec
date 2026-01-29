package instancing

import (
	_ "embed"
	"engine/modules/camera"
	"engine/modules/groups"
	"engine/modules/render"
	"engine/modules/transform"
	"engine/services/assets"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/graphics/program"
	"engine/services/graphics/shader"
	"engine/services/graphics/texturearray"
	"engine/services/graphics/vao/vbo"
	"engine/services/logger"
	"engine/services/media/window"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

//go:embed s.vert
var vertSource string

//go:embed s.frag
var fragSource string

type locations struct {
	Camera       int32 `uniform:"camera"`
	CameraGroups int32 `uniform:"cameraGroups"`
}

//

type system struct {
	EventsBuilder events.Builder    `inject:"1"`
	World         ecs.World         `inject:"1"`
	Render        render.Service    `inject:"1"`
	Camera        camera.Service    `inject:"1"`
	Groups        groups.Service    `inject:"1"`
	Transform     transform.Service `inject:"1"`

	Assets              assets.Assets                 `inject:"1"`
	Window              window.Api                    `inject:"1"`
	AssetsStorage       assets.AssetsStorage          `inject:"1"`
	Logger              logger.Logger                 `inject:"1"`
	VboFactory          vbo.VBOFactory[render.Vertex] `inject:"1"`
	TextureArrayFactory texturearray.Factory          `inject:"1"`

	// batches
	dirtyEntities   ecs.DirtySet
	entitiesBatches datastructures.SparseArray[ecs.EntityID, batchKey]
	batches         map[batchKey]*batch

	program   program.Program
	locations locations
}

func NewSystem(c ioc.Dic) render.SystemRenderer {
	return ecs.NewSystemRegister(func() error {
		vert, err := shader.NewShader(vertSource, shader.VertexShader)
		if err != nil {
			return err
		}
		defer vert.Release()

		frag, err := shader.NewShader(fragSource, shader.FragmentShader)
		if err != nil {
			return err
		}
		defer frag.Release()

		programID := gl.CreateProgram()
		gl.AttachShader(programID, vert.ID())
		gl.AttachShader(programID, frag.ID())

		p, err := program.NewProgram(programID, nil)
		if err != nil {
			return err
		}

		locations, err := program.GetProgramLocations[locations](p)
		if err != nil {
			return err
		}

		s := ioc.GetServices[*system](c)

		s.dirtyEntities = ecs.NewDirtySet()
		s.entitiesBatches = datastructures.NewSparseArray[ecs.EntityID, batchKey]()
		s.batches = make(map[batchKey]*batch)

		s.program = p
		s.locations = locations

		s.Render.Color().AddDirtySet(s.dirtyEntities)
		s.Render.TextureFrame().AddDirtySet(s.dirtyEntities)
		s.Transform.AddDirtySet(s.dirtyEntities)
		s.Render.Mesh().AddDirtySet(s.dirtyEntities)
		s.Render.Texture().AddDirtySet(s.dirtyEntities)

		events.ListenE(s.EventsBuilder, s.Listen)
		return nil
	})
}

func (s *system) Listen(render.RenderEvent) error {
	var err error
	// batch
	// for dirtyEntity in entities
	//  if exists than add (create batch if it doesn't exist)
	//  else remove
	for _, entity := range s.dirtyEntities.Get() {
		batchKey, batchKeyOk := batchKey{}, true
		if batchKeyOk {
			batchKey.mesh, batchKeyOk = s.Render.Mesh().Get(entity)
		}
		if batchKeyOk {
			batchKey.texture, batchKeyOk = s.Render.Texture().Get(entity)
		}

		oldBatchKey, oldBatchKeyOk := s.entitiesBatches.Get(entity)
		if oldBatchKeyOk && (!batchKeyOk || batchKey != oldBatchKey) {
			oldBatch := s.batches[oldBatchKey]
			oldBatch.Remove(entity)
			s.entitiesBatches.Remove(entity)
		}
		if !batchKeyOk {
			continue
		}
		batch, ok := s.batches[batchKey]
		if !ok {
			batch, err = s.NewBatch(batchKey)
			if err != nil {
				return err
			}
			s.batches[batchKey] = batch
		}
		batch.Upsert(entity)
		s.entitiesBatches.Set(entity, batchKey)
	}

	// render
	// for batch in batches
	//  bind everything
	//  for camera in cameras
	//   bind camera mat4
	//   render
	s.program.Use()
	for _, batch := range s.batches {
		batch.Bind()
		for _, camera := range s.Camera.Component().GetEntities() {
			gl.Viewport(s.Camera.GetViewport(camera))
			camMatrix := s.Camera.Mat4(camera)
			gl.UniformMatrix4fv(s.locations.Camera, 1, false, &camMatrix[0])

			camGroups, _ := s.Groups.Component().Get(camera)
			gl.Uniform1ui(s.locations.CameraGroups, camGroups.Mask)

			batch.Render()
		}
	}

	return nil
}
