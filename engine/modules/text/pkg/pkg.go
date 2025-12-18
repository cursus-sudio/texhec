package textpkg

import (
	"engine/modules/text"
	"engine/modules/text/internal/textrenderer"
	"engine/modules/text/internal/texttool"
	"engine/services/assets"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/graphics/texturearray"
	"engine/services/graphics/vao/vbo"
	"engine/services/logger"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
	"golang.org/x/image/font/opentype"
)

type pkg struct {
	defaultFontFamily func(c ioc.Dic) text.FontFamilyComponent
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
	defaultFontFamily func(c ioc.Dic) text.FontFamilyComponent,
	defaultFontSize text.FontSizeComponent,
	// defaultOverflow text.Overflow,
	defaultBreak text.BreakComponent,
	defaultTextAlign text.TextAlignComponent,
	defaultColor text.TextColorComponent,

	usedGlyphs datastructures.SparseSet[rune],
	// faceOptions opentype.FaceOptions,
	size float64,
	normalizedYBaseline float64,
) ioc.Pkg {
	return pkg{
		defaultFontFamily: defaultFontFamily,
		defaultFontSize:   defaultFontSize,
		// defaultOverflow:   defaultOverflow,
		defaultBreak:     defaultBreak,
		defaultTextAlign: defaultTextAlign,
		defaultColor:     defaultColor,
		usedGlyphs:       usedGlyphs,
		faceOptions: opentype.FaceOptions{
			Size: size,
			// DPI:  72,
			DPI: 78, // arbitrary number because it works for some reason (its a little bit rounded down)
		},
		yBaseline: int(size * normalizedYBaseline),
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[text.World, text.TextTool] {
		return texttool.NewTool(ioc.Get[logger.Logger](c))
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) textrenderer.FontService {
		return textrenderer.NewFontService(
			ioc.Get[assets.Assets](c),
			pkg.usedGlyphs,
			pkg.faceOptions,
			ioc.Get[logger.Logger](c),
			int(pkg.faceOptions.Size),
			pkg.yBaseline,
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) textrenderer.LayoutServiceFactory {
		return textrenderer.NewLayoutServiceFactory(
			ioc.Get[ecs.ToolFactory[text.World, text.TextTool]](c),
			ioc.Get[logger.Logger](c),
			ioc.Get[textrenderer.FontService](c),
			ioc.Get[textrenderer.FontKeys](c),
			pkg.defaultFontFamily(c),
			pkg.defaultFontSize,
			// pkg.defaultOverflow,
			pkg.defaultBreak,
			pkg.defaultTextAlign,
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) textrenderer.FontKeys {
		return textrenderer.NewFontKeys()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) text.System {
		return textrenderer.NewTextRendererRegister(
			ioc.Get[ecs.ToolFactory[text.World, text.TextTool]](c),
			ioc.Get[textrenderer.FontService](c),
			ioc.Get[vbo.VBOFactory[textrenderer.Glyph]](c),
			ioc.Get[textrenderer.LayoutServiceFactory](c),
			ioc.Get[logger.Logger](c),
			pkg.defaultFontFamily(c).FontFamily,
			pkg.defaultColor,
			ioc.Get[texturearray.Factory](c),
			ioc.Get[textrenderer.FontKeys](c),
			1,
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) vbo.VBOFactory[textrenderer.Glyph] {
		return func() vbo.VBOSetter[textrenderer.Glyph] {
			vbo := vbo.NewVBO[textrenderer.Glyph](func() {
				var i uint32 = 0

				gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false,
					int32(unsafe.Sizeof(textrenderer.Glyph{})), uintptr(unsafe.Offsetof(textrenderer.Glyph{}.Pos)))
				gl.EnableVertexAttribArray(i)
				i++

				gl.VertexAttribIPointerWithOffset(i, 1, gl.INT,
					int32(unsafe.Sizeof(textrenderer.Glyph{})), uintptr(unsafe.Offsetof(textrenderer.Glyph{}.Glyph)))
				gl.EnableVertexAttribArray(i)
				i++
			})
			return vbo
		}
	})
}
