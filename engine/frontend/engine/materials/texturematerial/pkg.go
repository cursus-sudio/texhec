package texturematerial

import (
	_ "embed"
	"frontend/services/assets"
	"frontend/services/media/window"

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

var TextureMaterial assets.AssetID = "materials/texturematerial"

func (Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b assets.AssetsStorageBuilder) assets.AssetsStorageBuilder {
		b.RegisterAsset(TextureMaterial, func() (assets.StorageAsset, error) {
			material := newTextureMaterial(
				vertSource,
				fragSource,
				ioc.Get[window.Api](c),
				ioc.Get[assets.Assets](c),
				nil,
			)
			return material.Material(), nil
		})
		return b
	})
}
