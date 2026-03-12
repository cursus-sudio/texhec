package definitionspkg

import (
	"core/modules/definitions"
	"engine/modules/assets"
	"engine/modules/collider"
	"engine/modules/registry"
	"engine/modules/render"
	"engine/modules/transition"
	"engine/services/ecs"
	"engine/services/graphics/vao/ebo"
	"engine/services/logger"
	"image"
	"image/color"
	_ "image/png"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	// register specific files
	ioc.WrapService(b, func(c ioc.Dic, b assets.Service) {
		b.Register("blank texture", func(_ assets.PathComponent) (assets.Asset, error) {
			img := image.NewRGBA(image.Rect(0, 0, 1, 1))
			white := color.RGBA{255, 255, 255, 255}
			img.Set(0, 0, white)
			asset, err := render.NewTextureAsset(img)
			return asset, err
		})
		b.Register("square mesh", func(_ assets.PathComponent) (assets.Asset, error) {
			vertices := []render.Vertex{
				{Pos: [3]float32{1, 1, 1}, TexturePos: [2]float32{1, 1}},
				{Pos: [3]float32{1, -1, 1}, TexturePos: [2]float32{1, 0}},
				{Pos: [3]float32{-1, -1, 1}, TexturePos: [2]float32{0, 0}},
				{Pos: [3]float32{-1, 1, 1}, TexturePos: [2]float32{0, 1}},
			}

			indices := []ebo.Index{
				0, 1, 2,
				0, 2, 3,
			}
			asset := render.NewMeshAsset(vertices, indices)
			return asset, nil
		})

		b.Register("square collider", func(_ assets.PathComponent) (assets.Asset, error) {
			asset := collider.NewColliderAsset(
				[]collider.AABB{collider.NewAABB(mgl32.Vec3{-1, -1}, mgl32.Vec3{1, 1})},
				[]collider.Range{collider.NewRange(collider.Leaf, 0, 2)},
				[]collider.Polygon{
					collider.NewPolygon(mgl32.Vec3{-1, -1, 0}, mgl32.Vec3{+1, -1, 0}, [3]float32{-1, +1, 0}),
					collider.NewPolygon(mgl32.Vec3{+1, +1, 0}, mgl32.Vec3{+1, -1, 0}, [3]float32{-1, +1, 0}),
				})
			return asset, nil
		})
	})

	// register assets
	ioc.RegisterSingleton(b, func(c ioc.Dic) definitions.Definitions {
		logger := ioc.Get[logger.Logger](c)
		registryService := ioc.Get[registry.Service](c)

		gameAssets, err := registry.GetRegistry[definitions.Definitions](registryService)
		logger.Warn(err)
		return gameAssets
	})

	//
	//
	//

	// animations

	transitions := map[string]func(t transition.Progress) transition.Progress{
		"linear": func(t transition.Progress) transition.Progress {
			return t
		},
		"my easing": func(t transition.Progress) transition.Progress {
			const n1 = 7.5625
			const d1 = 2.75

			if t < 1/d1 { // First segment of the bounce (rising curve)
				return n1 * t * t
			} else if t < 2/d1 { // Second segment (peak of the first bounce)
				t -= 1.5 / d1
				return n1*t*t + 0.75
			} else if t < 2.5/d1 { // Third segment (peak of the second, smaller bounce)
				t -= 2.25 / d1
				return n1*t*t + 0.9375
			} else { // Final segment (settling)
				t -= 2.625 / d1
				return n1*t*t + 0.984375
			}
		},
		"ease out elastic": func(t transition.Progress) transition.Progress {
			const c1 float64 = 10
			const c2 float64 = .75
			const c3 float64 = (2 * math.Pi) / 3
			if t == 0 {
				return 0
			}
			if t == 1 {
				return 1
			}
			x := float64(t)
			x = math.Pow(2, -c1*x)*
				math.Sin((x*c1-c2)*c3) +
				1
			return transition.Progress(x)
		},
	}

	ioc.WrapService(b, func(c ioc.Dic, b registry.Service) {
		b.Register("transition", func(entity ecs.EntityID, structTagValue string) {
			transitionService := ioc.Get[transition.Service](c)
			easing, ok := transitions[structTagValue]
			if !ok {
				easing = func(t transition.Progress) transition.Progress { return t }
			}
			transitionService.EasingFunction().Set(entity, transition.NewEasingFunction(easing))
		})
	})
}
