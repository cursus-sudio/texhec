package textrenderer

import (
	"engine/modules/groups"
	rendersys "engine/modules/render"
	"engine/modules/text"
	"engine/services/assets"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/graphics/program"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type locations struct {
	Mvp    int32 `uniform:"mvp"`
	Color  int32 `uniform:"u_color"`
	Offset int32 `uniform:"offset"`
}

type textRenderer struct {
	*textRendererRegister

	program   program.Program
	locations locations

	defaultColor text.TextColorComponent

	fontsBatches datastructures.SparseArray[FontKey, fontBatch]

	dirtyEntities  ecs.DirtySet
	layoutsBatches datastructures.SparseArray[ecs.EntityID, layoutBatch]
}

func (s *textRenderer) ensureFontExists(asset assets.AssetID) error {
	key := s.FontsKeys.GetKey(asset)
	if batch, ok := s.fontsBatches.Get(key); ok {
		batch.Release()
		s.fontsBatches.Remove(key)
	}

	font, err := s.FontService.AssetFont(asset)
	if err != nil {
		return err
	}
	batch, err := NewFontBatch(s.TextureArrayFactory, font)
	if err != nil {
		return err
	}
	s.fontsBatches.Set(key, batch)
	return nil
}

func (s *textRenderer) ListenRender(rendersys.RenderEvent) {
	if dirtyEntities := s.dirtyEntities.Get(); len(dirtyEntities) != 0 {
		// ensure fonts exist
		// get used fonts
		fonts := datastructures.NewSparseArray[FontKey, assets.AssetID]()
		fonts.Set(s.FontsKeys.GetKey(s.defaultTextAsset), s.defaultTextAsset)
		for _, font := range s.Text.FontFamily().GetEntities() {
			family, ok := s.Text.FontFamily().Get(font)
			if !ok {
				continue
			}
			fonts.Set(s.FontsKeys.GetKey(family.FontFamily), family.FontFamily)
		}

		// we don't remove unused fonts so i'll leave this commented
		// remove unused fonts
		// for _, key := range s.fontsBatches.GetIndices() {
		// 	if _, ok := fonts.Get(key); ok {
		// 		continue
		// 	}
		// 	batch, ok := s.fontsBatches.Get(key)
		// 	if !ok {
		// 		continue
		// 	}
		// 	fonts.Remove(key)
		// 	batch.Release()
		// 	s.fontsBatches.Remove(key)
		// }

		// add freshly added fonts
		for _, value := range fonts.GetValues() {
			s.Logger.Warn(s.ensureFontExists(value))
		}

		//

		// ensure layouts exist
		// add batches
		for _, entity := range dirtyEntities {
			if prevBatch, ok := s.layoutsBatches.Get(entity); ok {
				prevBatch.Release()
				s.layoutsBatches.Remove(entity)
			}

			layout, err := s.LayoutService.EntityLayout(entity)
			if err != nil {
				continue
			}

			batch := NewLayoutBatch(s.VboFactory, layout)
			s.layoutsBatches.Set(entity, batch)
		}
	}

	// render layouts
	s.program.Bind()
	for _, cameraEntity := range s.Camera.Component().GetEntities() {
		cameraGroups, _ := s.Groups.Component().Get(cameraEntity)
		cameraMatrix := s.Camera.Mat4(cameraEntity)
		gl.Viewport(s.Camera.GetViewport(cameraEntity))

		for _, entity := range s.layoutsBatches.GetIndices() {
			entityColor, ok := s.Text.Color().Get(entity)
			if !ok {
				entityColor = s.defaultColor
			}

			entityGroups, ok := s.Groups.Component().Get(entity)
			if !ok {
				entityGroups = groups.DefaultGroups()
			}
			if !cameraGroups.SharesAnyGroup(entityGroups) {
				continue
			}

			layout, _ := s.layoutsBatches.Get(entity)
			font, ok := s.fontsBatches.Get(layout.Layout.Font)
			if !ok {
				if prevBatch, ok := s.layoutsBatches.Get(entity); ok {
					prevBatch.Release()
					s.layoutsBatches.Remove(entity)
				}
				continue
			}

			pos, _ := s.Transform.AbsolutePos().Get(entity)
			rot, _ := s.Transform.AbsoluteRotation().Get(entity)
			size, _ := s.Transform.AbsoluteSize().Get(entity)

			{
				translation := mgl32.Translate3D(pos.Pos.Elem())
				rotation := rot.Rotation.Mat4()
				scale := mgl32.Scale3D(
					float32(layout.Layout.FontSize),
					float32(layout.Layout.FontSize),
					size.Size.Z()/2,
				)
				entityModel := translation.Mul4(rotation).Mul4(scale)
				mvp := cameraMatrix.Mul4(entityModel)
				gl.UniformMatrix4fv(s.locations.Mvp, 1, false, &mvp[0])
			}
			gl.Uniform4fv(s.locations.Color, 1, &entityColor.Color[0])
			{
				offset := mgl32.Vec2{
					(-size.Size.X() / 2) / float32(layout.Layout.FontSize),
					(size.Size.Y()/2 - float32(layout.Layout.FontSize)) / float32(layout.Layout.FontSize),
				}
				gl.Uniform2f(s.locations.Offset, offset.X(), offset.Y())
			}

			// apply changes on batch
			font.textures.Bind()
			font.glyphsWidth.Bind()
			layout.vao.Bind()

			gl.DrawArrays(gl.POINTS, 0, layout.verticesCount)
		}
	}
}
