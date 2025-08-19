package texturematerial

import (
	"frontend/services/graphics/vao/ebo"
	"frontend/services/graphics/vao/vbo"
)

type Vertex struct {
	Pos        [3]float32
	TexturePos [2]float32
}

var vertexByteSize int = 8 * 4

type Mesh struct {
	verts []vbo.Vertex
	idx   []ebo.Index
}

type PackedMesh struct {
	vertices []vbo.Vertex
	indices  []ebo.Index
	ranges   []MeshRange
}

var indexByteSize int = 4

func Pack(meshes ...Mesh) PackedMesh {
	p := PackedMesh{}
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
