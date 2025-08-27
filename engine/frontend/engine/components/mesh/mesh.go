package mesh

import (
	"frontend/services/assets"
	"frontend/services/graphics/vao/ebo"
)

type Mesh struct {
	ID assets.AssetID
}

func NewMesh(id assets.AssetID) Mesh {
	return Mesh{ID: id}
}

//

type MeshStorageAsset[Vertex any] interface {
	Vertices() []Vertex
	Indices() []ebo.Index
}

type meshStorageAsset[Vertex any] struct {
	vertices []Vertex
	indices  []ebo.Index
}

func NewMeshStorageAsset[Vertex any](
	vertices []Vertex,
	indices []ebo.Index,
) MeshStorageAsset[Vertex] {
	return &meshStorageAsset[Vertex]{
		vertices: vertices,
		indices:  indices,
	}
}

func (asset *meshStorageAsset[Vertex]) Vertices() []Vertex   { return asset.vertices }
func (asset *meshStorageAsset[Vertex]) Indices() []ebo.Index { return asset.indices }
