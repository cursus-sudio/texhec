package triangle

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

func createVAO() uint32 {
	vertices := []float32{ // Position (x, y, z)
		0.0, 1.0, 0.0, // Top
		1,
		-0.5, -0.5, 0.0, // Bottom-left
		0,
		0.5, -0.5, 0.0, // Bottom-right
		0,
	}
	verticiesSize := len(vertices) * 4 // 4 because float32 is 4 bytes
	var argumentsCount int32 = 4
	var stride int32 = argumentsCount * 4

	var triangleVAO uint32
	var VBO uint32

	// generate vao and vbo
	gl.GenVertexArrays(1, &triangleVAO)
	gl.GenBuffers(1, &VBO)

	// bind vao and vbo
	gl.BindVertexArray(triangleVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)

	// transfer verticies to VBO
	gl.BufferData(gl.ARRAY_BUFFER, verticiesSize, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Configure vertex attributes
	// Position attribute (layout location 0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, stride, 0)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 1, gl.FLOAT, false, stride, 3*4)
	gl.EnableVertexAttribArray(1)

	// New attribute (vec2, layout location 1)
	// Offset is the size of the position data (3 floats * 4 bytes/float)
	// gl.VertexAttribPointerWithOffset(1, 1, gl.FLOAT, false, stride, 3*4)
	// gl.EnableVertexAttribArray(1)

	// unbind VAO and VBO by assigning false values
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return triangleVAO
}
