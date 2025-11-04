package textpkg

import (
	"frontend/modules/camera"
	"frontend/modules/text"
	"frontend/modules/text/internal"
	"frontend/services/assets"
	"frontend/services/graphics/texturearray"
	"frontend/services/graphics/vao/vbo"
	"shared/services/datastructures"
	"shared/services/ecs"
	"shared/services/logger"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
	"golang.org/x/image/font/opentype"
)

type pkg struct {
	defaultFontFamily text.FontFamilyComponent
	defaultFontSize   text.FontSizeComponent
	// defaultOverflow   text.Overflow
	defaultBreak     text.BreakComponent
	defaultTextAlign text.TextAlignComponent
	defaultColor     text.TextColorComponent

	usedGlyphs  datastructures.SparseSet[rune]
	faceOptions opentype.FaceOptions
	yBaseline   int
}

func Package(
	defaultFontFamily text.FontFamilyComponent,
	defaultFontSize text.FontSizeComponent,
	// defaultOverflow text.Overflow,
	defaultBreak text.BreakComponent,
	defaultTextAlign text.TextAlignComponent,
	defaultColor text.TextColorComponent,

	usedGlyphs datastructures.SparseSet[rune],
	faceOptions opentype.FaceOptions,
	yBaseline int,
) ioc.Pkg {
	return pkg{
		defaultFontFamily: defaultFontFamily,
		defaultFontSize:   defaultFontSize,
		// defaultOverflow:   defaultOverflow,
		defaultBreak:     defaultBreak,
		defaultTextAlign: defaultTextAlign,
		defaultColor:     defaultColor,
		usedGlyphs:       usedGlyphs,
		faceOptions:      faceOptions,
		yBaseline:        yBaseline,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) internal.FontService {
		return internal.NewFontService(
			ioc.Get[assets.Assets](c),
			pkg.usedGlyphs,
			pkg.faceOptions,
			ioc.Get[logger.Logger](c),
			int(pkg.faceOptions.Size),
			pkg.yBaseline,
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) internal.LayoutServiceFactory {
		return internal.NewLayoutServiceFactory(
			ioc.Get[logger.Logger](c),
			ioc.Get[internal.FontService](c),
			ioc.Get[internal.FontKeys](c),
			pkg.defaultFontFamily,
			pkg.defaultFontSize,
			// pkg.defaultOverflow,
			pkg.defaultBreak,
			pkg.defaultTextAlign,
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) internal.FontKeys {
		return internal.NewFontKeys()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) text.System {
		return internal.NewTextRendererRegister(
			ioc.Get[ecs.ToolFactory[camera.CameraTool]](c),
			ioc.Get[internal.FontService](c),
			ioc.Get[vbo.VBOFactory[internal.Glyph]](c),
			ioc.Get[internal.LayoutServiceFactory](c),
			ioc.Get[logger.Logger](c),
			pkg.defaultFontFamily.FontFamily,
			pkg.defaultColor,
			ioc.Get[texturearray.Factory](c),
			ioc.Get[internal.FontKeys](c),
			1,
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) vbo.VBOFactory[internal.Glyph] {
		return func() vbo.VBOSetter[internal.Glyph] {
			vbo := vbo.NewVBO[internal.Glyph](func() {
				var i uint32 = 0

				gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false,
					int32(unsafe.Sizeof(internal.Glyph{})), uintptr(unsafe.Offsetof(internal.Glyph{}.Pos)))
				gl.EnableVertexAttribArray(i)
				i++

				gl.VertexAttribIPointerWithOffset(i, 1, gl.INT,
					int32(unsafe.Sizeof(internal.Glyph{})), uintptr(unsafe.Offsetof(internal.Glyph{}.Glyph)))
				gl.EnableVertexAttribArray(i)
				i++
			})
			return vbo
		}
	})
}
