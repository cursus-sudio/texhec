package internal

import (
	"frontend/engine/camera"
	"frontend/engine/groups"
	rendersys "frontend/engine/render"
	"frontend/engine/transform"
	"frontend/services/assets"
	"frontend/services/graphics/program"
	"frontend/services/graphics/texturearray"
	"shared/services/datastructures"
	"shared/services/ecs"
	"shared/services/logger"
	"shared/utils/httperrors"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type locations struct {
	Mvp    int32 `uniform:"mvp"`
	Offset int32 `uniform:"offset"`
}

type textRenderer struct {
	world          ecs.World
	groupsArray    ecs.ComponentsArray[groups.Groups]
	transformArray ecs.ComponentsArray[transform.Transform]
	cameraQuery    ecs.LiveQuery

	logger      logger.Logger
	cameraCtors camera.CameraTool
	fontService FontService

	program   program.Program
	locations locations

	textureFactory texturearray.Factory

	fontKeys     FontKeys
	fontsBatches datastructures.SparseArray[FontKey, fontBatch]

	layoutsBatches datastructures.SparseArray[ecs.EntityID, layoutBatch]
}

func (s *textRenderer) ensureOnlyFontsExist(assets []assets.AssetID) error {
	wantedKeys := datastructures.NewSparseSet[FontKey]()
	for _, asset := range assets {
		wantedKeys.Add(s.fontKeys.GetKey(asset))
	}
	existingKeys := datastructures.NewSparseSet[FontKey]()
	for _, key := range s.fontsBatches.GetIndices() {
		existingKeys.Add(key)
	}

	notUsedkeys := datastructures.NewSparseSet[FontKey]()
	for _, existingKey := range existingKeys.GetIndices() {
		isWanted := wantedKeys.Get(existingKey)
		if isWanted {
			continue
		}
		notUsedkeys.Add(existingKey)
	}

	keysToAdd := datastructures.NewSparseSet[FontKey]()
	for _, wantedKey := range wantedKeys.GetIndices() {
		exists := existingKeys.Get(wantedKey)
		if exists {
			continue
		}
		keysToAdd.Add(wantedKey)
	}

	for _, key := range notUsedkeys.GetIndices() {
		batch, _ := s.fontsBatches.Get(key)
		batch.Release()
		s.fontsBatches.Remove(key)

	}

	for _, key := range wantedKeys.GetIndices() {
		if err := s.ensureFontKeyExists(key); err != nil {
			return err
		}
	}

	return nil
}

func (s *textRenderer) ensureFontKeyExists(key FontKey) error {
	asset, ok := s.fontKeys.GetAsset(key)
	if !ok {
		return httperrors.Err500
	}
	if batch, ok := s.fontsBatches.Get(key); ok {
		batch.Release()
		s.fontsBatches.Remove(key)
	}

	font, err := s.fontService.AssetFont(asset)
	if err != nil {
		return err
	}
	batch, err := NewFontBatch(s.textureFactory, font)
	if err != nil {
		return err
	}
	s.fontsBatches.Set(key, batch)
	return nil
}

func (s *textRenderer) ensureFontExists(asset assets.AssetID) error {
	key := s.fontKeys.GetKey(asset)
	if batch, ok := s.fontsBatches.Get(key); ok {
		batch.Release()
		s.fontsBatches.Remove(key)
	}

	font, err := s.fontService.AssetFont(asset)
	if err != nil {
		return err
	}
	batch, err := NewFontBatch(s.textureFactory, font)
	if err != nil {
		return err
	}
	s.fontsBatches.Set(key, batch)
	return nil
}

func (s *textRenderer) Listen(rendersys.RenderEvent) {
	s.program.Use()

	// render

	for _, entity := range s.layoutsBatches.GetIndices() {
		layout, _ := s.layoutsBatches.Get(entity)
		font, ok := s.fontsBatches.Get(layout.Layout.Font)
		if !ok {
			s.layoutsBatches.Remove(entity)
			continue
		}

		entityTransform, err := s.transformArray.GetComponent(entity)
		if err != nil {
			entityTransform = transform.NewTransform()
		}

		entityGroups, err := s.groupsArray.GetComponent(entity)
		if err != nil {
			entityGroups = groups.DefaultGroups()
		}

		// apply changes on batch
		font.textures.Use()
		gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, font.glyphsWidth.ID())
		layout.vao.Use()

		{
			offset := mgl32.Vec2{
				(-entityTransform.Size.X() / 2) / float32(layout.Layout.FontSize),
				(entityTransform.Size.Y()/2 - float32(layout.Layout.FontSize)) / float32(layout.Layout.FontSize),
			}
			gl.Uniform2f(s.locations.Offset, offset.X(), offset.Y())
		}

		entityTransform.SetSize(mgl32.Vec3{
			float32(layout.Layout.FontSize) * 2,
			float32(layout.Layout.FontSize) * 2,
			float32(layout.Layout.FontSize) * 2,
		})

		for _, cameraEntity := range s.cameraQuery.Entities() {
			camera, err := s.cameraCtors.Get(cameraEntity)
			if err != nil {
				continue
			}

			cameraGroups, err := s.groupsArray.GetComponent(cameraEntity)
			if err != nil {
				cameraGroups = groups.DefaultGroups()
			}

			if !cameraGroups.SharesAnyGroup(entityGroups) {
				continue
			}

			mvp := camera.Mat4().Mul4(entityTransform.Mat4())
			gl.UniformMatrix4fv(s.locations.Mvp, 1, false, &mvp[0])

			gl.DrawArrays(gl.POINTS, 0, layout.verticesCount)
		}
	}
}

func (s textRenderer) Release() {
	for _, batch := range s.fontsBatches.GetValues() {
		batch.Release()
	}

	for _, batch := range s.layoutsBatches.GetValues() {
		batch.Release()
	}
}
