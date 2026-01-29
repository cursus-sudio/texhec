package render

import (
	"engine/services/assets"
	"engine/services/graphics/vao/ebo"
)

type MeshComponent struct {
	ID assets.AssetID
}

func NewMesh(id assets.AssetID) MeshComponent {
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

type meshStorageAsset struct {
	vertices []Vertex
	indices  []ebo.Index
}

func NewMeshStorageAsset(
	vertices []Vertex,
	indices []ebo.Index,
) MeshAsset {
	asset := &meshStorageAsset{
		vertices: vertices,
		indices:  indices,
	}
	return asset
}

func (asset *meshStorageAsset) Vertices() []Vertex   { return asset.vertices }
func (asset *meshStorageAsset) Indices() []ebo.Index { return asset.indices }
func (a *meshStorageAsset) Release()                 {}
