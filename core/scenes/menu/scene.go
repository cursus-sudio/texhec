package menuscene

import (
	gameassets "core/assets"
	gamescenes "core/scenes"
	"frontend/modules/animation"
	"frontend/modules/camera"
	"frontend/modules/collider"
	"frontend/modules/genericrenderer"
	"frontend/modules/inputs"
	"frontend/modules/render"
	scenessys "frontend/modules/scenes"
	"frontend/modules/text"
	"frontend/modules/transform"
	"frontend/services/scenes"
	"shared/services/ecs"
	"slices"
	"strings"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) LoadObjects(b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadObjects, func(c ioc.Dic, b gamescenes.MenuBuilder) gamescenes.MenuBuilder {
		b.OnLoad(func(world scenes.SceneCtx) {
			cameraEntity := world.NewEntity()
			ecs.SaveComponent(world, cameraEntity, camera.NewOrtho(-1000, +1000, 1))

			signature := world.NewEntity()
			ecs.SaveComponent(world, signature, transform.NewPos(5, 5, 0))
			ecs.SaveComponent(world, signature, transform.NewSize(100, 50, 1))
			ecs.SaveComponent(world, signature, transform.NewPivotPoint(0, .5, .5))
			ecs.SaveComponent(world, signature, transform.NewParent(cameraEntity, transform.RelativePos))
			ecs.SaveComponent(world, signature, transform.NewParentPivotPoint(0, 0, .5))

			ecs.SaveComponent(world, signature, text.TextComponent{Text: "menu"})
			ecs.SaveComponent(world, signature, text.FontSizeComponent{FontSize: 32})
			ecs.SaveComponent(world, signature, text.BreakComponent{Break: text.BreakNone})

			background := world.NewEntity()
			ecs.SaveComponent(world, background, transform.NewParent(cameraEntity, transform.RelativePos|transform.RelativeSize))
			ecs.SaveComponent(world, background, render.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world, background, render.NewTexture(gameassets.ForestTileTextureID))
			ecs.SaveComponent(world, background, genericrenderer.PipelineComponent{})
			ecs.SaveComponent(world, background, animation.NewAnimationComponent(
				gameassets.ChangeColorsAnimation,
				time.Second,
			))
			// ecs.SaveComponent(world, background, animation.NewLoopComponent())

			buttonArea := world.NewEntity()
			ecs.SaveComponent(world, buttonArea, transform.NewSize(500, 200, 1))
			ecs.SaveComponent(world, buttonArea, transform.NewParent(cameraEntity, transform.RelativePos))

			type Button struct {
				Text    string
				OnClick any
			}
			buttons := []Button{
				{Text: "play", OnClick: scenessys.NewChangeSceneEvent(gamescenes.GameID)},
				{Text: "settings", OnClick: scenessys.NewChangeSceneEvent(gamescenes.SettingsID)},
				{Text: "credits", OnClick: scenessys.NewChangeSceneEvent(gamescenes.CreditsID)},
				{Text: "exit", OnClick: inputs.QuitEvent{}},
			}
			slices.Reverse(buttons)

			for i, button := range buttons {
				btn := world.NewEntity()
				normalizedIndex := float32(i) / (float32(len(buttons)) - 1)
				ecs.SaveComponent(world, btn, transform.NewSize(500, 50, 2))
				ecs.SaveComponent(world, btn, transform.NewParent(buttonArea, transform.RelativePos))
				ecs.SaveComponent(world, btn, transform.NewParentPivotPoint(.5, normalizedIndex, .5))

				ecs.SaveComponent(world, btn, render.NewMesh(gameassets.SquareMesh))
				ecs.SaveComponent(world, btn, render.NewTexture(gameassets.WaterTileTextureID))
				ecs.SaveComponent(world, btn, render.NewTextureFrameComponent(1))
				ecs.SaveComponent(world, btn, genericrenderer.PipelineComponent{})
				ecs.SaveComponent(world, btn, animation.NewAnimationComponent(
					gameassets.ButtonAnimation,
					time.Second*2,
				))
				ecs.SaveComponent(world, btn, animation.NewLoopComponent())

				ecs.SaveComponent(world, btn, inputs.NewMouseLeftClick(button.OnClick))
				ecs.SaveComponent(world, btn, collider.NewCollider(gameassets.SquareColliderID))
				ecs.SaveComponent(world, btn, inputs.KeepSelectedComponent{})

				ecs.SaveComponent(world, btn, text.TextComponent{Text: strings.ToUpper(button.Text)})
				ecs.SaveComponent(world, btn, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
				ecs.SaveComponent(world, btn, text.FontSizeComponent{FontSize: 32})
			}
		})

		return b
	})
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.MenuBuilder { return scenes.NewSceneBuilder() })
	gamescenes.AddDefaults[gamescenes.MenuBuilder](b)

	pkg.LoadObjects(b)
}
