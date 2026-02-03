package tilerenderer

import (
	"core/modules/tile"
	_ "embed"
	"engine"
	"engine/modules/render"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/graphics/buffers"
	"engine/services/graphics/program"
	"engine/services/graphics/texturearray"
	"engine/services/graphics/vao"
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
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

type system struct {
	engine.World `inject:"1"`
	Tile         tile.Service `inject:"1"`

	program        program.Program
	locations      locations
	ids            datastructures.SparseArray[tile.Type, uint32]
	textureArray   texturearray.TextureArray
	texturesBuffer buffers.Buffer[mgl32.Vec2] // [index, amount]
	// texturesSizeBuffer buffers.Buffer[int32]
	vao vao.VAO

	dirtySet ecs.DirtySet
	batches  datastructures.SparseArray[ecs.EntityID, Batch]
}

type locations struct {
	Mvp    int32 `uniform:"mvp"`    // mat4
	Width  int32 `uniform:"width"`  // uint
	Height int32 `uniform:"height"` // uint
	// widthInv and heightInv is 2/width and 2/height
	WidthInv  int32 `uniform:"widthInv"`  // float
	HeightInv int32 `uniform:"heightInv"` // float
}

func (s *system) ListenRender(render render.RenderEvent) {
	dirtyEntities := s.dirtySet.Get()
	for _, entity := range dirtyEntities {
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

	s.texturesBuffer.Bind()

	s.program.Bind()
	s.vao.Bind()
	s.textureArray.Bind()

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
		// verticesCount := grid.Width() * grid.Height()
		gl.DrawArrays(gl.POINTS, 0, int32(verticesCount))
	}
}
