package triangle

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)

var shadersDir string = "shaders/gradient"

//go:embed shaders/*
var shaders embed.FS

type Shaders struct {
	vertex   string
	fragment string
}

func NewShaders(fs embed.FS) (Shaders, error) {
	vert, err := fs.ReadFile(filepath.Join(shadersDir, "s.vert"))
	if err != nil {
		return Shaders{}, err
	}
	frag, err := fs.ReadFile(filepath.Join(shadersDir, "s.frag"))
	if err != nil {
		return Shaders{}, err
	}
	return Shaders{
		vertex:   string(vert) + "\x00",
		fragment: string(frag) + "\x00",
	}, nil
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
