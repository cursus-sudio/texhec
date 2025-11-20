package internal

import (
	"frontend/modules/camera"
	"frontend/modules/groups"
	rendersys "frontend/modules/render"
	"frontend/modules/text"
	"frontend/modules/transform"
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
	Color  int32 `uniform:"u_color"`
	Offset int32 `uniform:"offset"`
}

type textRenderer struct {
	world                ecs.World
	groupsArray          ecs.ComponentsArray[groups.GroupsComponent]
	colorArray           ecs.ComponentsArray[text.TextColorComponent]
	transformTransaction transform.TransformTransaction
	cameraQuery          ecs.LiveQuery

	logger      logger.Logger
	cameraCtors camera.CameraTool
	fontService FontService

	program   program.Program
	locations locations

	defaultColor text.TextColorComponent

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

	for _, entity := range s.layoutsBatches.GetIndices() {
		layout, _ := s.layoutsBatches.Get(entity)
		font, ok := s.fontsBatches.Get(layout.Layout.Font)
		if !ok {
			s.layoutsBatches.Remove(entity)
			continue
		}

		entityTransform := s.transformTransaction.GetEntity(entity)
		pos, err := entityTransform.AbsolutePos().Get()
		if err != nil {
			s.logger.Warn(err)
			continue
		}
		rot, err := entityTransform.AbsoluteRotation().Get()
		if err != nil {
			s.logger.Warn(err)
			continue
		}
		size, err := entityTransform.AbsoluteSize().Get()
		if err != nil {
			s.logger.Warn(err)
			continue
		}
		entityColor, err := s.colorArray.GetComponent(entity)
		if err != nil {
			entityColor = s.defaultColor
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
				(-size.Size.X() / 2) / float32(layout.Layout.FontSize),
				(size.Size.Y()/2 - float32(layout.Layout.FontSize)) / float32(layout.Layout.FontSize),
			}
			gl.Uniform2f(s.locations.Offset, offset.X(), offset.Y())
		}

		translation := mgl32.Translate3D(pos.Pos.Elem())
		rotation := rot.Rotation.Mat4()
		scale := mgl32.Scale3D(
			float32(layout.Layout.FontSize),
			float32(layout.Layout.FontSize),
			float32(layout.Layout.FontSize),
		)
		entityMvp := translation.Mul4(rotation).Mul4(scale)

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

			mvp := camera.Mat4().Mul4(entityMvp)
			gl.UniformMatrix4fv(s.locations.Mvp, 1, false, &mvp[0])
			gl.Uniform4fv(s.locations.Color, 1, &entityColor.Color[0])

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
