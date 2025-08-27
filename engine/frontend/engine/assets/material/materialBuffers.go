package material

import (
	"frontend/services/datastructures"
	"frontend/services/ecs"
	"frontend/services/graphics"
	"frontend/services/graphics/buffers"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type materialBuffers struct {
	entities        datastructures.Set[ecs.EntityID]
	modelBuffer     buffers.Buffer[mgl32.Mat4]
	modelProjBuffer buffers.Buffer[int32]
	modelTexBuffer  buffers.Buffer[int32]
	cmdBuffer       buffers.Buffer[graphics.DrawElementsIndirectCommand]
	// currently there is 1 entity 1 command
	// TODO add instancing

	projBuffer buffers.Buffer[mgl32.Mat4]
}

func newMaterialBuffers(projectionsCount int) *materialBuffers {
	m := &materialBuffers{}
	var buffer uint32

	m.entities = datastructures.NewSet[ecs.EntityID]()

	gl.GenBuffers(1, &buffer)
	m.cmdBuffer = buffers.NewBuffer[graphics.DrawElementsIndirectCommand](
		gl.DRAW_INDIRECT_BUFFER, gl.DYNAMIC_DRAW, buffer)

	gl.GenBuffers(1, &buffer)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 1, buffer)
	m.modelTexBuffer = buffers.NewBuffer[int32](
		gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)

	gl.GenBuffers(1, &buffer)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 2, buffer)
	m.modelBuffer = buffers.NewBuffer[mgl32.Mat4](
		gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)

	gl.GenBuffers(1, &buffer)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 3, buffer)
	m.modelProjBuffer = buffers.NewBuffer[int32](
		gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)

	gl.GenBuffers(1, &buffer)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 4, buffer)
	m.projBuffer = buffers.NewBuffer[mgl32.Mat4](gl.SHADER_STORAGE_BUFFER, gl.STATIC_DRAW, buffer)

	for i := 0; i < projectionsCount; i++ {
		m.projBuffer.Add(mgl32.Ident4())
	}

	return m
}

func (m *materialBuffers) Release() {
	m.modelBuffer.Release()
	m.modelProjBuffer.Release()
	m.modelTexBuffer.Release()
	m.cmdBuffer.Release()

	m.projBuffer.Release()
}

func (m *materialBuffers) Flush() {
	m.cmdBuffer.Flush()
	m.modelTexBuffer.Flush()
	m.modelBuffer.Flush()
	m.modelProjBuffer.Flush()

	m.projBuffer.Flush()
}

func (m *materialBuffers) Upsert(
	entity ecs.EntityID,
	cmd graphics.DrawElementsIndirectCommand,
	textureIndex int32,
	model mgl32.Mat4,
	projectionIndex int32,
) {
	if index, ok := m.entities.GetIndex(entity); ok {
		m.cmdBuffer.Set(index, cmd)
		m.modelTexBuffer.Set(index, textureIndex)
		m.modelBuffer.Set(index, model)
		m.modelProjBuffer.Set(index, projectionIndex)
		return
	}
	m.entities.Add(entity)
	m.cmdBuffer.Add(cmd)
	m.modelTexBuffer.Add(textureIndex)
	m.modelBuffer.Add(model)
	m.modelProjBuffer.Add(projectionIndex)
}

func (m *materialBuffers) Remove(entities []ecs.EntityID) {
	indices := []int{}
	for _, entity := range entities {
		index, ok := m.entities.GetIndex(entity)
		if !ok {
			continue
		}
		indices = append(indices, index)
	}
	m.entities.Remove(indices...)
	m.cmdBuffer.Remove(indices...)
	m.modelTexBuffer.Remove(indices...)
	m.modelBuffer.Remove(indices...)
	m.modelProjBuffer.Remove(indices...)
}
