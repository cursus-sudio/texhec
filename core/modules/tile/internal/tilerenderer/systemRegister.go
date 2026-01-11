package tilerenderer

import (
	"core/modules/definition"
	"core/modules/tile"
	"engine"
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

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
)

type TileData struct {
	PosX, PosY float32
	Type       definition.DefinitionID
}

//

type TileRenderSystemRegister struct {
	window              window.Api
	textures            datastructures.SparseArray[uint32, image.Image]
	textureArrayFactory texturearray.Factory
	vboFactory          vbo.VBOFactory[TileData]
	assets              assets.Assets

	engine.World
	Definition definition.Service

	tileSize  int32
	gridDepth float32
	layers    int32

	groups groups.GroupsComponent
}

func NewTileRenderSystemRegister(
	world engine.World,
	definition definition.Service,
	textureArrayFactory texturearray.Factory,
	logger logger.Logger,
	window window.Api,
	vboFactory vbo.VBOFactory[TileData],
	assets assets.Assets,
	tileSize int32,
	gridDepth float32,
	layers int32,
	groups groups.GroupsComponent,
) *TileRenderSystemRegister {
	return &TileRenderSystemRegister{
		window:              window,
		textures:            datastructures.NewSparseArray[uint32, image.Image](),
		textureArrayFactory: textureArrayFactory,
		vboFactory:          vboFactory,
		assets:              assets,

		World:      world,
		Definition: definition,

		tileSize:  tileSize,
		gridDepth: gridDepth,
		layers:    layers,

		groups: groups,
	}
}

func (service *TileRenderSystemRegister) AddType(addedAssets datastructures.SparseArray[definition.DefinitionID, assets.AssetID]) {
	for _, assetIndex := range addedAssets.GetIndices() {
		asset, _ := addedAssets.Get(assetIndex)
		texture, err := assets.GetAsset[render.TextureAsset](service.assets, asset)
		if err != nil {
			continue
		}

		service.textures.Set(uint32(assetIndex), texture.Images()[0])
	}
}

func (factory *TileRenderSystemRegister) Register() error {
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

	dirtySet := ecs.NewDirtySet()
	tilePosArray := ecs.GetComponentsArray[tile.PosComponent](factory.World)
	factory.Definition.Link().AddDirtySet(dirtySet)
	tilePosArray.AddDirtySet(dirtySet)

	s := system{
		program:   p,
		locations: locations,
		window:    factory.window,

		textureArray: textureArray,
		rendered:     datastructures.NewSparseArray[ecs.EntityID, tile.PosComponent](),
		layers:       layers,

		tileSize:  factory.tileSize,
		gridDepth: factory.gridDepth,

		dirtySet:     dirtySet,
		World:        factory.World,
		Definition:   factory.Definition,
		gridGroups:   factory.groups,
		tilePosArray: tilePosArray,
	}

	events.Listen(factory.EventsBuilder, s.Listen)
	return nil
}
