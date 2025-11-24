package textrenderer

import (
	"frontend/services/graphics/buffers"
	"frontend/services/graphics/texturearray"

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

	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 5, buffer)
	glyphsWidth := buffers.NewBuffer[float32](
		gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)

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
