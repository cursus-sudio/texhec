package triangle

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)

// createProgram links shaders into a program
func createProgram(vertexShader, fragmentShader uint32) (uint32, error) {
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		infoLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(infoLog))

		return 0, fmt.Errorf("failed to link program: %v", infoLog)
	}

	return program, nil
}

func createShaderProgram(shaders Shaders) (uint32, error) {
	// 1. Compile Shaders
	vertexShader, err := compileShader(shaders.vertex, gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("failed to compile vertex shader: %v", err)
	}
	fragmentShader, err := compileShader(shaders.fragment, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, fmt.Errorf("failed to compile fragment shader: %v", err)
	}

	// 2. Create Shader Program
	shaderProgram, err := createProgram(vertexShader, fragmentShader)
	if err != nil {
		return 0, fmt.Errorf("failed to link shader program: %v", err)
	}
	gl.DeleteShader(vertexShader)   // Delete individual shaders after linking
	gl.DeleteShader(fragmentShader) // as they are now part of the program
	return shaderProgram, nil
}
