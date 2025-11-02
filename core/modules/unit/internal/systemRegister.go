package internal

import (
	"core/modules/unit"
	"frontend/modules/camera"
	"frontend/modules/groups"
	"frontend/modules/render"
	"frontend/modules/transform"
	"frontend/services/assets"
	"frontend/services/graphics/program"
	"frontend/services/graphics/shader"
	"frontend/services/graphics/texturearray"
	"frontend/services/graphics/vao"
	"frontend/services/graphics/vao/ebo"
	"frontend/services/graphics/vao/vbo"
	"image"
	"shared/services/datastructures"
	"shared/services/ecs"
	"shared/services/logger"
	"sync"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

type global struct {
	program      program.Program
	textureArray texturearray.TextureArray
	vao          vao.VAO
}

func (g global) Release() {
	g.program.Release()
	g.textureArray.Release()
	g.vao.Release()
}

//

type UnitRenderSystemRegister struct {
	logger              logger.Logger
	textures            datastructures.SparseArray[uint32, image.Image]
	textureArrayFactory texturearray.Factory
	vboFactory          vbo.VBOFactory[unit.UnitComponent]
	assetsStorage       assets.AssetsStorage

	unitSize  int32
	gridDepth float32

	groups             groups.GroupsComponent
	cameraCtorsFactory ecs.ToolFactory[camera.CameraTool]
}

func NewUnitRenderSystemRegister(
	textureArrayFactory texturearray.Factory,
	logger logger.Logger,
	vboFactory vbo.VBOFactory[unit.UnitComponent],
	assetsStorage assets.AssetsStorage,
	unitSize int32,
	gridDepth float32,
	groups groups.GroupsComponent,
	cameraCtorsFactory ecs.ToolFactory[camera.CameraTool],
) UnitRenderSystemRegister {
	return UnitRenderSystemRegister{
		logger:              logger,
		textures:            datastructures.NewSparseArray[uint32, image.Image](),
		textureArrayFactory: textureArrayFactory,
		vboFactory:          vboFactory,
		assetsStorage:       assetsStorage,

		unitSize:  unitSize,
		gridDepth: gridDepth,

		groups:             groups,
		cameraCtorsFactory: cameraCtorsFactory,
	}
}

func (service UnitRenderSystemRegister) AddType(addedAssets datastructures.SparseArray[uint32, assets.AssetID]) {
	for _, assetIndex := range addedAssets.GetIndices() {
		asset, _ := addedAssets.Get(assetIndex)
		texture, err := assets.StorageGet[render.TextureAsset](service.assetsStorage, asset)
		if err != nil {
			continue
		}

		service.textures.Set(assetIndex, texture.Image())
	}
}

func (factory UnitRenderSystemRegister) Register(w ecs.World) error {
	vert, err := shader.NewShader(vertSource, shader.VertexShader)
	if err != nil {
		return err
	}
	defer vert.Release()

	geom, err := shader.NewShader(geomSource, shader.GeomShader)
	if err != nil {
		return err
	}
	defer geom.Release()

	frag, err := shader.NewShader(fragSource, shader.FragmentShader)
	if err != nil {
		return err
	}
	defer frag.Release()

	programID := gl.CreateProgram()
	gl.AttachShader(programID, vert.ID())
	gl.AttachShader(programID, geom.ID())
	gl.AttachShader(programID, frag.ID())

	p, err := program.NewProgram(programID, nil)
	if err != nil {
		return err
	}

	locations, err := program.GetProgramLocations[locations](p)
	if err != nil {
		return err
	}

	textureArray, err := factory.textureArrayFactory.New(factory.textures)
	if err != nil {
		return err
	}

	VBO := factory.vboFactory()
	var EBO ebo.EBO = nil
	VAO := vao.NewVAO(VBO, EBO)

	changeMutex := &sync.Mutex{}
	units := datastructures.NewSparseArray[ecs.EntityID, unit.UnitComponent]()

	g := global{p, textureArray, VAO}
	w.SaveGlobal(g)

	s := system{
		program:   p,
		locations: locations,

		logger: factory.logger,

		textureArray:  textureArray,
		vao:           VAO,
		vertices:      VBO,
		verticesCount: 0,

		unitSize:  factory.unitSize,
		gridDepth: factory.gridDepth,

		world:       w,
		groupsArray: ecs.GetComponentsArray[groups.GroupsComponent](w.Components()),
		gridGroups:  factory.groups,
		cameraQuery: w.Query().Require(ecs.GetComponentType(camera.OrthoComponent{})).Build(),
		cameraCtors: factory.cameraCtorsFactory.Build(w),

		changed:     false,
		changeMutex: changeMutex,
		units:       units,
	}

	unitArray := ecs.GetComponentsArray[unit.UnitComponent](w.Components())
	transformArray := ecs.GetComponentsArray[transform.TransformComponent](w.Components())
	groupsArray := ecs.GetComponentsArray[groups.GroupsComponent](w.Components())

	onChangeOrAdd := func(ei []ecs.EntityID) {
		changeMutex.Lock()
		defer changeMutex.Unlock()
		s.changed = true

		transformTransaction := transformArray.Transaction()
		groupsTransaction := groupsArray.Transaction()

		for _, entity := range ei {
			unit, err := unitArray.GetComponent(entity)
			if err != nil {
				continue
			}
			units.Set(entity, unit)

			transformTransaction.SaveComponent(entity, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{float32(factory.unitSize), float32(factory.unitSize), 1}).
				SetPos(mgl32.Vec3{
					float32(factory.unitSize)*float32(unit.Pos.X) + float32(factory.unitSize)/2,
					float32(factory.unitSize)*float32(unit.Pos.Y) + float32(factory.unitSize)/2,
					factory.gridDepth,
				}).Val())
			groupsTransaction.SaveComponent(entity, factory.groups)

		}

		err := ecs.FlushMany(transformTransaction, groupsTransaction)
		if err != nil {
			factory.logger.Error(err)
		}
	}
	unitArray.OnAdd(onChangeOrAdd)
	unitArray.OnChange(onChangeOrAdd)
	unitArray.OnRemove(func(ei []ecs.EntityID) {
		changeMutex.Lock()
		defer changeMutex.Unlock()
		s.changed = true

		for _, entity := range ei {
			units.Remove(entity)
		}
	})

	events.Listen(w.EventsBuilder(), s.Listen)
	return nil
}
