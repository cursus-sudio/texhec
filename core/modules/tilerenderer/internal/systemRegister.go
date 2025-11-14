package internal

import (
	"core/modules/tilerenderer"
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

type TileData struct {
	Pos  tilerenderer.TilePosComponent
	Type tilerenderer.TileTextureComponent
}

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

type TileRenderSystemRegister struct {
	logger              logger.Logger
	textures            datastructures.SparseArray[uint32, image.Image]
	textureArrayFactory texturearray.Factory
	vboFactory          vbo.VBOFactory[TileData]
	assetsStorage       assets.AssetsStorage

	tileSize  int32
	gridDepth float32

	groups             groups.GroupsComponent
	cameraCtorsFactory ecs.ToolFactory[camera.CameraTool]
}

func NewTileRenderSystemRegister(
	textureArrayFactory texturearray.Factory,
	logger logger.Logger,
	vboFactory vbo.VBOFactory[TileData],
	assetsStorage assets.AssetsStorage,
	tileSize int32,
	gridDepth float32,
	groups groups.GroupsComponent,
	cameraCtorsFactory ecs.ToolFactory[camera.CameraTool],
) TileRenderSystemRegister {
	return TileRenderSystemRegister{
		logger:              logger,
		textures:            datastructures.NewSparseArray[uint32, image.Image](),
		textureArrayFactory: textureArrayFactory,
		vboFactory:          vboFactory,
		assetsStorage:       assetsStorage,

		tileSize:  tileSize,
		gridDepth: gridDepth,

		groups:             groups,
		cameraCtorsFactory: cameraCtorsFactory,
	}
}

func (service TileRenderSystemRegister) AddType(addedAssets datastructures.SparseArray[uint32, assets.AssetID]) {
	for _, assetIndex := range addedAssets.GetIndices() {
		asset, _ := addedAssets.Get(assetIndex)
		texture, err := assets.StorageGet[render.TextureAsset](service.assetsStorage, asset)
		if err != nil {
			continue
		}

		service.textures.Set(assetIndex, texture.Images()[0])
	}
}

func (factory TileRenderSystemRegister) Register(w ecs.World) error {
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
	tiles := datastructures.NewSparseArray[ecs.EntityID, TileData]()

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

		tileSize:  factory.tileSize,
		gridDepth: factory.gridDepth,

		world:       w,
		groupsArray: ecs.GetComponentsArray[groups.GroupsComponent](w.Components()),
		gridGroups:  factory.groups,
		cameraQuery: w.Query().Require(ecs.GetComponentType(camera.OrthoComponent{})).Build(),
		cameraCtors: factory.cameraCtorsFactory.Build(w),

		changed:     false,
		changeMutex: changeMutex,
		tiles:       tiles,
	}

	tileTypeArray := ecs.GetComponentsArray[tilerenderer.TileTextureComponent](w.Components())
	tilePosArray := ecs.GetComponentsArray[tilerenderer.TilePosComponent](w.Components())
	transformArray := ecs.GetComponentsArray[transform.TransformComponent](w.Components())
	groupsArray := ecs.GetComponentsArray[groups.GroupsComponent](w.Components())

	onChangeOrAdd := func(ei []ecs.EntityID) {
		changeMutex.Lock()
		defer changeMutex.Unlock()
		s.changed = true

		transformTransaction := transformArray.Transaction()
		groupsTransaction := groupsArray.Transaction()

		for _, entity := range ei {
			tileType, err := tileTypeArray.GetComponent(entity)
			if err != nil {
				continue
			}
			tilePos, err := tilePosArray.GetComponent(entity)
			if err != nil {
				continue
			}
			tile := TileData{tilePos, tileType}
			tiles.Set(entity, tile)

			transformTransaction.SaveComponent(entity, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{float32(factory.tileSize), float32(factory.tileSize), 1}).
				SetPos(mgl32.Vec3{
					float32(factory.tileSize)*float32(tile.Pos.X) + float32(factory.tileSize)/2,
					float32(factory.tileSize)*float32(tile.Pos.Y) + float32(factory.tileSize)/2,
					factory.gridDepth + float32(tile.Pos.Z),
				}).Val())
			groupsTransaction.SaveComponent(entity, factory.groups)
		}

		factory.logger.Warn(ecs.FlushMany(transformTransaction, groupsTransaction))
	}
	tileTypeArray.OnAdd(onChangeOrAdd)
	tileTypeArray.OnChange(onChangeOrAdd)
	tileTypeArray.OnRemove(func(ei []ecs.EntityID) {
		changeMutex.Lock()
		defer changeMutex.Unlock()
		s.changed = true

		for _, entity := range ei {
			tiles.Remove(entity)
		}
	})

	events.Listen(w.EventsBuilder(), s.Listen)
	return nil
}
