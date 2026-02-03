package textrenderer

import (
	"engine/services/graphics/buffers"
	"engine/services/graphics/texturearray"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type fontBatch struct {
	glyphsWidth buffers.Buffer[float32]
	textures    texturearray.TextureArray

	font Font
}

func NewFontBatch(
	textureArrayFactory texturearray.Factory,
	font Font,
) (fontBatch, error) {
	textureArray, err := textureArrayFactory.New(font.Images)
	if err != nil {
		return fontBatch{}, err
	}

	glyphsWidth := buffers.NewBuffer[float32](gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, 0)

	for _, index := range font.GlyphsWidth.GetIndices() {
		width, _ := font.GlyphsWidth.Get(index)
		glyphsWidth.Set(int(index), width)
	}
	glyphsWidth.Flush()

	return fontBatch{
		glyphsWidth: glyphsWidth,
		textures:    textureArray,
		font:        font,
	}, nil
}

func (b *fontBatch) Release() {
	b.glyphsWidth.Release()
	b.textures.Release()
}
