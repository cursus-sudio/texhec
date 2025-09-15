package mainpipeline

import (
	_ "embed"
	"frontend/services/graphics/program"
	"frontend/services/graphics/shader"
	"shared/services/datastructures"
	"shared/services/ecs"
	"sync"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type value struct {
	mutex *sync.RWMutex

	buffers *pipelineBuffers
	program program.Program

	projections datastructures.Set[ecs.ComponentType]
}

type pipelineRegister struct{ *value }

//go:embed s.vert
var vertSource string

//go:embed s.geom
var geomSource string

//go:embed s.frag
var fragSource string

func newRegister(projections datastructures.Set[ecs.ComponentType]) (pipelineRegister, error) {
	vert, err := shader.NewShader(vertSource, shader.VertexShader)
	if err != nil {
		return pipelineRegister{}, err
	}
	geom, err := shader.NewShader(geomSource, shader.GeomShader)
	if err != nil {
		return pipelineRegister{}, err
	}
	frag, err := shader.NewShader(fragSource, shader.FragmentShader)
	if err != nil {
		return pipelineRegister{}, err
	}
	programID := gl.CreateProgram()
	gl.AttachShader(programID, vert.ID())
	gl.AttachShader(programID, geom.ID())
	gl.AttachShader(programID, frag.ID())

	p, err := program.NewProgram(programID, nil)
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

	return pipelineRegister{&value{
		mutex:       &sync.RWMutex{},
		buffers:     newBuffers(),
		program:     p,
		projections: projections,
	}}, nil
}

func (register pipelineRegister) Release() {
	register.buffers.Release()
	register.program.Release()
}
