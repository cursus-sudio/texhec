package mesh

import (
	"frontend/services/assets"
	"frontend/services/graphics/vao/ebo"
	"frontend/services/graphics/vao/vbo"
)

type Mesh struct {
	ID assets.AssetID
}

func NewMesh(id assets.AssetID) Mesh {
	return Mesh{ID: id}
}

//

type MeshStorageAsset interface {
	Verticies() []vbo.Vertex
	Indicies() []ebo.Index
}

type meshStorageAsset struct {
	verticies []vbo.Vertex
	indicies  []ebo.Index
}

func NewMeshStorageAsset(
	verticies []vbo.Vertex,
	indicies []ebo.Index,
) MeshStorageAsset {
	return &meshStorageAsset{
		verticies: verticies,
		indicies:  indicies,
	}
}

func (asset *meshStorageAsset) Verticies() []vbo.Vertex { return asset.verticies }
func (asset *meshStorageAsset) Indicies() []ebo.Index   { return asset.indicies }
