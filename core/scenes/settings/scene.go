package settingsscene

import (
	gameassets "core/assets"
	gamescenes "core/scenes"
	"frontend/modules/anchor"
	"frontend/modules/audio"
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

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) LoadObjects(b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadObjects, func(c ioc.Dic, b gamescenes.SettingsBuilder) gamescenes.SettingsBuilder {
		b.OnLoad(func(world scenes.SceneCtx) {
			cameraEntity := world.NewEntity()
			ecs.SaveComponent(world.Components(), cameraEntity, transform.NewTransform())
			ecs.SaveComponent(world.Components(), cameraEntity, camera.NewDynamicOrtho(-1000, +1000, 1))

			signature := world.NewEntity()
			ecs.SaveComponent(world.Components(), signature, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{100, 50, 1}).Val())
			ecs.SaveComponent(world.Components(), signature, transform.NewPivotPoint(mgl32.Vec3{1, .5, .5}))
			ecs.SaveComponent(world.Components(), signature, anchor.NewParentAnchor(cameraEntity).Ptr().
				SetPivotPoint(mgl32.Vec3{0, 0, .5}).
				SetOffset(mgl32.Vec3{5, 5}).
				Val())

			ecs.SaveComponent(world.Components(), signature, text.TextComponent{Text: "settings"})
			ecs.SaveComponent(world.Components(), signature, text.FontSizeComponent{FontSize: 32})
			ecs.SaveComponent(world.Components(), signature, text.BreakComponent{Break: text.BreakNone})

			background := world.NewEntity()
			ecs.SaveComponent(world.Components(), background, anchor.NewParentAnchor(cameraEntity).Ptr().
				SetPivotPoint(mgl32.Vec3{.5, .5, .5}).
				SetRelativeTransform(transform.NewTransform().Ptr().SetSize(mgl32.Vec3{1, 1, 1}).Val()).Val(),
			)
			ecs.SaveComponent(world.Components(), background, render.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), background, render.NewTexture(gameassets.GroundTileTextureID))
			ecs.SaveComponent(world.Components(), background, genericrenderer.PipelineComponent{})

			buttonArea := world.NewEntity()
			ecs.SaveComponent(world.Components(), buttonArea, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{500, 200, 1}).Val())
			ecs.SaveComponent(world.Components(), buttonArea, anchor.NewParentAnchor(cameraEntity).Ptr().
				SetPivotPoint(mgl32.Vec3{.5, .5, .5}).Val())

			type Button struct {
				Text    string
				OnClick any
			}
			buttons := []Button{
				{Text: "mute", OnClick: audio.NewPlayEvent(gamescenes.EffectChannel, gameassets.AudioID)},
				{Text: "keybinds", OnClick: nil},
				{Text: "return to menu", OnClick: scenessys.NewChangeSceneEvent(gamescenes.MenuID)},
			}
			slices.Reverse(buttons)

			for i, button := range buttons {
				btn := world.NewEntity()
				normalizedIndex := float32(i) / (float32(len(buttons)) - 1)
				ecs.SaveComponent(world.Components(), btn, transform.NewTransform().Ptr().
					SetSize(mgl32.Vec3{500, 50, 1}).Val())
				ecs.SaveComponent(world.Components(), btn, anchor.NewParentAnchor(buttonArea).Ptr().
					SetPivotPoint(mgl32.Vec3{.5, normalizedIndex, .5}).
					Val())

				ecs.SaveComponent(world.Components(), btn, render.NewMesh(gameassets.SquareMesh))
				ecs.SaveComponent(world.Components(), btn, render.NewTexture(gameassets.WaterTileTextureID))
				ecs.SaveComponent(world.Components(), btn, genericrenderer.PipelineComponent{})

				ecs.SaveComponent(world.Components(), btn, inputs.NewMouseEvents().Ptr().
					AddLeftClickEvents(button.OnClick).Val())
				ecs.SaveComponent(world.Components(), btn, collider.NewCollider(gameassets.SquareColliderID))
				ecs.SaveComponent(world.Components(), btn, inputs.KeepSelectedComponent{})

				ecs.SaveComponent(world.Components(), btn, text.TextComponent{Text: strings.ToUpper(button.Text)})
				ecs.SaveComponent(world.Components(), btn, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
				ecs.SaveComponent(world.Components(), btn, text.FontSizeComponent{FontSize: 32})
			}
		})

		return b
	})
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.SettingsBuilder { return scenes.NewSceneBuilder() })
	gamescenes.AddDefaults[gamescenes.SettingsBuilder](b)

	pkg.LoadObjects(b)
}
