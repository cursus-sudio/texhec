package material

import (
	"frontend/services/graphics"
	"frontend/services/graphics/vao/ebo"
)

type MeshRange struct {
	firstIndex  uint32
	indexCount  uint32
	firstVertex uint32
}

func (r *MeshRange) DrawCommand(instanceCount uint32, firstInstance uint32) graphics.DrawElementsIndirectCommand {
	return graphics.NewDrawElementsIndirectCommand(
		r.indexCount,
		instanceCount,
		r.firstIndex,
		r.firstVertex,
		firstInstance,
	)
}

type Mesh[Vertex any] struct {
	verts []Vertex
	idx   []ebo.Index
}

type PackedMesh[Vertex any] struct {
	vertices []Vertex
	indices  []ebo.Index
	ranges   []MeshRange
}

func Pack[Vertex any](meshes ...Mesh[Vertex]) PackedMesh[Vertex] {
	p := PackedMesh[Vertex]{}
	for _, m := range meshes {
		var firstVertex = uint32(len(p.vertices))
		var firstIndex = uint32(len(p.indices))

		p.ranges = append(p.ranges, MeshRange{
			firstIndex:  firstIndex,
			indexCount:  uint32(len(m.idx)),
			firstVertex: firstVertex,
		})

		p.vertices = append(p.vertices, m.verts...)
		for _, i := range m.idx {
			p.indices = append(p.indices, i)
		}
	}
	return p
}
