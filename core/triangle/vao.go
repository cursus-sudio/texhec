package triangle

import (
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Vertex struct {
	Pos   [3]float32
	Color [4]float32
}

func createVAO() uint32 {
	vertices := []Vertex{
		{Pos: [3]float32{100, 0, 0}, Color: [4]float32{1, 0, 0, 1}},
		{Pos: [3]float32{200, 100, 0}, Color: [4]float32{0, 0, 1, 1}},
		{Pos: [3]float32{0, 100, 0}, Color: [4]float32{0, 1, 0, 1}},

		{Pos: [3]float32{400, 300, 0}, Color: [4]float32{1, 1, 0, 1}},
		{Pos: [3]float32{300, 400, 0}, Color: [4]float32{0, 1, 1, 1}},
		{Pos: [3]float32{500, 400, 0}, Color: [4]float32{1, 0, 1, 1}},
	}
	indices := []uint32{
		0, 1, 2,
		2, 3, 4,
		3, 4, 5,
	}
	verticiesSize := int(unsafe.Sizeof(vertices[0]) * uintptr(len(vertices)))

	var triangleVAO uint32
	var VBO uint32

	// vao
	gl.GenVertexArrays(1, &triangleVAO)
	gl.BindVertexArray(triangleVAO)

	// vbo
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)

	gl.BufferData(gl.ARRAY_BUFFER, verticiesSize, gl.Ptr(vertices), gl.STATIC_DRAW)

	fields := []func(i uint32){}
	fields = append(fields, func(i uint32) {
		gl.VertexAttribPointerWithOffset(i, 3, gl.FLOAT, false, int32(unsafe.Sizeof(Vertex{})),
			uintptr(unsafe.Offsetof(Vertex{}.Pos)))
	})
	fields = append(fields, func(i uint32) {
		gl.VertexAttribPointerWithOffset(i, 3, gl.FLOAT, false, int32(unsafe.Sizeof(Vertex{})),
			uintptr(unsafe.Offsetof(Vertex{}.Color)))
	})

	for i, field := range fields {
		index := uint32(i)
		field(index)
		gl.EnableVertexAttribArray(index)
	}

	// ebo
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)                                             // Bind EBO as ELEMENT_ARRAY_BUFFER
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW) // *4 for uint32 size

	// unbind VAO and VBO by assigning false values
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return triangleVAO
}
