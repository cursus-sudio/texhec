package settingsscene

import (
	gameassets "core/assets"
	"core/modules/settings"
	gamescenes "core/scenes"
	"engine/modules/camera"
	"engine/modules/genericrenderer"
	"engine/modules/hierarchy"
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
		b.OnLoad(func(world ecs.World) {
			cameraEntity := world.NewEntity()
			ecs.SaveComponent(world, cameraEntity, camera.NewOrtho(-1000, +1000))

			signature := world.NewEntity()
			ecs.SaveComponent(world, signature, transform.NewPos(5, 5, 0))
			ecs.SaveComponent(world, signature, transform.NewSize(100, 50, 1))
			ecs.SaveComponent(world, signature, transform.NewPivotPoint(0, .5, .5))
			ecs.SaveComponent(world, signature, hierarchy.NewParent(cameraEntity))
			ecs.SaveComponent(world, signature, transform.NewParent(transform.RelativePos))
			ecs.SaveComponent(world, signature, transform.NewParentPivotPoint(0, 0, .5))

			ecs.SaveComponent(world, signature, text.TextComponent{Text: "settings"})
			ecs.SaveComponent(world, signature, text.FontSizeComponent{FontSize: 32})
			ecs.SaveComponent(world, signature, text.BreakComponent{Break: text.BreakNone})

			background := world.NewEntity()
			ecs.SaveComponent(world, background, hierarchy.NewParent(cameraEntity))
			ecs.SaveComponent(world, background, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))
			ecs.SaveComponent(world, background, render.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world, background, render.NewTexture(gameassets.Tiles.Ground))
			ecs.SaveComponent(world, background, genericrenderer.PipelineComponent{})

			buttonArea := world.NewEntity()
			ecs.SaveComponent(world, buttonArea, transform.NewSize(500, 200, 2))
			// ecs.SaveComponent(world, buttonArea, transform.NewPos(0, 0, -50))
			ecs.SaveComponent(world, buttonArea, hierarchy.NewParent(cameraEntity))
			ecs.SaveComponent(world, buttonArea, transform.NewParent(transform.RelativePos))

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
