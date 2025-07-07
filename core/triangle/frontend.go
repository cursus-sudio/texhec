package triangle

import (
	"fmt"
	"strings"

	_ "embed"
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

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		infoLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(infoLog))

		return 0, fmt.Errorf("failed to compile %v: %v", source, infoLog)
	}

	return shader, nil
}

//go:embed vertexShader.glsl
var vertexShaderSource string

// Fragment Shader source code

//go:embed fragmentShader.glsl
var fragmentShaderSource string

func createShaderProgram() (uint32, error) {
	// 1. Compile Shaders
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("failed to compile vertex shader: %v", err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
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

func createVAO() uint32 {
	vertices := []float32{
		// Position (x, y, z)
		0.0, 0.5, 0.0, // Top
		-0.5, -0.5, 0.0, // Bottom-left
		0.5, -0.5, 0.0, // Bottom-right
	}

	var triangleVAO uint32

	// 3. Setup VAO and VBO
	var VBO uint32
	gl.GenVertexArrays(1, &triangleVAO) // Generate 1 VAO
	gl.GenBuffers(1, &VBO)              // Generate 1 VBO

	gl.BindVertexArray(triangleVAO) // Bind the VAO first

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO) // Bind the VBO
	// Transfer vertex data to VBO
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Configure vertex attributes
	// position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)

	// Unbind VBO and VAO (optional, but good practice)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return triangleVAO
}

type triangleTools struct {
	ShaderProgram uint32
	TriangleVAO   uint32
}

func NewTools() (*triangleTools, error) {
	shaderProgram, err := createShaderProgram()
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
