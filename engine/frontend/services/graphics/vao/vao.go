package vao

import (
	"frontend/services/graphics/vao/ebo"
	"frontend/services/graphics/vao/vbo"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type VAO interface {
	ID() uint32
	Release()
	Draw()
}

type vao struct {
	id  uint32
	vbo vbo.VBO
	ebo ebo.EBO
}

func NewVAO(VBO vbo.VBO, EBO ebo.EBO) VAO {
	var id uint32
	gl.GenVertexArrays(1, &id)

	gl.BindVertexArray(id)
	VBO.Configure()
	EBO.Configure()
	gl.BindVertexArray(0)

	return &vao{
		id:  id,
		vbo: VBO,
		ebo: EBO,
	}
}

func (vao *vao) ID() uint32 { return vao.id }

func (vao *vao) ReleaseVAO() {
	gl.DeleteVertexArrays(1, &vao.id)
}

func (vao *vao) Release() {
	vao.vbo.Release()
	vao.ebo.Release()
	vao.ReleaseVAO()
}

func (vao *vao) Draw() {
	gl.BindVertexArray(vao.id)
	gl.DrawElementsWithOffset(gl.TRIANGLES, int32(vao.ebo.Len()), gl.UNSIGNED_INT, 0)
	gl.BindVertexArray(0)
}
