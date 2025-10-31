package textsys

import (
	_ "embed"
	"frontend/engine/components/groups"
	"frontend/engine/components/projection"
	"frontend/engine/components/text"
	"frontend/engine/components/transform"
	"frontend/engine/tools/cameras"
	"frontend/services/assets"
	"frontend/services/graphics/program"
	"frontend/services/graphics/shader"
	"frontend/services/graphics/texturearray"
	"frontend/services/graphics/vao/vbo"
	"shared/services/datastructures"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
)

//go:embed shader.vert
var vertSource string

//go:embed shader.geom
var geomSource string

//go:embed shader.frag
var fragSource string

type TextRendererRegister ecs.SystemRegister
type textRendererRegister struct {
	cameraCtorsFactory   ecs.ToolFactory[cameras.CameraResolver]
	fontService          FontService
	vboFactory           vbo.VBOFactory[Glyph]
	layoutServiceFactory LayoutServiceFactory
	logger               logger.Logger
	defaultTextAsset     assets.AssetID
	textureArrayFactory  texturearray.Factory

	fontsKeys FontKeys

	removeOncePerNCalls uint16
}

func newTextRendererRegister(
	cameraCtorsFactory ecs.ToolFactory[cameras.CameraResolver],
	fontService FontService,
	vboFactory vbo.VBOFactory[Glyph],
	layoutServiceFactory LayoutServiceFactory,
	logger logger.Logger,
	defaultTextAsset assets.AssetID,
	textureArrayFactory texturearray.Factory,
	fontsKeys FontKeys,
	removeOncePerNCalls uint16,
) TextRendererRegister {
	return &textRendererRegister{
		cameraCtorsFactory:   cameraCtorsFactory,
		fontService:          fontService,
		vboFactory:           vboFactory,
		layoutServiceFactory: layoutServiceFactory,
		logger:               logger,
		defaultTextAsset:     defaultTextAsset,
		textureArrayFactory:  textureArrayFactory,
		fontsKeys:            fontsKeys,
		removeOncePerNCalls:  removeOncePerNCalls,
	}
}

func (f *textRendererRegister) Register(w ecs.World) error {
	vert, err := shader.NewShader(vertSource, shader.VertexShader)
	if err != nil {
		return err
	}
	defer vert.Release()

	geom, err := shader.NewShader(geomSource, shader.GeomShader)
	if err != nil {
		return err
	}
	defer geom.Release()

	frag, err := shader.NewShader(fragSource, shader.FragmentShader)
	if err != nil {
		return err
	}
	defer frag.Release()

	programID := gl.CreateProgram()
	gl.AttachShader(programID, vert.ID())
	gl.AttachShader(programID, geom.ID())
	gl.AttachShader(programID, frag.ID())

	p, err := program.NewProgram(programID, nil)
	if err != nil {
		return err
	}

	locations, err := program.GetProgramLocations[locations](p)
	if err != nil {
		p.Release()
		return err
	}

	renderer := textRenderer{
		world:          w,
		transformArray: ecs.GetComponentsArray[transform.Transform](w.Components()),
		groupsArray:    ecs.GetComponentsArray[groups.Groups](w.Components()),
		cameraQuery:    w.Query().Require(ecs.GetComponentType(projection.Ortho{})).Build(),

		logger:      f.logger,
		cameraCtors: f.cameraCtorsFactory.Build(w),
		fontService: f.fontService,

		program:   p,
		locations: locations,

		textureFactory: f.textureArrayFactory,

		fontKeys:     f.fontsKeys,
		fontsBatches: datastructures.NewSparseArray[FontKey, fontBatch](),

		layoutsBatches: datastructures.NewSparseArray[ecs.EntityID, layoutBatch](),
	}

	query := w.Query().Require(
		ecs.GetComponentType(text.Text{}),
		ecs.GetComponentType(transform.Transform{}),
	).Build()

	addOrChangeListener := func(ei []ecs.EntityID) {
		for _, entity := range ei {
			if prevBatch, ok := renderer.layoutsBatches.Get(entity); ok {
				prevBatch.Release()
				renderer.layoutsBatches.Remove(entity)
			}

			layout, err := f.layoutServiceFactory.New(w).EntityLayout(entity)
			if err != nil {
				continue
			}

			batch := newLayoutBatch(f.vboFactory, layout)
			renderer.layoutsBatches.Set(entity, batch)
		}
	}
	rmListener := func(ei []ecs.EntityID) {
		for _, entity := range ei {
			if prevBatch, ok := renderer.layoutsBatches.Get(entity); ok {
				prevBatch.Release()
			}
			renderer.layoutsBatches.Remove(entity)
		}
	}

	query.OnAdd(addOrChangeListener)
	query.OnChange(addOrChangeListener)
	query.OnRemove(rmListener)

	arrays := []ecs.AnyComponentArray{
		ecs.GetComponentsArray[text.Break](w.Components()),
		ecs.GetComponentsArray[text.FontFamily](w.Components()),
		// ecs.GetComponentsArray[text.Overflow](w.Components()),
		ecs.GetComponentsArray[text.FontSize](w.Components()),
		ecs.GetComponentsArray[text.TextAlign](w.Components()),
	}

	for _, array := range arrays {
		array.OnAdd(addOrChangeListener)
		array.OnChange(addOrChangeListener)
		array.OnRemove(rmListener)
	}

	fontArray := ecs.GetComponentsArray[text.FontFamily](w.Components())
	addFonts := func(ei []ecs.EntityID) {
		for _, entity := range ei {
			family, err := fontArray.GetComponent(entity)
			if err != nil {
				continue
			}

			if err := renderer.ensureFontExists(family.FontFamily); err != nil {
				f.logger.Error(err)
			}
		}
	}
	if err := renderer.ensureFontExists(f.defaultTextAsset); err != nil {
		p.Release()
		return err
	}
	fontArray.OnAdd(addFonts)
	fontArray.OnChange(addFonts)
	{
		fontFamilyArray := ecs.GetComponentsArray[text.FontFamily](w.Components())
		var i uint16 = 0
		removeUnused := func(_ []ecs.EntityID) {
			i++
			if i < f.removeOncePerNCalls {
				return
			}

			i = 0
			entities := query.Entities()
			assets := []assets.AssetID{f.defaultTextAsset}
			for _, entity := range entities {
				comp, err := fontFamilyArray.GetComponent(entity)
				if err != nil {
					continue
				}
				assets = append(assets, comp.FontFamily)
			}
			if err := renderer.ensureOnlyFontsExist(assets); err != nil {
				f.logger.Error(err)
			}
		}
		fontArray.OnChange(removeUnused)
		fontArray.OnRemove(removeUnused)
	}

	w.SaveGlobal(renderer)

	events.Listen(w.EventsBuilder(), renderer.Listen)

	return nil
}
