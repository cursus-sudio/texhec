package shader

import (
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"
)

var (
	VertexShader   uint32 = gl.VERTEX_SHADER
	GeomShader     uint32 = gl.GEOMETRY_SHADER
	FragmentShader uint32 = gl.FRAGMENT_SHADER
)

type Shader interface {
	ID() uint32
	Release()
}

type shader struct {
	id uint32
}

func NewShader(shaderSource string, shaderType uint32) (Shader, error) {
	s, err := compileShader(shaderSource+"\x00", shaderType)
	if err != nil {
		return nil, fmt.Errorf("failed to compile vertex shader: %v", err)
	}

	return &shader{id: s}, nil
}

func (shader *shader) ID() uint32 { return shader.id }

func (shader *shader) Release() {
	gl.DeleteShader(shader.id)
}
