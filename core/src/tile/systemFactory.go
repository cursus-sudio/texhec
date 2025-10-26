package tile

import (
	"frontend/engine/components/groups"
	"frontend/engine/components/projection"
	"frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/engine/tools/cameras"
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
)

type register struct {
	program      program.Program
	textureArray texturearray.TextureArray
	vao          vao.VAO
}

func (r register) Release() {
	r.program.Release()
	r.textureArray.Release()
	r.vao.Release()
}

//

type TileRenderSystemFactory struct {
	logger              logger.Logger
	textures            datastructures.SparseArray[uint32, image.Image]
	textureArrayFactory texturearray.Factory
	vboFactory          vbo.VBOFactory[TileComponent]
	assetsStorage       assets.AssetsStorage

	tileSize  int32
	gridDepth float32

	groups      groups.Groups
	cameraCtors cameras.CameraConstructors
}

func newTileRenderSystemFactory(
	textureArrayFactory texturearray.Factory,
	logger logger.Logger,
	vboFactory vbo.VBOFactory[TileComponent],
	assetsStorage assets.AssetsStorage,
	tileSize int32,
	gridDepth float32,
	groups groups.Groups,
	cameraCtors cameras.CameraConstructors,
) TileRenderSystemFactory {
	return TileRenderSystemFactory{
		logger:              logger,
		textures:            datastructures.NewSparseArray[uint32, image.Image](),
		textureArrayFactory: textureArrayFactory,
		vboFactory:          vboFactory,
		assetsStorage:       assetsStorage,

		tileSize:  tileSize,
		gridDepth: gridDepth,

		groups:      groups,
		cameraCtors: cameraCtors,
	}
}

func (factory TileRenderSystemFactory) AddType(addedAssets datastructures.SparseArray[uint32, assets.AssetID]) {
	for _, assetIndex := range addedAssets.GetIndices() {
		asset, _ := addedAssets.Get(assetIndex)
		texture, err := assets.StorageGet[texture.TextureAsset](factory.assetsStorage, asset)
		if err != nil {
			continue
		}

		factory.textures.Set(assetIndex, texture.Image())
	}
}

func (factory TileRenderSystemFactory) NewSystem(world ecs.World) (ecs.SystemRegister, error) {
	vert, err := shader.NewShader(vertSource, shader.VertexShader)
	if err != nil {
		return nil, err
	}
	defer vert.Release()

	geom, err := shader.NewShader(geomSource, shader.GeomShader)
	if err != nil {
		return nil, err
	}
	defer geom.Release()

	frag, err := shader.NewShader(fragSource, shader.FragmentShader)
	if err != nil {
		return nil, err
	}
	defer frag.Release()

	programID := gl.CreateProgram()
	gl.AttachShader(programID, vert.ID())
	gl.AttachShader(programID, geom.ID())
	gl.AttachShader(programID, frag.ID())

	p, err := program.NewProgram(programID, nil)
	if err != nil {
		return nil, err
	}

	locations, err := program.GetProgramLocations[locations](p)
	if err != nil {
		return nil, err
	}

	textureArray, err := factory.textureArrayFactory.New(factory.textures)
	if err != nil {
		return nil, err
	}

	VBO := factory.vboFactory()
	var EBO ebo.EBO = nil
	VAO := vao.NewVAO(VBO, EBO)

	changeMutex := &sync.Mutex{}
	tiles := datastructures.NewSparseArray[ecs.EntityID, TileComponent]()

	r := register{p, textureArray, VAO}
	world.SaveRegister(r)

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

		world:       world,
		groupsArray: ecs.GetComponentsArray[groups.Groups](world.Components()),
		gridGroups:  factory.groups,
		cameraQuery: world.QueryEntitiesWithComponents(ecs.GetComponentType(projection.Ortho{})),
		cameraCtors: factory.cameraCtors,

		changed:     false,
		changeMutex: changeMutex,
		tiles:       tiles,
	}

	tileArray := ecs.GetComponentsArray[TileComponent](world.Components())
	transformArray := ecs.GetComponentsArray[transform.Transform](world.Components())
	groupsArray := ecs.GetComponentsArray[groups.Groups](world.Components())

	query := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(TileComponent{}),
	)
	onChangeOrAdd := func(ei []ecs.EntityID) {
		changeMutex.Lock()
		defer changeMutex.Unlock()
		s.changed = true

		transformTransaction := transformArray.Transaction()
		groupsTransaction := groupsArray.Transaction()

		for _, entity := range ei {
			tile, err := tileArray.GetComponent(entity)
			if err != nil {
				continue
			}
			tiles.Set(entity, tile)

			transformTransaction.SaveComponent(entity, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{float32(factory.tileSize), float32(factory.tileSize), 1}).
				SetPos(mgl32.Vec3{
					float32(factory.tileSize)*float32(tile.Pos.X) + float32(factory.tileSize)/2,
					float32(factory.tileSize)*float32(tile.Pos.Y) + float32(factory.tileSize)/2,
					factory.gridDepth,
				}).Val())
			groupsTransaction.SaveComponent(entity, factory.groups)
		}

		err := ecs.FlushMany(transformTransaction, groupsTransaction)
		if err != nil {
			factory.logger.Error(err)
		}
	}
	query.OnAdd(onChangeOrAdd)
	query.OnChange(onChangeOrAdd)
	query.OnRemove(func(ei []ecs.EntityID) {
		changeMutex.Lock()
		defer changeMutex.Unlock()
		s.changed = true

		for _, entity := range ei {
			tiles.Remove(entity)
		}
	})

	return &s, nil
}
