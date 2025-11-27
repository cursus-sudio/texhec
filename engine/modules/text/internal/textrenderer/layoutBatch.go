package textrenderer

import (
	"engine/services/graphics/vao"
	"engine/services/graphics/vao/vbo"
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
