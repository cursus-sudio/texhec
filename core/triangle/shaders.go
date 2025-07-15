package triangle

import (
	"core/triangle/abstractions/shader"
	"embed"
	"path/filepath"
)

var shadersDir string = "shaders_old/texture"

//go:embed shaders_old/*
var shaders embed.FS

type Shaders struct {
	vertex   shader.Shader
	fragment shader.Shader
}

func NewShaders(fs embed.FS) (Shaders, error) {
	vertSource, err := fs.ReadFile(filepath.Join(shadersDir, "s.vert"))
	if err != nil {
		return Shaders{}, err
	}
	vert, err := shader.NewShader(string(vertSource), shader.VertexShader)
	if err != nil {
		return Shaders{}, err
	}
	fragSource, err := fs.ReadFile(filepath.Join(shadersDir, "s.frag"))
	if err != nil {
		return Shaders{}, err
	}
	frag, err := shader.NewShader(string(fragSource), shader.FragmentShader)
	if err != nil {
		return Shaders{}, err
	}
	return Shaders{
		vertex:   vert,
		fragment: frag,
	}, nil
}
