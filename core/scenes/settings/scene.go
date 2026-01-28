package settingsscene

import (
	gameassets "core/assets"
	"core/modules/settings"
	gamescenes "core/scenes"
	"engine/modules/camera"
	"engine/modules/layout"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/ecs"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.SettingsBuilder {
		gameassets := ioc.Get[gameassets.GameAssets](c)
		return func(sceneParent ecs.EntityID) {
			world := ioc.GetServices[gamescenes.World](c)
			cameraEntity := world.NewEntity()
			world.Hierarchy.SetParent(cameraEntity, sceneParent)
			world.Camera.Ortho().Set(cameraEntity, camera.NewOrtho(-1000, +1000))

			signature := world.NewEntity()
			world.Hierarchy.SetParent(signature, cameraEntity)
			world.Transform.Pos().Set(signature, transform.NewPos(5, 5, 0))
			world.Transform.Size().Set(signature, transform.NewSize(100, 50, 1))
			world.Transform.PivotPoint().Set(signature, transform.NewPivotPoint(0, .5, .5))
			world.Transform.Parent().Set(signature, transform.NewParent(transform.RelativePos))
			world.Transform.ParentPivotPoint().Set(signature, transform.NewParentPivotPoint(0, 0, .5))

			world.Text.Content().Set(signature, text.TextComponent{Text: "settings"})
			world.Text.FontSize().Set(signature, text.FontSizeComponent{FontSize: 32})
			world.Text.Break().Set(signature, text.BreakComponent{Break: text.BreakNone})

			background := world.NewEntity()
			world.Hierarchy.SetParent(background, cameraEntity)
			world.Transform.Parent().Set(background, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))
			world.Render.Mesh().Set(background, render.NewMesh(gameassets.SquareMesh))
			world.Render.Texture().Set(background, render.NewTexture(gameassets.Hud.Background))
			world.Render.Render(background)

			buttonArea := world.NewEntity()
			world.Transform.Size().Set(buttonArea, transform.NewSize(500, 200, 1))
			world.Hierarchy.SetParent(buttonArea, cameraEntity)
			world.Transform.Parent().Set(buttonArea, transform.NewParent(transform.RelativePos))

			world.Layout.Order().Set(buttonArea, layout.NewOrder(layout.OrderVectical))
			world.Layout.Align().Set(buttonArea, layout.NewAlign(.5, .5))
			world.Layout.Gap().Set(buttonArea, layout.NewGap(10))

			events.Emit(world.Events, settings.EnterSettingsForParentEvent{Parent: buttonArea})
		}
	})
}
