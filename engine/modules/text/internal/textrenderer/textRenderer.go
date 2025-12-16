package textrenderer

import (
	"engine/modules/camera"
	"engine/modules/groups"
	rendersys "engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/assets"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/graphics/program"
	"engine/services/graphics/texturearray"
	"engine/services/logger"

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

	world       ecs.World
	cameraArray ecs.ComponentsArray[camera.CameraComponent]
	groupsArray ecs.ComponentsArray[groups.GroupsComponent]
	text        text.Interface
	transform   transform.Interface

	logger      logger.Logger
	cameraCtors camera.Interface
	fontService FontService

	program   program.Program
	locations locations

	defaultColor text.TextColorComponent

	textureFactory texturearray.Factory

	fontKeys     FontKeys
	fontsBatches datastructures.SparseArray[FontKey, fontBatch]

	dirtyEntities  ecs.DirtySet
	layoutsBatches datastructures.SparseArray[ecs.EntityID, layoutBatch]
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

	dirtyEntities := s.dirtyEntities.Get()

	// ensure fonts exist
	if len(dirtyEntities) != 0 {
		// get used fonts
		fonts := datastructures.NewSparseArray[FontKey, assets.AssetID]()
		fonts.Set(s.fontKeys.GetKey(s.defaultTextAsset), s.defaultTextAsset)
		for _, font := range s.text.FontFamily().GetEntities() {
			family, ok := s.text.FontFamily().GetComponent(font)
			if !ok {
				continue
			}
			fonts.Set(s.fontKeys.GetKey(family.FontFamily), family.FontFamily)
		}

		// remove unused fonts
		for _, key := range s.fontsBatches.GetIndices() {
			if _, ok := fonts.Get(key); ok {
				continue
			}
			batch, ok := s.fontsBatches.Get(key)
			if !ok {
				continue
			}
			fonts.Remove(key)
			batch.Release()
			s.fontsBatches.Remove(key)
		}

		// add freshly added fonts
		for _, value := range fonts.GetValues() {
			s.logger.Warn(s.ensureFontExists(value))
		}
	}

	// ensure layouts exist
	// add batches
	for _, entity := range dirtyEntities {
		if prevBatch, ok := s.layoutsBatches.Get(entity); ok {
			prevBatch.Release()
			s.layoutsBatches.Remove(entity)
		}

		layout, err := s.layoutServiceFactory.New(s.world).EntityLayout(entity)
		if err != nil {
			continue
		}

		batch := NewLayoutBatch(s.vboFactory, layout)
		s.layoutsBatches.Set(entity, batch)
	}

	// render layouts
	for _, entity := range s.layoutsBatches.GetIndices() {
		layout, _ := s.layoutsBatches.Get(entity)
		font, ok := s.fontsBatches.Get(layout.Layout.Font)
		if !ok {
			if prevBatch, ok := s.layoutsBatches.Get(entity); ok {
				prevBatch.Release()
				s.layoutsBatches.Remove(entity)
			}
			continue
		}

		pos, ok := s.transform.AbsolutePos().GetComponent(entity)
		if !ok {
			continue
		}
		rot, ok := s.transform.AbsoluteRotation().GetComponent(entity)
		if !ok {
			continue
		}
		size, ok := s.transform.AbsoluteSize().GetComponent(entity)
		if !ok {
			continue
		}
		entityColor, ok := s.text.TextColor().GetComponent(entity)
		if !ok {
			entityColor = s.defaultColor
		}

		entityGroups, ok := s.groupsArray.GetComponent(entity)
		if !ok {
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
			size.Size.Z()/2,
		)
		entityMvp := translation.Mul4(rotation).Mul4(scale)

		for _, cameraEntity := range s.cameraArray.GetEntities() {
			camera, err := s.cameraCtors.GetObject(cameraEntity)
			if err != nil {
				continue
			}

			cameraGroups, ok := s.groupsArray.GetComponent(cameraEntity)
			if !ok {
				cameraGroups = groups.DefaultGroups()
			}

			if !cameraGroups.SharesAnyGroup(entityGroups) {
				continue
			}

			mvp := camera.Mat4().Mul4(entityMvp)
			gl.UniformMatrix4fv(s.locations.Mvp, 1, false, &mvp[0])
			gl.Uniform4fv(s.locations.Color, 1, &entityColor.Color[0])
			gl.Viewport(camera.Viewport())

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
