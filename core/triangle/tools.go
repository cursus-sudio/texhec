package triangle

import "github.com/go-gl/gl/v4.5-core/gl"

type triangleTools struct {
	ShaderProgram uint32
	TriangleVAO   uint32
}

func NewTriangleTools() (*triangleTools, error) {
	shaders, err := NewShaders(shaders)
	if err != nil {
		panic(err.Error())
	}
	shaderProgram, err := createShaderProgram(shaders)
	if err != nil {
		panic(err.Error())
	}
	triangleVAO := createVAO()
	if err := gl.GetError(); err != gl.NO_ERROR {
		panic(err)
	}
	return &triangleTools{
		ShaderProgram: shaderProgram,
		TriangleVAO:   triangleVAO,
	}, nil
}
