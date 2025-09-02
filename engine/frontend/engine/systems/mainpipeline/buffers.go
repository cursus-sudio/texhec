package mainpipeline

import (
	"frontend/services/datastructures"
	"frontend/services/ecs"
	"frontend/services/graphics"
	"frontend/services/graphics/buffers"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Vertex struct {
	Pos [3]float32
	// normal [3]float32
	TexturePos [2]float32
	// color [4]float32
	// vertexGroups (for animation) []VertexGroupWeight {Name string; weight float32} (weights should add up to one)
}

type pipelineBuffers struct {
	entities        datastructures.Set[ecs.EntityID]
	matrixBuffer    buffers.Buffer[mgl32.Mat4]
	modelProjBuffer buffers.Buffer[int32]
	modelTexBuffer  buffers.Buffer[int32]
	cmdBuffer       buffers.Buffer[graphics.DrawElementsIndirectCommand]
	// currently there is 1 entity 1 command
	// TODO add instancing

	orthoBuffer       buffers.Buffer[mgl32.Mat4]
	perspectiveBuffer buffers.Buffer[mgl32.Mat4]
}

func newBuffers() *pipelineBuffers {
	m := &pipelineBuffers{}
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
	m.matrixBuffer = buffers.NewBuffer[mgl32.Mat4](
		gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)

	gl.GenBuffers(1, &buffer)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 3, buffer)
	m.modelProjBuffer = buffers.NewBuffer[int32](
		gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)

	gl.GenBuffers(1, &buffer)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 4, buffer)
	m.orthoBuffer = buffers.NewBuffer[mgl32.Mat4](
		gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)

	gl.GenBuffers(1, &buffer)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 5, buffer)
	m.perspectiveBuffer = buffers.NewBuffer[mgl32.Mat4](
		gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, buffer)

	return m
}

func (m *pipelineBuffers) Release() {
	m.matrixBuffer.Release()
	m.modelProjBuffer.Release()
	m.modelTexBuffer.Release()
	m.cmdBuffer.Release()

	m.orthoBuffer.Release()
	m.perspectiveBuffer.Release()
}

func (m *pipelineBuffers) Flush() {
	m.cmdBuffer.Flush()
	m.modelTexBuffer.Flush()
	m.matrixBuffer.Flush()
	m.modelProjBuffer.Flush()

	m.orthoBuffer.Flush()
	m.perspectiveBuffer.Flush()
}

func (m *pipelineBuffers) Upsert(
	entity ecs.EntityID,
	cmd graphics.DrawElementsIndirectCommand,
	textureIndex int32,
	modelMatrix mgl32.Mat4,
	projectionIndex int32,
) {
	if index, ok := m.entities.GetIndex(entity); ok {
		m.cmdBuffer.Set(index, cmd)
		m.modelTexBuffer.Set(index, textureIndex)
		m.matrixBuffer.Set(index, modelMatrix)
		m.modelProjBuffer.Set(index, projectionIndex)
		return
	}
	m.entities.Add(entity)
	m.cmdBuffer.Add(cmd)
	m.modelTexBuffer.Add(textureIndex)
	m.matrixBuffer.Add(modelMatrix)
	m.modelProjBuffer.Add(projectionIndex)
}

func (m *pipelineBuffers) Remove(entities []ecs.EntityID) {
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
	m.matrixBuffer.Remove(indices...)
	m.modelProjBuffer.Remove(indices...)
}
