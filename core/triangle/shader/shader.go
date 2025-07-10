package shader

import (
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"
)

var (
	VertexShader   uint32 = gl.VERTEX_SHADER
	FragmentShader uint32 = gl.FRAGMENT_SHADER
)

type Shader struct {
	ID uint32
}

func NewShader(shaderSource string, shaderType uint32) (Shader, error) {
	shader, err := compileShader(shaderSource+"\x00", shaderType)
	if err != nil {
		return Shader{}, fmt.Errorf("failed to compile vertex shader: %v", err)
	}

	return Shader{ID: shader}, nil
}

func (shader Shader) Release() {
	gl.DeleteShader(shader.ID)
}
