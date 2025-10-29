package settingsscene

import (
	gameassets "core/assets"
	gamescenes "core/scenes"
	"frontend/engine/components/anchor"
	"frontend/engine/components/camera"
	"frontend/engine/components/collider"
	"frontend/engine/components/mesh"
	"frontend/engine/components/mouse"
	"frontend/engine/components/projection"
	"frontend/engine/components/text"
	"frontend/engine/components/texture"
	"frontend/engine/components/transform"
	"frontend/engine/systems/genericrenderer"
	"frontend/engine/systems/scenes"
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
	ioc.WrapService(b, scenes.LoadObjects, func(c ioc.Dic, b gamescenes.SettingsBuilder) gamescenes.SettingsBuilder {
		b.OnLoad(func(world scenes.SceneCtx) {
			cameraEntity := world.NewEntity()
			ecs.SaveComponent(world.Components(), cameraEntity, transform.NewTransform())
			ecs.SaveComponent(world.Components(), cameraEntity, projection.NewDynamicOrtho(-1000, +1000, 1))
			ecs.SaveComponent(world.Components(), cameraEntity,
				camera.NewCamera(ecs.GetComponentType(projection.Ortho{})))

			signature := world.NewEntity()
			ecs.SaveComponent(world.Components(), signature, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{100, 50, 1}).Val())
			ecs.SaveComponent(world.Components(), signature, transform.NewPivotPoint(mgl32.Vec3{1, .5, .5}))
			ecs.SaveComponent(world.Components(), signature, anchor.NewParentAnchor(cameraEntity).Ptr().
				SetPivotPoint(mgl32.Vec3{0, 0, .5}).
				SetOffset(mgl32.Vec3{5, 5}).
				Val())

			ecs.SaveComponent(world.Components(), signature, text.Text{Text: "settings"})
			ecs.SaveComponent(world.Components(), signature, text.FontSize{FontSize: 32})
			ecs.SaveComponent(world.Components(), signature, text.Break{Break: text.BreakNone})

			background := world.NewEntity()
			ecs.SaveComponent(world.Components(), background, anchor.NewParentAnchor(cameraEntity).Ptr().
				SetPivotPoint(mgl32.Vec3{.5, .5, .5}).
				SetRelativeTransform(transform.NewTransform().Ptr().SetSize(mgl32.Vec3{1, 1, 1}).Val()).Val(),
			)
			ecs.SaveComponent(world.Components(), background, mesh.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), background, texture.NewTexture(gameassets.GroundTileTextureID))
			ecs.SaveComponent(world.Components(), background, genericrenderersys.PipelineComponent{})

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
				{Text: "mute", OnClick: nil},
				{Text: "keybinds", OnClick: nil},
				{Text: "return to menu", OnClick: scenessys.NewChangeSceneEvent(gamescenes.MenuID)},
			}
			slices.Reverse(buttons)

			for i, button := range buttons {
				entity := world.NewEntity()
				normalizedIndex := float32(i) / (float32(len(buttons)) - 1)
				ecs.SaveComponent(world.Components(), entity, transform.NewTransform().Ptr().
					SetSize(mgl32.Vec3{500, 50, 1}).Val())
				ecs.SaveComponent(world.Components(), entity, anchor.NewParentAnchor(buttonArea).Ptr().
					SetPivotPoint(mgl32.Vec3{.5, normalizedIndex, .5}).
					Val())

				ecs.SaveComponent(world.Components(), entity, mesh.NewMesh(gameassets.SquareMesh))
				ecs.SaveComponent(world.Components(), entity, texture.NewTexture(gameassets.WaterTileTextureID))
				ecs.SaveComponent(world.Components(), entity, genericrenderersys.PipelineComponent{})

				ecs.SaveComponent(world.Components(), entity, mouse.NewMouseEvents().AddLeftClickEvents(button.OnClick))
				ecs.SaveComponent(world.Components(), entity, collider.NewCollider(gameassets.SquareColliderID))

				ecs.SaveComponent(world.Components(), entity, text.Text{Text: strings.ToUpper(button.Text)})
				ecs.SaveComponent(world.Components(), entity, text.TextAlign{Vertical: .5, Horizontal: .5})
				ecs.SaveComponent(world.Components(), entity, text.FontSize{FontSize: 32})
			}
		})

		return b
	})
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.SettingsBuilder { return scenes.NewSceneBuilder() })
	gamescenes.AddDefaults[gamescenes.SettingsBuilder](b)

	pkg.LoadObjects(b)
}
