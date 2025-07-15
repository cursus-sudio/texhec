package triangle

import (
	"bytes"
	"core/triangle/abstractions/program"
	"core/triangle/abstractions/texture"
	"core/triangle/abstractions/vao"
	"core/triangle/abstractions/vao/ebo"
	"core/triangle/abstractions/vao/vbo"
	_ "embed"

	"github.com/go-gl/gl/v4.5-core/gl"
)

//go:embed square.png
var textureSource []byte

type triangleTools struct {
	Program   program.Program
	VAO       vao.VAO
	Texture   texture.Texture
	Locations locations
}

func NewTriangleTools() (*triangleTools, error) {
	// program
	shaders, err := NewShaders(shaders)
	if err != nil {
		panic(err.Error())
	}
	p, err := program.NewProgram(shaders.vertex, shaders.fragment)
	if err != nil {
		panic(err.Error())
	}
	shaders.vertex.Release()
	shaders.fragment.Release()

	// vao
	VBO := vbo.NewVBO()
	VBO.SetVertices([]vbo.Vertex{
		{Pos: [3]float32{0, 0, 0}, TexturePos: [2]float32{0, 0}},
		{Pos: [3]float32{100, 0, 0}, TexturePos: [2]float32{1, 0}},
		{Pos: [3]float32{0, 100, 0}, TexturePos: [2]float32{0, 1}},
		{Pos: [3]float32{100, 100, 0}, TexturePos: [2]float32{1, 1}},
	})
	EBO := ebo.NewEBO()
	EBO.SetIndices([]ebo.Index{
		0, 1, 3,
		0, 2, 3,
	})
	VAO := vao.NewVAO(VBO, EBO)

	// texture
	t, err := texture.NewTexture(bytes.NewReader(textureSource))
	if err != nil {
		panic(err.Error())
	}

	return &triangleTools{
		Program: p,
		VAO:     VAO,
		Texture: t,
		Locations: locations{
			Resolution: gl.GetUniformLocation(p.ID, gl.Str("resolution\x00")),
			Model:      gl.GetUniformLocation(p.ID, gl.Str("model\x00")),
			Camera:     gl.GetUniformLocation(p.ID, gl.Str("camera\x00")),
			Projection: gl.GetUniformLocation(p.ID, gl.Str("projection\x00")),
		},
	}, nil
}

type locations struct {
	Resolution int32
	Model      int32
	Camera     int32
	Projection int32
}
