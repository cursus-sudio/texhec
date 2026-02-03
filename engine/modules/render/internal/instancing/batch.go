package instancing

import (
	"engine/modules/render"
	"engine/services/assets"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/graphics/buffers"
	"engine/services/graphics/texturearray"
	"engine/services/graphics/vao"
	"engine/services/graphics/vao/ebo"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type batchKey struct {
	texture render.TextureComponent
	mesh    render.MeshComponent
}

type batch struct {
	system       *system
	VAO          vao.VAO
	TextureArray texturearray.TextureArray
	Dirty        bool

	// buffers (model, color, frame)
	Entities datastructures.Set[ecs.EntityID]
	Models   buffers.Buffer[mgl32.Mat4]
	Colors   buffers.Buffer[mgl32.Vec4]
	Frames   buffers.Buffer[int32]
	Groups   buffers.Buffer[uint32]
}

func (s *system) NewBatch(batchKey batchKey) (*batch, error) {
	// mesh
	VAO, ok := s.meshes[batchKey.mesh.ID]
	if !ok {
		meshAsset, err := assets.GetAsset[render.MeshAsset](s.Assets, batchKey.mesh.ID)
		if err != nil {
			return nil, err
		}
		VBO := s.VboFactory()
		VBO.SetVertices(meshAsset.Vertices())
		EBO := ebo.NewEBO()
		EBO.SetIndices(meshAsset.Indices())
		VAO = vao.NewVAO(VBO, EBO)
		s.meshes[batchKey.mesh.ID] = VAO
	}

	// texture
	textureArr, ok := s.textures[batchKey.texture.Asset]
	if !ok {
		textureAsset, err := assets.GetAsset[render.TextureAsset](s.Assets, batchKey.texture.Asset)
		if err != nil {
			return nil, err
		}
		textureArr, err = s.TextureArrayFactory.NewFromSlice(textureAsset.Images())
		if err != nil {
			return nil, err
		}
		s.textures[batchKey.texture.Asset] = textureArr
	}

	// batch
	batch := &batch{
		system:       s,
		VAO:          VAO,
		TextureArray: textureArr,
		Dirty:        true,

		Entities: datastructures.NewSet[ecs.EntityID](),
	}

	// buffers
	batch.Models = buffers.NewBuffer[mgl32.Mat4](gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, 0)
	batch.Colors = buffers.NewBuffer[mgl32.Vec4](gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, 1)
	batch.Frames = buffers.NewBuffer[int32](gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, 2)
	batch.Groups = buffers.NewBuffer[uint32](gl.SHADER_STORAGE_BUFFER, gl.DYNAMIC_DRAW, 3)

	return batch, nil
}

//

func (s *batch) Upsert(entity ecs.EntityID) {
	index, ok := s.Entities.GetIndex(entity)
	for !ok {
		s.Entities.Add(entity)
		index, ok = s.Entities.GetIndex(entity)
	}
	model := s.system.Transform.Mat4(entity)
	color, _ := s.system.Render.Color().Get(entity)
	textureFrame, _ := s.system.Render.TextureFrame().Get(entity)
	groups, _ := s.system.Groups.Component().Get(entity)

	frame := int32(textureFrame.GetFrame(s.TextureArray.ImagesCount))

	s.Dirty = true
	s.Models.Set(index, model)
	s.Colors.Set(index, color.Color)
	s.Frames.Set(index, frame)
	s.Groups.Set(index, groups.Mask)
}

func (s *batch) Remove(entity ecs.EntityID) {
	index, ok := s.Entities.GetIndex(entity)
	if !ok {
		return
	}

	s.Dirty = true
	s.Entities.Remove(index)
	s.Models.Remove(index)
	s.Colors.Remove(index)
	s.Frames.Remove(index)
	s.Groups.Remove(index)
}

//

func (s *batch) Render() {
	if s.Dirty {
		s.Dirty = false
		s.Models.Flush()
		s.Colors.Flush()
		s.Frames.Flush()
		s.Groups.Flush()
	}

	s.VAO.Bind()
	s.TextureArray.Bind()

	s.Models.Bind()
	s.Colors.Bind()
	s.Frames.Bind()
	s.Groups.Bind()

	gl.DrawElementsInstanced(
		gl.TRIANGLES,
		int32(s.VAO.EBO().Len()),
		gl.UNSIGNED_INT,
		nil,
		int32(len(s.Entities.Get())),
	)
}
