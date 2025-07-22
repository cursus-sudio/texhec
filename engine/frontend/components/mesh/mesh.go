package mesh

import (
	"frontend/components/transform"
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
	Size() transform.Size
	Verticies() []vbo.Vertex
	Indicies() []ebo.Index
}

type meshStorageAsset struct {
	size      transform.Size
	verticies []vbo.Vertex
	indicies  []ebo.Index
}

func NewMeshStorageAsset(
	size transform.Size,
	verticies []vbo.Vertex,
	indicies []ebo.Index,
) MeshStorageAsset {
	return &meshStorageAsset{
		size:      size,
		verticies: verticies,
		indicies:  indicies,
	}
}

func (asset *meshStorageAsset) Size() transform.Size    { return asset.size }
func (asset *meshStorageAsset) Verticies() []vbo.Vertex { return asset.verticies }
func (asset *meshStorageAsset) Indicies() []ebo.Index   { return asset.indicies }
func (asset *meshStorageAsset) Cache() (assets.CachedAsset, error) {
	vbo := vbo.NewVBO()
	vbo.SetVertices(asset.Verticies())

	ebo := ebo.NewEBO()
	ebo.SetIndices(asset.Indicies())

	vao := vao.NewVAO(vbo, ebo)
	return &meshCachedAsset{size: asset.size, vao: vao}, nil
}

//

type MeshCachedAsset interface {
	assets.CachedAsset
	Size() transform.Size
	VAO() vao.VAO
}

type meshCachedAsset struct {
	size transform.Size
	vao  vao.VAO
}

func (asset *meshCachedAsset) Size() transform.Size { return asset.size }
func (asset *meshCachedAsset) VAO() vao.VAO         { return asset.vao }
func (asset *meshCachedAsset) Release()             { asset.vao.Release() }
