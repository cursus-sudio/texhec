package render

import (
	"engine/modules/assets"
	"engine/services/graphics/vao/ebo"
)

type MeshComponent struct {
	ID assets.ID
}

func NewMesh(id assets.ID) MeshComponent {
	return MeshComponent{ID: id}
}

//

type Vertex struct {
	Pos [3]float32
	// normal [3]float32
	TexturePos [2]float32
	// color [4]float32
	// vertexGroups (for animation) []VertexGroupWeight {Name string; weight float32} (weights should add up to one)
}

//

type MeshAsset interface {
	Vertices() []Vertex
	Indices() []ebo.Index
}

type meshAsset struct {
	vertices []Vertex
	indices  []ebo.Index
}

func NewMeshAsset(
	vertices []Vertex,
	indices []ebo.Index,
) MeshAsset {
	asset := &meshAsset{
		vertices: vertices,
		indices:  indices,
	}
	return asset
}

func (asset *meshAsset) Vertices() []Vertex   { return asset.vertices }
func (asset *meshAsset) Indices() []ebo.Index { return asset.indices }
func (a *meshAsset) Release()                 {}
