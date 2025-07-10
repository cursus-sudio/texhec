package triangle

import (
	"bytes"
	"core/triangle/program"
	"core/triangle/texture"
	"core/triangle/vao"
	"core/triangle/vao/ebo"
	"core/triangle/vao/vbo"
	_ "embed"
)

//go:embed square.png
var textureSource []byte

type triangleTools struct {
	Program program.Program
	VAO     vao.VAO
	Texture texture.Texture
}

// if err := gl.GetError(); err != gl.NO_ERROR {
// 	panic(err)
// }

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

	// VAO := vao.NewVAO(VBO, EBO)
	VAO := vao.NewVAO()
	VAO.SetVBO(&VBO)
	VAO.SetEBO(&EBO)

	// texture
	t, err := texture.NewTexture(bytes.NewReader(textureSource))
	if err != nil {
		panic(err.Error())
	}

	return &triangleTools{
		Program: p,
		VAO:     VAO,
		Texture: t,
	}, nil
}
