package texture

import (
	"io"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Texture struct {
	ID uint32
}

func NewTexture(file io.Reader) (Texture, error) {
	t, err := newTexture(file)
	return Texture{
		ID: t,
	}, err
}

func (t *Texture) Draw(draw func()) {
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
	draw()
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (t *Texture) Release() {
	gl.DeleteTextures(1, &t.ID)
}
