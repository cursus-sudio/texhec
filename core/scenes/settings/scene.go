package settingsscene

import (
	gameassets "core/assets"
	"core/modules/settings"
	gamescenes "core/scenes"
	"engine/modules/camera"
	"engine/modules/genericrenderer"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/scenes"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) LoadObjects(b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadObjects, func(c ioc.Dic, b gamescenes.SettingsBuilder) gamescenes.SettingsBuilder {
		gameassets := ioc.Get[gameassets.GameAssets](c)
		worldResolver := ioc.Get[gamescenes.WorldResolver](c)
		b.OnLoad(func(rawWorld ecs.World) {
			world := worldResolver(rawWorld)
			cameraEntity := world.NewEntity()
			world.Camera().Ortho().Set(cameraEntity, camera.NewOrtho(-1000, +1000))

			signature := world.NewEntity()
			world.Transform().Pos().Set(signature, transform.NewPos(5, 5, 0))
			world.Transform().Size().Set(signature, transform.NewSize(100, 50, 1))
			world.Transform().PivotPoint().Set(signature, transform.NewPivotPoint(0, .5, .5))
			world.Hierarchy().SetParent(signature, cameraEntity)
			world.Transform().Parent().Set(signature, transform.NewParent(transform.RelativePos))
			world.Transform().ParentPivotPoint().Set(signature, transform.NewParentPivotPoint(0, 0, .5))

			world.Text().Content().Set(signature, text.TextComponent{Text: "settings"})
			world.Text().FontSize().Set(signature, text.FontSizeComponent{FontSize: 32})
			world.Text().Break().Set(signature, text.BreakComponent{Break: text.BreakNone})

			background := world.NewEntity()
			world.Hierarchy().SetParent(background, cameraEntity)
			world.Transform().Parent().Set(background, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))
			world.Render().Mesh().Set(background, render.NewMesh(gameassets.SquareMesh))
			world.Render().Texture().Set(background, render.NewTexture(gameassets.Tiles.Ground))
			world.GenericRenderer().Pipeline().Set(background, genericrenderer.PipelineComponent{})

			buttonArea := world.NewEntity()
			world.Transform().Size().Set(buttonArea, transform.NewSize(500, 200, 2))
			world.Hierarchy().SetParent(buttonArea, cameraEntity)
			world.Transform().Parent().Set(buttonArea, transform.NewParent(transform.RelativePos))

			events.Emit(world.Events(), settings.EnterSettingsForParentEvent{Parent: buttonArea})
		})

		return b
	})
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.SettingsBuilder { return scenes.NewSceneBuilder() })
	gamescenes.AddDefaults[gamescenes.SettingsBuilder](b)

	pkg.LoadObjects(b)
}
