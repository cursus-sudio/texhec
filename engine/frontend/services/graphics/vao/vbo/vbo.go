package vbo

import (
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type VBO interface {
	ID() uint32
	Len() int
	Configure()
	Release()
}

type VBOSetter[Vertex any] interface {
	VBO
	SetVertices(vertices []Vertex)
}

type VBOFactory[Vertex any] func() VBOSetter[Vertex]

type vbo[Vertex any] struct {
	id        uint32
	len       int
	configure func()
}

func NewVBO[Vertex any](configure func()) VBOSetter[Vertex] {
	var id uint32
	gl.GenBuffers(1, &id)
	return &vbo[Vertex]{
		id:        id,
		len:       0,
		configure: configure,
	}
}

func (vbo *vbo[Vertex]) ID() uint32 { return vbo.id }
func (vbo *vbo[Vertex]) Len() int   { return vbo.len }

func (vbo *vbo[Vertex]) Configure() {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.id)
	vbo.configure()
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (vbo *vbo[Vertex]) Release() {
	gl.DeleteBuffers(1, &vbo.id)
}

func (vbo *vbo[Vertex]) SetVertices(vertices []Vertex) {
	vbo.len = len(vertices)
	verticesSize := int(unsafe.Sizeof(vertices[0]) * uintptr(vbo.len))
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.id)
	var ptr unsafe.Pointer
	if vbo.len != 0 {
		ptr = gl.Ptr(vertices)
	}
	gl.BufferData(gl.ARRAY_BUFFER, verticesSize, ptr, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
