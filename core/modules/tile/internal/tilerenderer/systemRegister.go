package tilerenderer

import (
	"core/modules/tile"
	"engine"
	"engine/services/assets"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/graphics/program"
	"engine/services/graphics/shader"
	"engine/services/graphics/texturearray"
	"engine/services/graphics/vao"
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

//

type TileRenderSystemRegister struct {
	TextureArrayFactory texturearray.Factory `inject:"1"`
	engine.World        `inject:"1"`
	Tile                tile.Service `inject:"1"`

	C        ioc.Dic
	ids      datastructures.SparseArray[tile.Type, uint32]
	textures datastructures.SparseArray[uint32, image.Image]
}

func NewTileRenderSystemRegister(c ioc.Dic) *TileRenderSystemRegister {
	s := ioc.GetServices[*TileRenderSystemRegister](c)
	s.C = c
	s.ids = datastructures.NewSparseArray[tile.Type, uint32]()
	s.textures = datastructures.NewSparseArray[uint32, image.Image]()
	return s
}

func (service *TileRenderSystemRegister) AddType(addedAssets datastructures.SparseArray[tile.Type, assets.AssetID]) {
	for _, assetIndex := range addedAssets.GetIndices() {
		id := uint32(len(service.ids.GetIndices()))
		service.ids.Set(assetIndex, id)
		asset, _ := addedAssets.Get(assetIndex)
		texture, err := assets.GetAsset[tile.BiomAsset](service.Assets, asset)
		if err != nil {
			service.Logger.Warn(err)
			continue
		}

		base := id*15 + 1
		for i, img := range texture.Images() {
			service.textures.Set(base+uint32(i), img)
		}
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

	dirtySet := ecs.NewDirtySet()
	factory.Tile.Grid().AddDirtySet(dirtySet)

	s := ioc.GetServices[*system](factory.C)

	s.program = p
	s.vao = vao.NewVAO(nil, nil)
	s.locations = locations
	s.ids = factory.ids
	s.textureArray = textureArray

	s.dirtySet = dirtySet
	s.batches = datastructures.NewSparseArray[ecs.EntityID, Batch]()

	events.Listen(factory.EventsBuilder, s.Listen)
	return nil
}
