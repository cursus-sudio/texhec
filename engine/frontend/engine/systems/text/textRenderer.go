package textsys

import (
	"fmt"
	"frontend/engine/components/groups"
	"frontend/engine/components/projection"
	"frontend/engine/components/transform"
	rendersys "frontend/engine/systems/render"
	"frontend/engine/tools/cameras"
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

// make a pipeline which:
// - takes these components and returns vertices (glyphs)
// - takes assets and returns (size, uv, texture in array)

type locations struct {
	Mvp    int32 `location:"mvp"`
	Height int32 `location:"height"`
}

type TextRenderer struct {
	world          ecs.World
	groupsArray    ecs.ComponentsArray[groups.Groups]
	transformArray ecs.ComponentsArray[transform.Transform]
	cameraQuery    ecs.LiveQuery

	logger      logger.Logger
	cameraCtors cameras.CameraConstructors
	fontService FontService

	program   program.Program
	locations locations

	textureFactory texturearray.Factory

	// ensure font key service exists with available keys tracker
	fontKeys     FontKeys
	fontsBatches datastructures.SparseArray[FontKey, fontBatch]

	layoutsBatches datastructures.SparseArray[ecs.EntityID, layoutBatch]
}

// this may be used in textRendererFactory.go

func (s *TextRenderer) ensureOnlyFontsExist(assets []assets.AssetID) error {
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

func (s *TextRenderer) ensureFontKeyExists(key FontKey) error {
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
	batch, err := newFontBatch(s.textureFactory, font)
	if err != nil {
		return err
	}
	s.fontsBatches.Set(key, batch)
	return nil
}

func (s *TextRenderer) ensureFontExists(asset assets.AssetID) error {
	key := s.fontKeys.GetKey(asset)
	if batch, ok := s.fontsBatches.Get(key); ok {
		batch.Release()
		s.fontsBatches.Remove(key)
	}

	font, err := s.fontService.AssetFont(asset)
	if err != nil {
		return err
	}
	batch, err := newFontBatch(s.textureFactory, font)
	if err != nil {
		return err
	}
	s.fontsBatches.Set(key, batch)
	return nil
}

func (s *TextRenderer) Listen(rendersys.RenderEvent) {
	s.program.Use()

	// s.logger.Info(fmt.Sprintf(
	// 	"entities length during render is %v; ptr is %p\n",
	// 	len(s.layoutsBatches.GetIndices()),
	// 	s.layoutsBatches,
	// ))

	for _, entity := range s.layoutsBatches.GetIndices() {
		layout, _ := s.layoutsBatches.Get(entity)
		font, ok := s.fontsBatches.Get(layout.Layout.Font)
		if !ok {
			s.logger.Info(fmt.Sprintf("fonts len is %v\n", len(s.fontsBatches.GetIndices())))
			s.layoutsBatches.Remove(entity)
			continue
		}

		entityTransform, err := s.transformArray.GetComponent(entity)
		if err != nil {
			entityTransform = transform.NewTransform()
		}
		entityTransform.SetSize(mgl32.Vec3{float32(layout.Layout.FontSize), float32(layout.Layout.FontSize), 1})

		entityGroups, err := s.groupsArray.GetComponent(entity)
		if err != nil {
			entityGroups = groups.DefaultGroups()
		}

		// apply changes on batch
		font.textures.Use()
		gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, font.glyphsWidth.ID())
		layout.vao.Use()

		// s.logger.Info(fmt.Sprintf("glyphs are %v with widths %v\n", layout.Layout.Glyphs, font.font.GlyphsWidth.GetValues()))
		// height := entityTransform.Size.Y()
		var height float32 = 0

		for _, cameraEntity := range s.cameraQuery.Entities() {
			camera, err := s.cameraCtors.Get(cameraEntity, ecs.GetComponentType(projection.Ortho{}))
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
			entityTransform.SetSize(mgl32.Vec3{
				float32(layout.Layout.FontSize),
				float32(layout.Layout.FontSize),
				float32(layout.Layout.FontSize),
			})

			mvp := camera.Mat4().Mul4(entityTransform.Mat4())
			gl.UniformMatrix4fv(s.locations.Mvp, 1, false, &mvp[0])
			gl.Uniform1f(s.locations.Height, height)

			gl.Enable(gl.BLEND)
			gl.DepthMask(false)
			gl.DrawArrays(gl.POINTS, 0, layout.verticesCount)
			gl.Disable(gl.BLEND)
			gl.DepthMask(true)
		}
	}
}

func (s TextRenderer) Release() {
	for _, batch := range s.fontsBatches.GetValues() {
		batch.Release()
	}

	for _, batch := range s.layoutsBatches.GetValues() {
		batch.Release()
	}
}
