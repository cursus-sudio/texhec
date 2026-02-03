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
	"github.com/ogiusek/ioc/v2"
)

//go:embed shader.vert
var vertSource string

//go:embed shader.geom
var geomSource string

//go:embed shader.frag
var fragSource string

type textRendererRegister struct {
	EventsBuilder       events.Builder        `inject:"1"`
	World               ecs.World             `inject:"1"`
	Camera              camera.Service        `inject:"1"`
	Groups              groups.Service        `inject:"1"`
	Transform           transform.Service     `inject:"1"`
	Text                text.Service          `inject:"1"`
	FontService         FontService           `inject:"1"`
	VboFactory          vbo.VBOFactory[Glyph] `inject:"1"`
	LayoutService       LayoutService         `inject:"1"`
	Logger              logger.Logger         `inject:"1"`
	TextureArrayFactory texturearray.Factory  `inject:"1"`
	FontsKeys           FontKeys              `inject:"1"`

	defaultTextAsset    assets.AssetID
	defaultColor        text.TextColorComponent
	removeOncePerNCalls uint16
}

func NewTextRenderer(c ioc.Dic,
	defaultTextAsset assets.AssetID,
	defaultColor text.TextColorComponent,
	removeOncePerNCalls uint16,
) text.SystemRenderer {
	s := ioc.GetServices[*textRendererRegister](c)
	s.defaultTextAsset = defaultTextAsset
	s.defaultColor = defaultColor
	s.removeOncePerNCalls = removeOncePerNCalls
	return s
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

		program:   p,
		locations: locations,

		defaultColor: f.defaultColor,

		fontsBatches: datastructures.NewSparseArray[FontKey, fontBatch](),

		dirtyEntities:  ecs.NewDirtySet(),
		layoutsBatches: datastructures.NewSparseArray[ecs.EntityID, layoutBatch](),
	}

	renderer.Transform.AddDirtySet(renderer.dirtyEntities)

	arrays := []ecs.AnyComponentArray{
		ecs.GetComponentsArray[text.TextComponent](f.World),
		ecs.GetComponentsArray[text.BreakComponent](f.World),
		ecs.GetComponentsArray[text.FontFamilyComponent](f.World),
		// ecs.GetComponentsArray[text.Overflow](w),
		ecs.GetComponentsArray[text.FontSizeComponent](f.World),
		ecs.GetComponentsArray[text.TextAlignComponent](f.World),
	}

	for _, array := range arrays {
		array.AddDirtySet(renderer.dirtyEntities)
	}

	events.Listen(f.EventsBuilder, renderer.ListenRender)

	return nil
}
