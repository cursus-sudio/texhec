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
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type TileData struct {
	PosX, PosY float32
	Type       definition.DefinitionID
}

//

type TileRenderSystemRegister struct {
	TextureArrayFactory texturearray.Factory     `inject:"1"`
	VboFactory          vbo.VBOFactory[TileData] `inject:"1"`

	engine.World `inject:"1"`
	Definition   definition.Service `inject:"1"`

	textures datastructures.SparseArray[uint32, image.Image]

	tileSize  int32
	gridDepth float32
	layers    int32

	groups groups.GroupsComponent
}

func NewTileRenderSystemRegister(c ioc.Dic,
	tileSize int32,
	gridDepth float32,
	layers int32,
	groups groups.GroupsComponent,
) *TileRenderSystemRegister {
	s := ioc.GetServices[*TileRenderSystemRegister](c)
	s.textures = datastructures.NewSparseArray[uint32, image.Image]()
	s.tileSize = tileSize
	s.gridDepth = gridDepth
	s.layers = layers
	s.groups = groups
	return s
}

func (service *TileRenderSystemRegister) AddType(addedAssets datastructures.SparseArray[definition.DefinitionID, assets.AssetID]) {
	for _, assetIndex := range addedAssets.GetIndices() {
		asset, _ := addedAssets.Get(assetIndex)
		texture, err := assets.GetAsset[render.TextureAsset](service.Assets, asset)
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

	textureArray, err := factory.TextureArrayFactory.New(factory.textures)
	if err != nil {
		return err
	}

	layers := []*layer{}
	for i := 0; i < int(factory.layers); i++ {
		VBO := factory.VboFactory()
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
