package tilerenderer

import (
	"core/modules/tile"
	"engine"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/graphics/buffers"
	"engine/services/graphics/program"
	"engine/services/graphics/shader"
	"engine/services/graphics/texturearray"
	"engine/services/graphics/vao"
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

//

type TileRenderSystemRegister struct {
	TextureArrayFactory texturearray.Factory `inject:"1"`
	engine.World        `inject:"1"`
	Tile                tile.Service `inject:"1"`

	C ioc.Dic
}

func NewTileRenderSystemRegister(c ioc.Dic) *TileRenderSystemRegister {
	s := ioc.GetServices[*TileRenderSystemRegister](c)
	s.C = c

	return s
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

	gridDirtySet := ecs.NewDirtySet()
	factory.Tile.Grid().AddDirtySet(gridDirtySet)

	tileDirtySet := ecs.NewDirtySet()
	factory.Tile.Tile().AddDirtySet(tileDirtySet)

	s := ioc.GetServices[*system](factory.C)

	s.program = p
	s.vao = vao.NewVAO(nil, nil)
	s.locations = locations
	s.ids = datastructures.NewSparseArray[tile.ID, uint32]()
	s.lodTextureArrays = []texturearray.TextureArray{}

	s.texturesBuffer = buffers.NewBuffer[mgl32.Vec2](gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, 1)

	s.tileTextures = datastructures.NewSparseArray[uint32, mgl32.Vec2]()
	s.textures = datastructures.NewSparseArray[uint32, image.Image]()

	s.tilesDirtySet = tileDirtySet
	s.gridDirtySet = gridDirtySet
	s.batches = datastructures.NewSparseArray[ecs.EntityID, Batch]()

	events.Listen(factory.EventsBuilder, s.ListenRender)
	return nil
}
