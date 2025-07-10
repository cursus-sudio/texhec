package program

import (
	"core/triangle/shader"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Program struct {
	ID uint32
}

func NewProgram(vertexShader, fragmentShader shader.Shader) (Program, error) {
	program, err := createProgram(vertexShader.ID, fragmentShader.ID)
	return Program{ID: program}, err
}

func (p *Program) Draw(draw func()) {
	gl.UseProgram(p.ID)
	draw()
	gl.UseProgram(0)
}

func (p *Program) Release() {
	gl.DeleteProgram(p.ID)
}
