package mesh

import (
	"frontend/services/assets"
	"frontend/services/graphics/vao"
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
	assets.StorageAsset
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
func (asset *meshStorageAsset) Cache() (assets.CachedAsset, error) {
	vbo := vbo.NewVBO()
	vbo.SetVertices(asset.Verticies())

	ebo := ebo.NewEBO()
	ebo.SetIndices(asset.Indicies())

	vao := vao.NewVAO(vbo, ebo)
	return &meshCachedAsset{vao: vao}, nil
}

//

type MeshCachedAsset interface {
	assets.CachedAsset
	VAO() vao.VAO
}

type meshCachedAsset struct {
	vao vao.VAO
}

func (asset *meshCachedAsset) VAO() vao.VAO { return asset.vao }
func (asset *meshCachedAsset) Release()     { asset.vao.Release() }
