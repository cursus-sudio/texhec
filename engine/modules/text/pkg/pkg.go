package textpkg

import (
	"engine/modules/assets"
	"engine/modules/text"
	"engine/modules/text/internal/textrenderer"
	"engine/modules/text/internal/textservice"
	"engine/services/datastructures"
	"engine/services/ecs"
	gtexture "engine/services/graphics/texture"
	"engine/services/graphics/vao/vbo"
	"engine/services/logger"
	"image"
	"image/color"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
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
	ioc.RegisterSingleton(b, func(c ioc.Dic) text.Service {
		return textservice.NewService(
			ioc.Get[ecs.World](c),
			ioc.Get[logger.Logger](c),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) textrenderer.FontService {
		return textrenderer.NewFontService(
			ioc.Get[assets.Service](c),
			pkg.usedGlyphs,
			pkg.faceOptions,
			ioc.Get[logger.Logger](c),
			int(pkg.faceOptions.Size),
			pkg.yBaseline,
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) textrenderer.LayoutService {
		return textrenderer.NewLayoutService(
			c,
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

	ioc.RegisterSingleton(b, func(c ioc.Dic) text.SystemRenderer {
		return textrenderer.NewTextRenderer(
			c,
			pkg.defaultFontFamily(c).FontFamily,
			pkg.defaultColor,
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
			})
			return vbo
		}
	})

	ioc.WrapService(b, func(c ioc.Dic, b assets.Extensions) {
		getLetterImage := func(drawer font.Drawer, letter rune) *image.RGBA {
			var text = string(letter)
			textBounds, _ := drawer.BoundString(text)

			cellSize := int(pkg.faceOptions.Size)
			rect := image.Rect(0, 0, cellSize, cellSize)
			img := image.NewRGBA(rect)
			drawer.Dst = img

			dotX := fixed.I(0) - textBounds.Min.X
			dotY := fixed.I(pkg.yBaseline)

			drawer.Dot = fixed.Point26_6{
				X: dotX,
				Y: dotY,
			}

			drawer.DrawString(text)
			return img
		}
		b.Register("ttf", func(id assets.Path) (any, error) {
			source, err := os.ReadFile(string(id))
			if err != nil {
				return nil, err
			}
			rawFont, err := opentype.Parse(source)
			if err != nil {
				return nil, err
			}
			fontFace, err := opentype.NewFace(rawFont, &pkg.faceOptions)
			if err != nil {
				return nil, err
			}

			glyphs := text.Glyphs{
				GlyphsWidth: datastructures.NewSparseArray[uint32, float32](),
				Images:      datastructures.NewSparseArray[uint32, image.Image](),
			}
			for _, glyph := range pkg.usedGlyphs.GetIndices() {
				glyphID := uint32(glyph)
				_, advance, _ := fontFace.GlyphBounds(glyph)
				width := float32(advance.Ceil()) / float32(pkg.faceOptions.Size)
				glyphs.GlyphsWidth.Set(glyphID, width)

				drawer := font.Drawer{
					Src:  image.NewUniform(color.White),
					Face: fontFace,
				}
				image := getLetterImage(drawer, glyph)
				glyphs.Images.Set(glyphID, gtexture.NewImage(image).FlipV().Image())
			}

			asset := text.NewFontAsset(*rawFont, glyphs)
			return asset, nil
		})
	})
}
