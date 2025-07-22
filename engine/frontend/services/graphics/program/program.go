package program

import (
	"frontend/services/graphics/shader"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Program[Locations any] interface {
	ID() uint32
	Use()
	Locations() Locations
	Release()
}

type program[Locations any] struct {
	id        uint32
	locations Locations
}

func NewProgram[Locations any](vertexShader, fragmentShader shader.Shader) (Program[Locations], error) {
	p, err := createProgram(vertexShader.ID(), fragmentShader.ID())
	if err != nil {
		return nil, err
	}
	locations := createLocations[Locations](p)
	return &program[Locations]{id: p, locations: locations}, nil
}

func (p *program[Locations]) ID() uint32 { return p.id }

func (p *program[Locations]) Use() {
	gl.UseProgram(p.id)
}

func (p *program[Locations]) Locations() Locations {
	return p.locations
}

func (p *program[Locations]) Release() {
	gl.DeleteProgram(p.id)
}
