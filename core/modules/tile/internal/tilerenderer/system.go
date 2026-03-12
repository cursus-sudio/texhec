package tilerenderer

import (
	"core/modules/tile"
	_ "embed"
	"engine"
	"engine/modules/assets"
	"engine/modules/render"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/graphics/buffers"
	"engine/services/graphics/program"
	"engine/services/graphics/shader"
	gtexture "engine/services/graphics/texture"
	"engine/services/graphics/texturearray"
	"engine/services/graphics/vao"
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

//go:embed shader.vert
var vertSource string

//go:embed shader.geom
var geomSource string

//go:embed shader.frag
var fragSource string

type TileType struct {
	Texture image.Image
}

type Batch struct {
	buffer buffers.Buffer[int32]
}

func (b *Batch) Release() {
	b.buffer.Release()
}

//

type locations struct {
	Mvp    int32 `uniform:"mvp"`    // mat4
	Width  int32 `uniform:"width"`  // uint
	Height int32 `uniform:"height"` // uint
	// widthInv and heightInv is 2/width and 2/height
	WidthInv  int32 `uniform:"widthInv"`  // float
	HeightInv int32 `uniform:"heightInv"` // float
}

type system struct {
	TextureArrayFactory texturearray.Factory `inject:"1"`
	engine.World        `inject:"1"`
	Tile                tile.Service `inject:"1"`

	program   program.Program
	locations locations
	ids       datastructures.SparseArray[tile.ID, uint32]
	// lod shrinks
	lodTextureArrays []texturearray.TextureArray
	texturesBuffer   buffers.Buffer[mgl32.Vec2] // [index, amount]
	vao              vao.VAO

	tileTextures datastructures.SparseArray[uint32, mgl32.Vec2]
	textures     datastructures.SparseArray[uint32, image.Image]

	tilesDirtySet ecs.DirtySet
	gridDirtySet  ecs.DirtySet
	batches       datastructures.SparseArray[ecs.EntityID, Batch]
}

func NewSystem(c ioc.Dic) error {
	s := ioc.GetServices[*system](c)

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

	s.program = p
	s.vao = vao.NewVAO(nil, nil)
	s.locations = locations
	s.ids = datastructures.NewSparseArray[tile.ID, uint32]()
	s.lodTextureArrays = []texturearray.TextureArray{}

	s.texturesBuffer = buffers.NewBuffer[mgl32.Vec2](gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, 1)

	s.tileTextures = datastructures.NewSparseArray[uint32, mgl32.Vec2]()
	s.textures = datastructures.NewSparseArray[uint32, image.Image]()

	s.tilesDirtySet = ecs.NewDirtySet()
	s.Tile.Tile().AddDirtySet(s.tilesDirtySet)

	s.gridDirtySet = ecs.NewDirtySet()
	s.Tile.Grid().AddDirtySet(s.gridDirtySet)

	s.batches = datastructures.NewSparseArray[ecs.EntityID, Batch]()

	events.Listen(s.EventsBuilder, s.ListenRender)
	return nil
}

func (s *system) ListenRender(render render.RenderEvent) {
	{ // rare reload. it reloads definitions, buffers, texture arrays (not optimal because currently its used once)
		dirtyTiles := s.tilesDirtySet.Get()
		for _, entity := range dirtyTiles {
			tileComp, ok := s.Tile.Tile().Get(entity)
			if !ok {
				continue
			}

			id := uint32(s.ids.Size())
			s.ids.Set(tileComp.ID, id)
			texture, err := assets.GetAsset[tile.BiomAsset](s.Assets, entity)
			if err != nil {
				s.Logger.Warn(err)
				continue
			}

			rangeBase := id*15 + 1
			for i, images := range texture.Images() {
				size := s.textures.Size()
				tileRange := mgl32.Vec2{float32(size), float32(len(images))}
				s.tileTextures.Set(rangeBase+uint32(i), tileRange)

				imageBase := size
				for i, img := range images {
					s.textures.Set(uint32(imageBase+i), img)
				}
			}
		}
		if len(dirtyTiles) != 0 {
			highLodTextureArray, err := s.TextureArrayFactory.New(s.textures)
			if err != nil {
				s.Logger.Warn(err)
				return
			}

			lowLodTextures := datastructures.NewSparseArray[uint32, image.Image]()
			for _, texture := range s.textures.GetIndices() {
				img, _ := s.textures.Get(texture)
				img = gtexture.NewImage(img).Scale(2, 2).Opaque().Image()
				lowLodTextures.Set(texture, img)
			}
			lowLodTextureArray, err := s.TextureArrayFactory.New(lowLodTextures)
			if err != nil {
				s.Logger.Warn(err)
				return
			}

			dirtySet := ecs.NewDirtySet()
			s.Tile.Grid().AddDirtySet(dirtySet)

			for _, t := range s.lodTextureArrays {
				t.Release()
			}

			s.lodTextureArrays = []texturearray.TextureArray{
				highLodTextureArray,
				lowLodTextureArray,
			}

			for _, id := range s.tileTextures.GetIndices() {
				value, _ := s.tileTextures.Get(id)
				s.texturesBuffer.Set(int(id), value)
			}
			s.texturesBuffer.Flush()
		}
		if len(s.lodTextureArrays) == 0 {
			return
		}
	}

	// reload per grid buffers
	for _, entity := range s.gridDirtySet.Get() {
		batch, batchOk := s.batches.Get(entity)
		grid, compOk := s.Tile.Grid().Get(entity)

		if !batchOk && !compOk {
			continue
		}
		if batchOk && !compOk {
			batch.Release()
			s.batches.Remove(entity)
			continue
		}
		if !batchOk && compOk {
			batch = Batch{
				buffers.NewBuffer[int32](gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, 0),
			}
			s.batches.Set(entity, batch)
		}

		for i, tile := range grid.GetTiles() {
			// there is a conflict
			// we use definitionID to define tile and textures used
			// but tile values are tile.Type diffrentiate it
			id, ok := s.ids.Get(tile)
			if !ok {
				continue
			}
			batch.buffer.Set(i, int32(id))
		}
		batch.buffer.Flush()
	}

	// render
	s.texturesBuffer.Bind()

	s.program.Bind()
	s.vao.Bind()
	var lod int
	if ortho, ok := s.Camera.Ortho().Get(render.Camera); ok && ortho.Zoom < .25 {
		lod = 1
	}
	s.lodTextureArrays[lod].Bind()

	cameraGroups, _ := s.Groups.Component().Get(render.Camera)
	cameraMatrix := s.Camera.Mat4(render.Camera)

	for _, entity := range s.batches.GetIndices() {
		batch, ok := s.batches.Get(entity)
		if !ok {
			continue
		}
		if groups, _ := s.Groups.Component().Get(entity); !cameraGroups.SharesAnyGroup(groups) {
			continue
		}
		batch.buffer.Bind()

		grid, _ := s.Tile.Grid().Get(entity)

		gl.Uniform1ui(s.locations.Width, uint32(grid.Width()))
		gl.Uniform1ui(s.locations.Height, uint32(grid.Height()))
		gl.Uniform1f(s.locations.WidthInv, 2/float32(grid.Width()))
		gl.Uniform1f(s.locations.HeightInv, 2/float32(grid.Height()))

		mvp := cameraMatrix.Mul4(s.Transform.Mat4(entity))
		gl.UniformMatrix4fv(s.locations.Mvp, 1, false, &mvp[0])

		verticesCount := (grid.Width() + 1) * (grid.Height() + 1)
		gl.DrawArrays(gl.POINTS, 0, int32(verticesCount))
	}
}
