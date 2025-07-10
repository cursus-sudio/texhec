package vbo

import (
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type VBO struct {
	ID  uint32
	Len int
}

type Vertex struct {
	Pos        [3]float32
	TexturePos [2]float32
}

func NewVBO() VBO {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	// gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false,
	// 	int32(unsafe.Sizeof(Vertex{})), uintptr(unsafe.Offsetof(Vertex{}.Pos)))
	// gl.EnableVertexAttribArray(0)
	//
	// gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false,
	// 	int32(unsafe.Sizeof(Vertex{})), uintptr(unsafe.Offsetof(Vertex{}.TexturePos)))
	// gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	return VBO{
		ID:  vbo,
		Len: 0,
	}
}

func (vbo *VBO) Configure() {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.ID)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false,
		int32(unsafe.Sizeof(Vertex{})), uintptr(unsafe.Offsetof(Vertex{}.Pos)))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false,
		int32(unsafe.Sizeof(Vertex{})), uintptr(unsafe.Offsetof(Vertex{}.TexturePos)))
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (vbo *VBO) Release() {
	gl.DeleteBuffers(1, &vbo.ID)
}

func (vbo *VBO) SetVertices(vertices []Vertex) {
	verticesLen := len(vertices)
	verticesSize := int(unsafe.Sizeof(vertices[0]) * uintptr(verticesLen))
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.ID)
	gl.BufferData(gl.ARRAY_BUFFER, verticesSize, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	vbo.Len = verticesLen
}
