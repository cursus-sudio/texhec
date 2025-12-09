package tilerenderer

import (
	"core/modules/definition"
	"core/modules/tile"
	"engine/modules/camera"
	"engine/modules/groups"
	"engine/modules/render"
	"engine/services/assets"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/graphics/program"
	"engine/services/graphics/shader"
	"engine/services/graphics/texturearray"
	"engine/services/graphics/vao"
	"engine/services/graphics/vao/ebo"
	"engine/services/graphics/vao/vbo"
	"engine/services/logger"
	"engine/services/media/window"
	"image"
	"sync"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
)

type TileData struct {
	PosX, PosY int32
	Type       definition.DefinitionID
}

type global struct {
	program      program.Program
	textureArray texturearray.TextureArray
	layers       []*layer
}

func (g global) Release() {
	g.program.Release()
	g.textureArray.Release()
	for _, layer := range g.layers {
		layer.vao.Release()
	}
}

//

type TileRenderSystemRegister struct {
	logger              logger.Logger
	window              window.Api
	textures            datastructures.SparseArray[uint32, image.Image]
	textureArrayFactory texturearray.Factory
	vboFactory          vbo.VBOFactory[TileData]
	assetsStorage       assets.AssetsStorage

	tileSize  int32
	gridDepth float32
	layers    int32

	groups             groups.GroupsComponent
	cameraCtorsFactory ecs.ToolFactory[camera.Tool]
}

func NewTileRenderSystemRegister(
	textureArrayFactory texturearray.Factory,
	logger logger.Logger,
	window window.Api,
	vboFactory vbo.VBOFactory[TileData],
	assetsStorage assets.AssetsStorage,
	tileSize int32,
	gridDepth float32,
	layers int32,
	groups groups.GroupsComponent,
	cameraCtorsFactory ecs.ToolFactory[camera.Tool],
) TileRenderSystemRegister {
	return TileRenderSystemRegister{
		logger:              logger,
		window:              window,
		textures:            datastructures.NewSparseArray[uint32, image.Image](),
		textureArrayFactory: textureArrayFactory,
		vboFactory:          vboFactory,
		assetsStorage:       assetsStorage,

		tileSize:  tileSize,
		gridDepth: gridDepth,
		layers:    layers,

		groups:             groups,
		cameraCtorsFactory: cameraCtorsFactory,
	}
}

func (service TileRenderSystemRegister) AddType(addedAssets datastructures.SparseArray[definition.DefinitionID, assets.AssetID]) {
	for _, assetIndex := range addedAssets.GetIndices() {
		asset, _ := addedAssets.Get(assetIndex)
		texture, err := assets.StorageGet[render.TextureAsset](service.assetsStorage, asset)
		if err != nil {
			continue
		}

		service.textures.Set(uint32(assetIndex), texture.Images()[0])
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

	changeMutex := &sync.Mutex{}
	layers := []*layer{}
	for i := 0; i < int(factory.layers); i++ {
		VBO := factory.vboFactory()
		var EBO ebo.EBO = nil
		VAO := vao.NewVAO(VBO, EBO)
		layer := &layer{
			VAO,
			VBO,
			0,
			true,
			datastructures.NewSparseArray[ecs.EntityID, TileData](),
		}
		layers = append(layers, layer)
	}

	g := global{p, textureArray, layers}
	w.SaveGlobal(g)

	s := system{
		program:   p,
		locations: locations,
		window:    factory.window,

		logger: factory.logger,

		textureArray: textureArray,
		layers:       layers,

		tileSize:  factory.tileSize,
		gridDepth: factory.gridDepth,

		world:       w,
		groupsArray: ecs.GetComponentsArray[groups.GroupsComponent](w),
		gridGroups:  factory.groups,
		cameraQuery: w.Query().Require(camera.OrthoComponent{}).Build(),
		cameraCtors: factory.cameraCtorsFactory.Build(w),
	}

	linkArray := ecs.GetComponentsArray[definition.DefinitionLinkComponent](w)
	posArray := ecs.GetComponentsArray[tile.PosComponent](w)

	onChangeOrAdd := func(ei []ecs.EntityID) {
		changeMutex.Lock()
		defer changeMutex.Unlock()

		for _, entity := range ei {
			tileType, err := linkArray.GetComponent(entity)
			if err != nil {
				continue
			}
			tilePos, err := posArray.GetComponent(entity)
			if err != nil {
				continue
			}
			layer := s.layers[tilePos.Layer]
			layer.changed = true
			tile := TileData{tilePos.X, tilePos.Y, tileType.DefinitionID}
			layer.tiles.Set(entity, tile)
		}
	}
	query := w.Query().
		Require(definition.DefinitionLinkComponent{}).
		Require(tile.PosComponent{}).
		Build()
	query.OnAdd(onChangeOrAdd)
	query.OnChange(onChangeOrAdd)
	query.OnRemove(func(ei []ecs.EntityID) {
		changeMutex.Lock()
		defer changeMutex.Unlock()

		for _, entity := range ei {
			tilePos, err := posArray.GetComponent(entity)
			if err != nil {
				continue
			}
			layer := s.layers[tilePos.Layer]
			layer.changed = true
			layer.tiles.Remove(entity)
		}
	})

	events.Listen(w.EventsBuilder(), s.Listen)
	return nil
}
