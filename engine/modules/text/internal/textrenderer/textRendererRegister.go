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
	eventsBuilder       events.Builder
	world               ecs.World
	camera              camera.Service
	groups              groups.Service
	transform           transform.Service
	text                text.Service
	fontService         FontService
	vboFactory          vbo.VBOFactory[Glyph]
	layoutService       LayoutService
	logger              logger.Logger
	textureArrayFactory texturearray.Factory

	defaultTextAsset assets.AssetID
	defaultColor     text.TextColorComponent

	fontsKeys FontKeys

	removeOncePerNCalls uint16
}

func NewTextRenderer(
	eventsBuilder events.Builder,
	world ecs.World,
	camera camera.Service,
	groups groups.Service,
	transform transform.Service,
	text text.Service,
	fontService FontService,
	vboFactory vbo.VBOFactory[Glyph],
	layoutService LayoutService,
	logger logger.Logger,
	defaultTextAsset assets.AssetID,
	defaultColor text.TextColorComponent,
	textureArrayFactory texturearray.Factory,
	fontsKeys FontKeys,
	removeOncePerNCalls uint16,
) text.System {
	return &textRendererRegister{
		eventsBuilder:       eventsBuilder,
		world:               world,
		camera:              camera,
		groups:              groups,
		transform:           transform,
		text:                text,
		fontService:         fontService,
		vboFactory:          vboFactory,
		layoutService:       layoutService,
		logger:              logger,
		textureArrayFactory: textureArrayFactory,

		defaultTextAsset: defaultTextAsset,
		defaultColor:     defaultColor,

		fontsKeys:           fontsKeys,
		removeOncePerNCalls: removeOncePerNCalls,
	}
}

func (f *textRendererRegister) Register() error {
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

	renderer := &textRenderer{
		textRendererRegister: f,

		world:     f.world,
		groups:    f.groups,
		camera:    f.camera,
		transform: f.transform,
		text:      f.text,

		logger:      f.logger,
		fontService: f.fontService,

		program:   p,
		locations: locations,

		defaultColor: f.defaultColor,

		textureFactory: f.textureArrayFactory,

		fontKeys:     f.fontsKeys,
		fontsBatches: datastructures.NewSparseArray[FontKey, fontBatch](),

		dirtyEntities:  ecs.NewDirtySet(),
		layoutsBatches: datastructures.NewSparseArray[ecs.EntityID, layoutBatch](),
	}

	renderer.transform.AddDirtySet(renderer.dirtyEntities)

	arrays := []ecs.AnyComponentArray{
		ecs.GetComponentsArray[text.TextComponent](f.world),
		ecs.GetComponentsArray[text.BreakComponent](f.world),
		ecs.GetComponentsArray[text.FontFamilyComponent](f.world),
		// ecs.GetComponentsArray[text.Overflow](w),
		ecs.GetComponentsArray[text.FontSizeComponent](f.world),
		ecs.GetComponentsArray[text.TextAlignComponent](f.world),
	}

	for _, array := range arrays {
		array.AddDirtySet(renderer.dirtyEntities)
	}

	events.Listen(f.eventsBuilder, renderer.Listen)

	return nil
}
