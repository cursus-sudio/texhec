package mainpipeline

import (
	"frontend/services/ecs"
	"frontend/services/graphics/program"
	"frontend/services/graphics/shader"
	"sync"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type materialWorldRegister struct {
	mutex   *sync.RWMutex
	buffers *materialBuffers
	program program.Program

	projections map[ecs.ComponentType]int32
}

func newMaterialWorldRegistry(projections map[ecs.ComponentType]int32) (materialWorldRegister, error) {
	vert, err := shader.NewShader(vertSource, shader.VertexShader)
	if err != nil {
		return materialWorldRegister{}, err
	}
	frag, err := shader.NewShader(fragSource, shader.FragmentShader)
	if err != nil {
		return materialWorldRegister{}, err
	}
	p, err := program.NewProgram(vert, frag, nil)
	if err != nil {
		vert.Release()
		frag.Release()
		return materialWorldRegister{}, err
	}
	vert.Release()
	frag.Release()

	p.Use()
	texLoc := gl.GetUniformLocation(p.ID(), gl.Str("texs\x00"))
	gl.Uniform1i(texLoc, 1)

	return materialWorldRegister{
		mutex:       &sync.RWMutex{},
		buffers:     newMaterialBuffers(len(projections)),
		program:     p,
		projections: projections,
	}, nil
}

func (register materialWorldRegister) Release() {
	register.buffers.Release()
	register.program.Release()
}
