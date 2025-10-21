package textsys

import (
	"frontend/engine/components/text"
	"frontend/engine/tools/cameras"
	"frontend/services/assets"
	"frontend/services/graphics/texturearray"
	"frontend/services/graphics/vao/vbo"
	"shared/services/datastructures"
	"shared/services/logger"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
	"golang.org/x/image/font/opentype"
)

type Pkg struct {
	defaultFontFamily text.FontFamily
	defaultFontSize   text.FontSize
	// defaultOverflow   text.Overflow
	defaultBreak     text.Break
	defaultTextAlign text.TextAlign

	usedGlyphs  datastructures.SparseSet[rune]
	faceOptions opentype.FaceOptions
	yBaseline   int
}

func Package(
	defaultFontFamily text.FontFamily,
	defaultFontSize text.FontSize,
	// defaultOverflow text.Overflow,
	defaultBreak text.Break,
	defaultTextAlign text.TextAlign,

	usedGlyphs datastructures.SparseSet[rune],
	faceOptions opentype.FaceOptions,
	yBaseline int,
) ioc.Pkg {
	return Pkg{
		defaultFontFamily: defaultFontFamily,
		defaultFontSize:   defaultFontSize,
		// defaultOverflow:   defaultOverflow,
		defaultBreak:     defaultBreak,
		defaultTextAlign: defaultTextAlign,
		usedGlyphs:       usedGlyphs,
		faceOptions:      faceOptions,
		yBaseline:        yBaseline,
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterTransient(b, func(c ioc.Dic) FontService {
		return newFontService(
			ioc.Get[assets.Assets](c),
			pkg.usedGlyphs,
			pkg.faceOptions,
			ioc.Get[logger.Logger](c),
			int(pkg.faceOptions.Size),
			pkg.yBaseline,
		)
	})

	ioc.RegisterTransient(b, func(c ioc.Dic) LayoutServiceFactory {
		return newLayoutServiceFactory(
			ioc.Get[logger.Logger](c),
			ioc.Get[FontService](c),
			ioc.Get[FontKeys](c),
			pkg.defaultFontFamily,
			pkg.defaultFontSize,
			// pkg.defaultOverflow,
			pkg.defaultBreak,
			pkg.defaultTextAlign,
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) FontKeys {
		return newFontKeys()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) TextRendererFactory {
		return newTextRendererFactory(
			ioc.Get[cameras.CameraConstructors](c),
			ioc.Get[FontService](c),
			ioc.Get[vbo.VBOFactory[Glyph]](c),
			ioc.Get[LayoutServiceFactory](c),
			ioc.Get[logger.Logger](c),
			pkg.defaultFontFamily.FontAsset,
			ioc.Get[texturearray.Factory](c),
			ioc.Get[FontKeys](c),
			1,
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) vbo.VBOFactory[Glyph] {
		return func() vbo.VBOSetter[Glyph] {
			vbo := vbo.NewVBO[Glyph](func() {
				var i uint32 = 0

				gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false,
					int32(unsafe.Sizeof(Glyph{})), uintptr(unsafe.Offsetof(Glyph{}.Pos)))
				gl.EnableVertexAttribArray(i)
				i++

				gl.VertexAttribIPointerWithOffset(i, 1, gl.INT,
					int32(unsafe.Sizeof(Glyph{})), uintptr(unsafe.Offsetof(Glyph{}.Glyph)))
				gl.EnableVertexAttribArray(i)
				i++
			})
			return vbo
		}
	})
}
