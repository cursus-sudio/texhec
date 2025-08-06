package texturematerial

import (
	_ "embed"
	"frontend/engine/components/projection"
	"frontend/services/assets"
	"frontend/services/media/window"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

//go:embed s.frag
var fragSource string

//go:embed s.vert
var vertSource string

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

var (
	TextureMaterial2D assets.AssetID = "materials/texturematerial2d"
	TextureMaterial3D assets.AssetID = "materials/texturematerial3d"
)

type X[T ~struct{}] struct{}
type x struct{}

func (Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b assets.AssetsStorageBuilder) assets.AssetsStorageBuilder {
		b.RegisterAsset(TextureMaterial3D, func() (assets.StorageAsset, error) {
			material := newTextureMaterial(
				vertSource,
				fragSource,
				ioc.Get[window.Api](c),
				ioc.Get[assets.Assets](c),
				nil,
				func(component projection.Perspecrive) mgl32.Mat4 {
					return mgl32.Perspective(
						component.FovY,
						component.AspectRatio,
						component.Near,
						component.Far,
					)
				},
			)
			return material.Material(), nil
		})
		b.RegisterAsset(TextureMaterial2D, func() (assets.StorageAsset, error) {
			material := newTextureMaterial(
				vertSource,
				fragSource,
				ioc.Get[window.Api](c),
				ioc.Get[assets.Assets](c),
				nil,
				func(component projection.Ortho) mgl32.Mat4 {
					return mgl32.Ortho(
						-float32(component.Width)/2,
						float32(component.Width)/2,
						-float32(component.Height)/2,
						float32(component.Height)/2,
						component.Near,
						component.Far,
					)
				},
			)
			return material.Material(), nil
		})
		return b
	})
}
