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
	cameraCtorsFactory   ecs.ToolFactory[camera.CameraTool]
	transformToolFactory ecs.ToolFactory[transform.TransformTool]
	textToolFactory      ecs.ToolFactory[text.TextTool]
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
	cameraCtorsFactory ecs.ToolFactory[camera.CameraTool],
	transformToolFactory ecs.ToolFactory[transform.TransformTool],
	textToolFactory ecs.ToolFactory[text.TextTool],
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
		textToolFactory:      textToolFactory,
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

	transformTool := f.transformToolFactory.Build(w).Transform()
	renderer := textRenderer{
		textRendererRegister: f,

		world:       w,
		cameraArray: ecs.GetComponentsArray[camera.CameraComponent](w),
		groupsArray: ecs.GetComponentsArray[groups.GroupsComponent](w),
		text:        f.textToolFactory.Build(w).Text(),
		transform:   transformTool,

		logger:      f.logger,
		cameraCtors: f.cameraCtorsFactory.Build(w).Camera(),
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

	transformTool.AddDirtySet(renderer.dirtyEntities)

	arrays := []ecs.AnyComponentArray{
		ecs.GetComponentsArray[text.TextComponent](w),
		ecs.GetComponentsArray[text.BreakComponent](w),
		ecs.GetComponentsArray[text.FontFamilyComponent](w),
		// ecs.GetComponentsArray[text.Overflow](w),
		ecs.GetComponentsArray[text.FontSizeComponent](w),
		ecs.GetComponentsArray[text.TextAlignComponent](w),
	}

	for _, array := range arrays {
		array.AddDirtySet(renderer.dirtyEntities)
	}

	w.SaveGlobal(renderer)

	events.Listen(w.EventsBuilder(), renderer.Listen)

	return nil
}
