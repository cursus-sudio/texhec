package material

import (
	"frontend/engine/tools/worldmesh"
	"frontend/engine/tools/worldtexture"
	"frontend/services/ecs"
	"frontend/services/graphics/program"
	"sync"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type materialWorldRegister struct {
	mutex   *sync.RWMutex
	buffers *materialBuffers

	projections map[ecs.ComponentType]int32
}

func newMaterialWorldRegistry(
	projections map[ecs.ComponentType]int32,
) materialWorldRegister {
	return materialWorldRegister{
		mutex:       &sync.RWMutex{},
		buffers:     newMaterialBuffers(len(projections)),
		projections: projections,
	}
}

func (register materialWorldRegister) Release() {
	register.buffers.Release()
}

func (register materialWorldRegister) Render(world ecs.World, p program.Program) error {
	mesh, err := ecs.GetRegister[worldmesh.WorldMeshRegister[Vertex]](world)
	if err != nil {
		return err
	}

	texture, err := ecs.GetRegister[worldtexture.WorldTextureRegister](world)
	if err != nil {
		return err
	}

	register.mutex.Lock()
	register.buffers.Flush()
	register.mutex.Unlock()

	p.Use()
	mesh.Mesh.Use()
	texture.Bind()
	cmds := register.buffers.cmdBuffer
	gl.BindBuffer(gl.DRAW_INDIRECT_BUFFER, cmds.ID())
	gl.MultiDrawElementsIndirect(gl.TRIANGLES, gl.UNSIGNED_INT, nil, int32(len(cmds.Get())), 0)

	return nil
}
