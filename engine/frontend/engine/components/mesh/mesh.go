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
	Verticies() []Vertex
	Indicies() []ebo.Index
}

type meshStorageAsset[Vertex any] struct {
	verticies []Vertex
	indicies  []ebo.Index
}

func NewMeshStorageAsset[Vertex any](
	verticies []Vertex,
	indicies []ebo.Index,
) MeshStorageAsset[Vertex] {
	return &meshStorageAsset[Vertex]{
		verticies: verticies,
		indicies:  indicies,
	}
}

func (asset *meshStorageAsset[Vertex]) Verticies() []Vertex   { return asset.verticies }
func (asset *meshStorageAsset[Vertex]) Indicies() []ebo.Index { return asset.indicies }
