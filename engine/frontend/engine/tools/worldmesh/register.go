package worldmesh

import (
	"errors"
	"frontend/engine/components/mesh"
	"frontend/services/assets"
	"frontend/services/graphics/vao"
	"frontend/services/graphics/vao/ebo"
	"frontend/services/graphics/vao/vbo"
)

type WorldMeshRegister[Vertex any] struct {
	Mesh       vao.VAO
	PackedMesh mesh.PackedMesh[Vertex]
	Ranges     map[assets.AssetID]*mesh.MeshRange
}

func (r WorldMeshRegister[Vertex]) Release() {
	r.Mesh.Release()
}

type RegisterFactory[Vertex any] interface {
	New(meshAssets ...assets.AssetID) (WorldMeshRegister[Vertex], error)
}

type registerFactory[Vertex any] struct {
	assetsStorage assets.AssetsStorage
	vboFactory    vbo.VBOFactory[Vertex]
}

func NewRegisterFactory[Vertex any](
	assetsStorage assets.AssetsStorage,
	vboFactory vbo.VBOFactory[Vertex],
) RegisterFactory[Vertex] {
	return &registerFactory[Vertex]{
		assetsStorage: assetsStorage,
		vboFactory:    vboFactory,
	}
}

func (r *registerFactory[Vertex]) New(meshAssets ...assets.AssetID) (WorldMeshRegister[Vertex], error) {
	rangesIndices := map[assets.AssetID]int{}

	meshesToPack := make([]mesh.MeshStorageAsset[Vertex], len(meshAssets))
	for i, assetID := range meshAssets {
		meshAsset, err := assets.StorageGet[mesh.MeshStorageAsset[Vertex]](r.assetsStorage, assetID)
		if err != nil {
			err := errors.Join(
				err,
				errors.New("creating world mesh register"),
			)
			return WorldMeshRegister[Vertex]{}, err
		}
		meshesToPack[i] = meshAsset
		rangesIndices[assetID] = i
	}
	packedMesh := mesh.Pack(meshesToPack...)

	ranges := map[assets.AssetID]*mesh.MeshRange{}
	for assetID, index := range rangesIndices {
		ranges[assetID] = &packedMesh.Ranges[index]
	}

	VBO := r.vboFactory()
	VBO.SetVertices(packedMesh.Vertices)

	EBO := ebo.NewEBO()
	EBO.SetIndices(packedMesh.Indices)

	register := WorldMeshRegister[Vertex]{
		Mesh:       vao.NewVAO(VBO, EBO),
		PackedMesh: packedMesh,
		Ranges:     ranges,
	}

	return register, nil
}
