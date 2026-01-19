package tilerenderer

import (
	"core/modules/tile"
	_ "embed"
	"engine"
	"engine/modules/render"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/graphics/program"
	"engine/services/graphics/texturearray"
	"engine/services/graphics/vao"
	"engine/services/graphics/vao/vbo"
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
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

type entityBatch struct {
	vao      vao.VAO
	vertices vbo.VBOSetter[tile.Type]
}

type system struct {
	engine.World `inject:"1"`
	Tile         tile.Service              `inject:"1"`
	VboFactory   vbo.VBOFactory[tile.Type] `inject:"1"`

	program      program.Program
	locations    locations
	textureArray texturearray.TextureArray

	dirtySet ecs.DirtySet
	batches  datastructures.SparseArray[ecs.EntityID, entityBatch]
}

type locations struct {
	Mvp   int32 `uniform:"mvp"`   // mat4
	Width int32 `uniform:"width"` // int
	// widthInv and heightInv is 2/width and 2/height
	WidthInv  int32 `uniform:"widthInv"`  // float
	HeightInv int32 `uniform:"heightInv"` // float
}

func (s *system) Listen(render.RenderEvent) {
	// before get
	dirtyEntities := s.dirtySet.Get()
	for _, entity := range dirtyEntities {
		batch, batchOk := s.batches.Get(entity)
		grid, compOk := s.Tile.Grid().Get(entity)

		if !batchOk && !compOk {
			continue
		}
		if batchOk && !compOk {
			batch.vao.Release()
			s.batches.Remove(entity)
			continue
		}
		if batchOk && compOk {
			batch.vertices.SetVertices(grid.GetTiles())
			continue
		}
		if !batchOk && compOk {
			VBO := s.VboFactory()
			VBO.SetVertices(grid.GetTiles())
			VAO := vao.NewVAO(VBO, nil)
			batch := entityBatch{
				VAO,
				VBO,
			}
			s.batches.Set(entity, batch)
			continue
		}
	}

	// render
	w, h := s.Window.Window().GetSize()
	defer func() { gl.Viewport(0, 0, w, h) }()

	s.program.Use()
	s.textureArray.Use()
	for _, entity := range s.batches.GetIndices() {
		batch, ok := s.batches.Get(entity)
		if !ok {
			continue
		}
		batch.vao.Use()

		grid, ok := s.Tile.Grid().Get(entity)
		if !ok {
			continue
		}

		matrix := s.Transform.Mat4(entity)
		groups, _ := s.Groups.Component().Get(entity)

		gl.Uniform1i(s.locations.Width, int32(grid.Width()))
		gl.Uniform1f(s.locations.WidthInv, 2/float32(grid.Width()))
		gl.Uniform1f(s.locations.HeightInv, 2/float32(grid.Height()))

		for _, cameraEntity := range s.Camera.Component().GetEntities() {
			cameraGroups, _ := s.Groups.Component().Get(cameraEntity)
			if !cameraGroups.SharesAnyGroup(groups) {
				continue
			}

			cameraMatrix := s.Camera.Mat4(cameraEntity)
			mvp := cameraMatrix.Mul4(matrix)
			gl.UniformMatrix4fv(s.locations.Mvp, 1, false, &mvp[0])

			gl.Viewport(s.Camera.GetViewport(cameraEntity))
			gl.DrawArrays(gl.POINTS, 0, int32(batch.vertices.Len()))
		}
	}
}
