package menuscene

import (
	gameassets "core/assets"
	gamescenes "core/scenes"
	"frontend/engine/anchor"
	"frontend/engine/camera"
	"frontend/engine/collider"
	"frontend/engine/genericrenderer"
	"frontend/engine/inputs"
	"frontend/engine/mesh"
	scenessys "frontend/engine/scenes"
	"frontend/engine/text"
	"frontend/engine/texture"
	"frontend/engine/transform"
	"frontend/services/scenes"
	"shared/services/ecs"
	"slices"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

func (Pkg) LoadObjects(b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadObjects, func(c ioc.Dic, b gamescenes.MenuBuilder) gamescenes.MenuBuilder {
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

			ecs.SaveComponent(world.Components(), signature, text.Text{Text: "menu"})
			ecs.SaveComponent(world.Components(), signature, text.FontSize{FontSize: 32})
			ecs.SaveComponent(world.Components(), signature, text.Break{Break: text.BreakNone})

			background := world.NewEntity()
			ecs.SaveComponent(world.Components(), background, anchor.NewParentAnchor(cameraEntity).Ptr().
				SetPivotPoint(mgl32.Vec3{.5, .5, .5}).
				SetRelativeTransform(transform.NewTransform().Ptr().SetSize(mgl32.Vec3{1, 1, 1}).Val()).Val(),
			)
			ecs.SaveComponent(world.Components(), background, mesh.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), background, texture.NewTexture(gameassets.ForestTileTextureID))
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
				{Text: "play", OnClick: scenessys.NewChangeSceneEvent(gamescenes.GameID)},
				{Text: "settings", OnClick: scenessys.NewChangeSceneEvent(gamescenes.SettingsID)},
				{Text: "credits", OnClick: scenessys.NewChangeSceneEvent(gamescenes.CreditsID)},
				{Text: "exit", OnClick: inputs.QuitEvent{}},
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

				ecs.SaveComponent(world.Components(), btn, mesh.NewMesh(gameassets.SquareMesh))
				ecs.SaveComponent(world.Components(), btn, texture.NewTexture(gameassets.WaterTileTextureID))
				ecs.SaveComponent(world.Components(), btn, genericrenderer.PipelineComponent{})

				ecs.SaveComponent(world.Components(), btn, inputs.NewMouseEvents().AddLeftClickEvents(button.OnClick))
				ecs.SaveComponent(world.Components(), btn, collider.NewCollider(gameassets.SquareColliderID))
				ecs.SaveComponent(world.Components(), btn, inputs.KeepSelected{})

				ecs.SaveComponent(world.Components(), btn, text.Text{Text: strings.ToUpper(button.Text)})
				ecs.SaveComponent(world.Components(), btn, text.TextAlign{Vertical: .5, Horizontal: .5})
				ecs.SaveComponent(world.Components(), btn, text.FontSize{FontSize: 32})
			}
		})

		return b
	})
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.MenuBuilder { return scenes.NewSceneBuilder() })
	gamescenes.AddDefaults[gamescenes.MenuBuilder](b)

	pkg.LoadObjects(b)
}
