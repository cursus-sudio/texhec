package mainpipeline

import (
	"frontend/services/datastructures"
	"frontend/services/ecs"
	"frontend/services/graphics/program"
	"frontend/services/graphics/shader"
	"sync"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type pipelineRegister struct {
	mutex   *sync.RWMutex
	buffers *pipelineBuffers
	program program.Program

	projections datastructures.Set[ecs.ComponentType]
}

func newRegister(projections datastructures.Set[ecs.ComponentType]) (pipelineRegister, error) {
	vert, err := shader.NewShader(vertSource, shader.VertexShader)
	if err != nil {
		return pipelineRegister{}, err
	}
	frag, err := shader.NewShader(fragSource, shader.FragmentShader)
	if err != nil {
		return pipelineRegister{}, err
	}
	p, err := program.NewProgram(vert, frag, nil)
	if err != nil {
		vert.Release()
		frag.Release()
		return pipelineRegister{}, err
	}
	vert.Release()
	frag.Release()

	p.Use()
	texLoc := gl.GetUniformLocation(p.ID(), gl.Str("texs\x00"))
	gl.Uniform1i(texLoc, 1)

	return pipelineRegister{
		mutex:       &sync.RWMutex{},
		buffers:     newBuffers(len(projections.Get())),
		program:     p,
		projections: projections,
	}, nil
}

func (register pipelineRegister) Release() {
	register.buffers.Release()
	register.program.Release()
}
