package vao

import (
	"core/triangle/vao/ebo"
	"core/triangle/vao/vbo"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type VAO struct {
	ID  uint32
	vbo *vbo.VBO
	ebo *ebo.EBO
}

func NewVAO() VAO {
	var vao uint32
	gl.GenVertexArrays(1, &vao)

	return VAO{
		ID:  vao,
		vbo: nil,
		ebo: nil,
	}
}

func (vao *VAO) Release() {
	gl.DeleteVertexArrays(1, &vao.ID)
}

func (vao *VAO) ReleaseAll() {
	vao.Release()
	vao.ReleaseVBO()
	vao.ReleaseEBO()
}

func (vao *VAO) ReleaseVBO() {
	if vao.vbo != nil {
		vao.vbo.Release()
	}
}

func (vao *VAO) ReleaseEBO() {
	if vao.ebo != nil {
		vao.ebo.Release()
	}
}

func (vao *VAO) SetVBO(VBO *vbo.VBO) {
	gl.BindVertexArray(vao.ID)
	if VBO != nil {
		VBO.Configure()
	}
	gl.BindVertexArray(0)
	vao.vbo = VBO
}

func (vao *VAO) SetEBO(ebo *ebo.EBO) {
	var eboID uint32
	if ebo != nil {
		eboID = ebo.ID
	}
	gl.BindVertexArray(vao.ID)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, eboID)
	gl.BindVertexArray(0)
	vao.ebo = ebo
}

func (vao *VAO) Draw() {
	if vao.ebo == nil || vao.vbo == nil {
		return
	}
	gl.BindVertexArray(vao.ID)
	gl.DrawElementsWithOffset(gl.TRIANGLES, int32(vao.ebo.Len), gl.UNSIGNED_INT, 0)
	gl.BindVertexArray(0)
}
