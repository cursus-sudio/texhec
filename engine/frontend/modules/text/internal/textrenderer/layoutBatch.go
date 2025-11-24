package textrenderer

import (
	"frontend/services/graphics/vao"
	"frontend/services/graphics/vao/vbo"
)

type layoutBatch struct {
	vao           vao.VAO
	vertices      vbo.VBOSetter[Glyph]
	verticesCount int32

	Layout Layout
}

func NewLayoutBatch(
	v vbo.VBOFactory[Glyph],
	layout Layout,
) layoutBatch {
	VBO := v()
	VBO.SetVertices(layout.Glyphs)
	VAO := vao.NewVAO(VBO, nil)
	return layoutBatch{
		vao:           VAO,
		vertices:      VBO,
		verticesCount: int32(len(layout.Glyphs)),

		Layout: layout,
	}
}

func (b layoutBatch) Release() {
	b.vao.Release()
}
