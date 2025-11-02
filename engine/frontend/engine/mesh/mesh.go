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

type MeshAsset[Vertex any] interface {
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
) MeshAsset[Vertex] {
	asset := &meshStorageAsset[Vertex]{
		vertices: vertices,
		indices:  indices,
	}
	return asset
}

func (asset *meshStorageAsset[Vertex]) Vertices() []Vertex   { return asset.vertices }
func (asset *meshStorageAsset[Vertex]) Indices() []ebo.Index { return asset.indices }
