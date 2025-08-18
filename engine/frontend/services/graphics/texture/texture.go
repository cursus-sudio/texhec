package texture

import (
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Texture interface {
	ID() uint32
	Use()
	Release()
}

type texture struct {
	id uint32
}

func NewTexture(img image.Image) (Texture, error) {
	t, err := newTexture(img)
	return &texture{
		id: t,
	}, err
}

func (t *texture) ID() uint32 { return t.id }

func (t *texture) Use() {
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

func (t *texture) Release() {
	gl.DeleteTextures(1, &t.id)
}
