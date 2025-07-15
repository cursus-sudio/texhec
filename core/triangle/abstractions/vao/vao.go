package vao

import (
	"core/triangle/abstractions/vao/ebo"
	"core/triangle/abstractions/vao/vbo"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type VAO struct {
	ID  uint32
	VBO vbo.VBO
	EBO ebo.EBO
}

func NewVAO(VBO vbo.VBO, EBO ebo.EBO) VAO {
	var vao uint32
	gl.GenVertexArrays(1, &vao)

	gl.BindVertexArray(vao)
	VBO.Configure()
	EBO.Configure()
	gl.BindVertexArray(0)

	return VAO{
		ID:  vao,
		VBO: VBO,
		EBO: EBO,
	}
}

func (vao *VAO) Release() {
	gl.DeleteVertexArrays(1, &vao.ID)
}

func (vao *VAO) ReleaseAll() {
	vao.VBO.Release()
	vao.EBO.Release()
	vao.Release()
}

func (vao *VAO) Draw() {
	gl.BindVertexArray(vao.ID)
	gl.DrawElementsWithOffset(gl.TRIANGLES, int32(vao.EBO.Len), gl.UNSIGNED_INT, 0)
	gl.BindVertexArray(0)
}
