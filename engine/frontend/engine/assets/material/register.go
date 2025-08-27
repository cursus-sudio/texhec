package material

import (
	meshcomponent "frontend/engine/components/mesh"
	"frontend/engine/components/projection"
	texturecomponent "frontend/engine/components/texture"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/graphics/buffers"
	"frontend/services/graphics/program"
	"frontend/services/graphics/vao"
	"frontend/services/graphics/vao/ebo"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"sync"
)

type materialWorldRegister struct {
	world ecs.World
	mutex *sync.RWMutex

	entitiesBuffers *entitiesBuffers

	projBuffer buffers.Buffer[mgl32.Mat4]

	// for modelCmdBuffer and modelTexBuffer
	meshes   map[assets.AssetID]int32
	textures map[assets.AssetID]int32

	// render data
	mesh    vao.VAO
	texture uint32

	packedMesh []MeshRange

	projections map[ecs.ComponentType]int32
}

func newMaterialWorldRegistry(
	meshes map[assets.AssetID]int32,
	textures map[assets.AssetID]int32,
	packedMesh PackedMesh[Vertex],
	projections map[ecs.ComponentType]int32,
	texture uint32,
) materialWorldRegister {
	return materialWorldRegister{
		mutex: &sync.RWMutex{},

		entitiesBuffers: newEntitiesBuffers(),
		projBuffer: func() buffers.Buffer[mgl32.Mat4] {
			var bufferID uint32
			gl.GenBuffers(1, &bufferID)
			gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 4, bufferID)
			buffer := buffers.NewBuffer[mgl32.Mat4](gl.SHADER_STORAGE_BUFFER, gl.STATIC_DRAW, bufferID)

			for i := 0; i < len(projections); i++ {
				buffer.Add(mgl32.Ident4())
			}
			return buffer
		}(),

		meshes:   meshes,
		textures: textures,

		mesh: func() vao.VAO {
			vbo := NewVBO()
			vbo.SetVertices(packedMesh.vertices)

			ebo := ebo.NewEBO()
			ebo.SetIndices(packedMesh.indices)

			return vao.NewVAO(vbo, ebo)
		}(),
		texture: texture,

		packedMesh:  packedMesh.ranges,
		projections: projections,
	}
}

func (register materialWorldRegister) CleanUp() {
	register.entitiesBuffers.CleanUp()
	register.projBuffer.Release()

	register.mesh.Release()
	gl.DeleteTextures(1, &register.texture)
}

func (register materialWorldRegister) Render(p program.Program) {
	register.mutex.Lock()
	register.entitiesBuffers.Flush()
	register.projBuffer.Flush()
	register.mutex.Unlock()

	p.Use()
	register.mesh.Use()
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, register.texture)
	cmds := register.entitiesBuffers.cmdBuffer
	gl.BindBuffer(gl.DRAW_INDIRECT_BUFFER, cmds.ID())
	gl.MultiDrawElementsIndirect(gl.TRIANGLES, gl.UNSIGNED_INT, nil, int32(len(cmds.Get())), 0)
}

func createRegister(
	world ecs.World,
	p program.Program,
	assetsStorage assets.AssetsStorage,
) (materialWorldRegister, error) {
	worldMeshesAndTextures, err := ecs.GetRegister[WorldMeshesAndTextures](world)
	if err != nil {
		return materialWorldRegister{}, err
	}

	// create textures
	textures := map[assets.AssetID]int32{}
	var texture uint32
	{
		images := make([]image.Image, len(worldMeshesAndTextures.Textures))
		w, h := 0, 0
		for i, assetID := range worldMeshesAndTextures.Textures {
			textureAsset, err := assets.StorageGet[texturecomponent.TextureStorageAsset](assetsStorage, assetID)
			if err != nil {
				continue
			}

			image := textureAsset.Image()
			if i == 0 {
				w, h = image.Bounds().Dx(), image.Bounds().Dy()
			}
			if w != image.Bounds().Dx() || h != image.Bounds().Dy() {
				return materialWorldRegister{}, ErrTexturesHaveToShareSize
			}
			images[i] = image
			textures[assetID] = int32(i)
		}
		texture = CreateTexs(w, h, images)
		texLoc := gl.GetUniformLocation(p.ID(), gl.Str("texs\x00"))
		gl.Uniform1i(texLoc, 1)
	}

	// create meshes
	meshes := map[assets.AssetID]int32{}
	var packedMesh PackedMesh[Vertex]
	{
		meshesToPack := make([]Mesh[Vertex], len(worldMeshesAndTextures.Meshes))
		for i, assetID := range worldMeshesAndTextures.Meshes {
			meshAsset, err := assets.StorageGet[meshcomponent.MeshStorageAsset[Vertex]](assetsStorage, assetID)
			if err != nil {
				continue
			}
			mesh := Mesh[Vertex]{
				meshAsset.Verticies(),
				meshAsset.Indicies(),
			}
			meshesToPack[i] = mesh
			meshes[assetID] = int32(i)
		}
		packedMesh = Pack(meshesToPack...)
	}

	register := newMaterialWorldRegistry(
		meshes,
		textures,
		packedMesh,
		map[ecs.ComponentType]int32{
			ecs.GetComponentType(projection.Ortho{}):       0,
			ecs.GetComponentType(projection.Perspective{}): 1,
		},
		texture,
	)

	return register, nil
}
