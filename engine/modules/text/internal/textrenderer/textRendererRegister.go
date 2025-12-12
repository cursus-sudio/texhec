package textrenderer

import (
	_ "embed"
	"engine/modules/camera"
	"engine/modules/groups"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/assets"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/graphics/program"
	"engine/services/graphics/shader"
	"engine/services/graphics/texturearray"
	"engine/services/graphics/vao/vbo"
	"engine/services/logger"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
)

//go:embed shader.vert
var vertSource string

//go:embed shader.geom
var geomSource string

//go:embed shader.frag
var fragSource string

type textRendererRegister struct {
	cameraCtorsFactory   ecs.ToolFactory[camera.Tool]
	transformToolFactory ecs.ToolFactory[transform.Tool]
	fontService          FontService
	vboFactory           vbo.VBOFactory[Glyph]
	layoutServiceFactory LayoutServiceFactory
	logger               logger.Logger
	textureArrayFactory  texturearray.Factory

	defaultTextAsset assets.AssetID
	defaultColor     text.TextColorComponent

	fontsKeys FontKeys

	removeOncePerNCalls uint16
}

func NewTextRendererRegister(
	cameraCtorsFactory ecs.ToolFactory[camera.Tool],
	transformToolFactory ecs.ToolFactory[transform.Tool],
	fontService FontService,
	vboFactory vbo.VBOFactory[Glyph],
	layoutServiceFactory LayoutServiceFactory,
	logger logger.Logger,
	defaultTextAsset assets.AssetID,
	defaultColor text.TextColorComponent,
	textureArrayFactory texturearray.Factory,
	fontsKeys FontKeys,
	removeOncePerNCalls uint16,
) text.System {
	return &textRendererRegister{
		cameraCtorsFactory:   cameraCtorsFactory,
		transformToolFactory: transformToolFactory,
		fontService:          fontService,
		vboFactory:           vboFactory,
		layoutServiceFactory: layoutServiceFactory,
		logger:               logger,
		textureArrayFactory:  textureArrayFactory,

		defaultTextAsset: defaultTextAsset,
		defaultColor:     defaultColor,

		fontsKeys:           fontsKeys,
		removeOncePerNCalls: removeOncePerNCalls,
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

	transformTool := f.transformToolFactory.Build(w)
	renderer := textRenderer{
		textRendererRegister: f,

		world:                w,
		colorArray:           ecs.GetComponentsArray[text.TextColorComponent](w),
		groupsArray:          ecs.GetComponentsArray[groups.GroupsComponent](w),
		transformTransaction: transformTool.Transaction(),
		cameraQuery:          w.Query().Require(camera.OrthoComponent{}).Build(),

		logger:      f.logger,
		cameraCtors: f.cameraCtorsFactory.Build(w),
		fontService: f.fontService,

		program:   p,
		locations: locations,

		defaultColor: f.defaultColor,

		textureFactory: f.textureArrayFactory,

		fontKeys:     f.fontsKeys,
		fontsBatches: datastructures.NewSparseArray[FontKey, fontBatch](),

		dirtyEntities:  datastructures.NewSparseSet[ecs.EntityID](),
		layoutsBatches: datastructures.NewSparseArray[ecs.EntityID, layoutBatch](),
	}

	query := transformTool.Query(w.Query()).
		Require(text.TextComponent{}).
		Build()

	onMod := func(ei []ecs.EntityID) {
		for _, entity := range ei {
			renderer.dirtyEntities.Add(entity)
		}
	}

	query.OnAdd(onMod)
	query.OnChange(onMod)
	query.OnRemove(onMod)

	arrays := []ecs.AnyComponentArray{
		ecs.GetComponentsArray[text.BreakComponent](w),
		ecs.GetComponentsArray[text.FontFamilyComponent](w),
		// ecs.GetComponentsArray[text.Overflow](w),
		ecs.GetComponentsArray[text.FontSizeComponent](w),
		ecs.GetComponentsArray[text.TextAlignComponent](w),
	}

	for _, array := range arrays {
		array.OnAdd(onMod)
		array.OnChange(onMod)
		array.OnRemove(onMod)
	}

	fontArray := ecs.GetComponentsArray[text.FontFamilyComponent](w)
	addFonts := func(ei []ecs.EntityID) {
		for _, entity := range ei {
			family, err := fontArray.GetComponent(entity)
			if err != nil {
				continue
			}

			f.logger.Warn(renderer.ensureFontExists(family.FontFamily))
		}
	}
	if err := renderer.ensureFontExists(f.defaultTextAsset); err != nil {
		p.Release()
		return err
	}
	fontArray.OnAdd(addFonts)
	fontArray.OnChange(addFonts)
	{
		fontFamilyArray := ecs.GetComponentsArray[text.FontFamilyComponent](w)
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
			f.logger.Warn(renderer.ensureOnlyFontsExist(assets))
		}
		fontArray.OnChange(removeUnused)
		fontArray.OnRemove(removeUnused)
	}

	w.SaveGlobal(renderer)

	events.Listen(w.EventsBuilder(), renderer.Listen)

	return nil
}
