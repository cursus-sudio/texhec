package material

import (
	"frontend/services/graphics/vao/vbo"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type VBO struct {
	id  uint32
	len int
}

type Vertex struct {
	Pos [3]float32
	// normal [3]float32
	TexturePos [2]float32
	// color [4]float32
	// vertexGroups (for animation) []VertexGroupWeight {Name string; weight float32} (weights should add up to one)
}

func NewVBO() vbo.VBOSetter[Vertex] {
	var id uint32
	gl.GenBuffers(1, &id)
	return &VBO{
		id:  id,
		len: 0,
	}
}

func (vbo *VBO) ID() uint32 { return vbo.id }
func (vbo *VBO) Len() int   { return vbo.len }

func (vbo *VBO) Configure() {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.id)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false,
		int32(unsafe.Sizeof(Vertex{})), uintptr(unsafe.Offsetof(Vertex{}.Pos)))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false,
		int32(unsafe.Sizeof(Vertex{})), uintptr(unsafe.Offsetof(Vertex{}.TexturePos)))
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (vbo *VBO) Release() {
	gl.DeleteBuffers(1, &vbo.id)
}

func (vbo *VBO) SetVertices(vertices []Vertex) {
	vbo.len = len(vertices)
	verticesSize := int(unsafe.Sizeof(vertices[0]) * uintptr(vbo.len))
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.id)
	gl.BufferData(gl.ARRAY_BUFFER, verticesSize, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
