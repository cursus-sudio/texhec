package material

import (
	"frontend/engine/components/projection"
	texturecomponent "frontend/engine/components/texture"
	"frontend/engine/tools/worldmesh"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/graphics/program"
	"image"
	"sync"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type materialWorldRegister struct {
	mutex   *sync.RWMutex
	buffers *materialBuffers

	textures map[assets.AssetID]int32
	texture  uint32

	projections map[ecs.ComponentType]int32
}

func newMaterialWorldRegistry(
	textures map[assets.AssetID]int32,
	projections map[ecs.ComponentType]int32,
	texture uint32,
) materialWorldRegister {
	return materialWorldRegister{
		mutex:       &sync.RWMutex{},
		buffers:     newMaterialBuffers(len(projections)),
		textures:    textures,
		texture:     texture,
		projections: projections,
	}
}

func (register materialWorldRegister) Release() {
	register.buffers.Release()
	gl.DeleteTextures(1, &register.texture)
}

func (register materialWorldRegister) Render(world ecs.World, p program.Program) error {
	geometry, err := ecs.GetRegister[worldmesh.WorldMeshRegister[Vertex]](world)
	if err != nil {
		return err
	}

	register.mutex.Lock()
	register.buffers.Flush()
	register.mutex.Unlock()

	p.Use()
	geometry.Mesh.Use()
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, register.texture)
	cmds := register.buffers.cmdBuffer
	gl.BindBuffer(gl.DRAW_INDIRECT_BUFFER, cmds.ID())
	gl.MultiDrawElementsIndirect(gl.TRIANGLES, gl.UNSIGNED_INT, nil, int32(len(cmds.Get())), 0)
	return nil
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

	register := newMaterialWorldRegistry(
		textures,
		map[ecs.ComponentType]int32{
			ecs.GetComponentType(projection.Ortho{}):       0,
			ecs.GetComponentType(projection.Perspective{}): 1,
		},
		texture,
	)

	return register, nil
}
