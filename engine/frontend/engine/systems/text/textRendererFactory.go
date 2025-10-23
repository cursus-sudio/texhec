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
)

//go:embed shader.vert
var vertSource string

//go:embed shader.geom
var geomSource string

//go:embed shader.frag
var fragSource string

type TextRendererFactory interface {
	New(ecs.World) (ecs.SystemRegister, error)
}

type textRendererFactory struct {
	cameraCtors          cameras.CameraConstructors
	fontService          FontService
	vboFactory           vbo.VBOFactory[Glyph]
	layoutServiceFactory LayoutServiceFactory
	logger               logger.Logger
	defaultTextAsset     assets.AssetID
	textureArrayFactory  texturearray.Factory

	fontsKeys FontKeys

	removeOncePerNCalls uint16
}

func newTextRendererFactory(
	cameraCtors cameras.CameraConstructors,
	fontService FontService,
	vboFactory vbo.VBOFactory[Glyph],
	layoutServiceFactory LayoutServiceFactory,
	logger logger.Logger,
	defaultTextAsset assets.AssetID,
	textureArrayFactory texturearray.Factory,
	fontsKeys FontKeys,
	removeOncePerNCalls uint16,
) TextRendererFactory {
	return &textRendererFactory{
		cameraCtors:          cameraCtors,
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

func (f *textRendererFactory) New(world ecs.World) (ecs.SystemRegister, error) {
	vert, err := shader.NewShader(vertSource, shader.VertexShader)
	if err != nil {
		return nil, err
	}
	defer vert.Release()

	geom, err := shader.NewShader(geomSource, shader.GeomShader)
	if err != nil {
		return nil, err
	}
	defer geom.Release()

	frag, err := shader.NewShader(fragSource, shader.FragmentShader)
	if err != nil {
		return nil, err
	}
	defer frag.Release()

	programID := gl.CreateProgram()
	gl.AttachShader(programID, vert.ID())
	gl.AttachShader(programID, geom.ID())
	gl.AttachShader(programID, frag.ID())

	p, err := program.NewProgram(programID, nil)
	if err != nil {
		return nil, err
	}

	locations, err := program.GetProgramLocations[locations](p)
	if err != nil {
		p.Release()
		return nil, err
	}

	renderer := textRenderer{
		world:          world,
		transformArray: ecs.GetComponentsArray[transform.Transform](world.Components()),
		groupsArray:    ecs.GetComponentsArray[groups.Groups](world.Components()),
		cameraQuery:    world.QueryEntitiesWithComponents(ecs.GetComponentType(projection.Ortho{})),

		logger:      f.logger,
		cameraCtors: f.cameraCtors,
		fontService: f.fontService,

		program:   p,
		locations: locations,

		textureFactory: f.textureArrayFactory,

		fontKeys:     f.fontsKeys,
		fontsBatches: datastructures.NewSparseArray[FontKey, fontBatch](),

		layoutsBatches: datastructures.NewSparseArray[ecs.EntityID, layoutBatch](),
	}

	query := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(text.Text{}),
		ecs.GetComponentType(transform.Transform{}),
	)

	addOrChangeListener := func(ei []ecs.EntityID) {
		for _, entity := range ei {
			if prevBatch, ok := renderer.layoutsBatches.Get(entity); ok {
				prevBatch.Release()
				renderer.layoutsBatches.Remove(entity)
			}

			layout, err := f.layoutServiceFactory.New(world).EntityLayout(entity)
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
		ecs.GetComponentsArray[text.Break](world.Components()),
		ecs.GetComponentsArray[text.FontFamily](world.Components()),
		// ecs.GetComponentsArray[text.Overflow](world.Components()),
		ecs.GetComponentsArray[text.FontSize](world.Components()),
		ecs.GetComponentsArray[text.TextAlign](world.Components()),
	}

	for _, array := range arrays {
		array.OnAdd(addOrChangeListener)
		array.OnChange(addOrChangeListener)
		array.OnRemove(rmListener)
	}

	fontArray := ecs.GetComponentsArray[text.FontFamily](world.Components())
	addFonts := func(ei []ecs.EntityID) {
		for _, entity := range ei {
			family, err := fontArray.GetComponent(entity)
			if err != nil {
				continue
			}

			if err := renderer.ensureFontExists(family.FontAsset); err != nil {
				f.logger.Error(err)
			}
		}
	}
	if err := renderer.ensureFontExists(f.defaultTextAsset); err != nil {
		p.Release()
		return nil, err
	}
	fontArray.OnAdd(addFonts)
	fontArray.OnChange(addFonts)
	{
		fontFamilyArray := ecs.GetComponentsArray[text.FontFamily](world.Components())
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
				assets = append(assets, comp.FontAsset)
			}
			if err := renderer.ensureOnlyFontsExist(assets); err != nil {
				f.logger.Error(err)
			}
		}
		fontArray.OnChange(removeUnused)
		fontArray.OnRemove(removeUnused)
	}

	world.SaveRegister(renderer)

	return &renderer, nil
}
