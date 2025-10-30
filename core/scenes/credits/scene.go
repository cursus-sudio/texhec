package creditsscene

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
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

func (Pkg) LoadObjects(b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadObjects, func(c ioc.Dic, b gamescenes.CreditsBuilder) gamescenes.CreditsBuilder {
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

			ecs.SaveComponent(world.Components(), signature, text.Text{Text: "credits"})
			ecs.SaveComponent(world.Components(), signature, text.FontSize{FontSize: 32})
			ecs.SaveComponent(world.Components(), signature, text.Break{Break: text.BreakNone})

			background := world.NewEntity()
			ecs.SaveComponent(world.Components(), background, anchor.NewParentAnchor(cameraEntity).Ptr().
				SetPivotPoint(mgl32.Vec3{.5, .5, .5}).
				SetRelativeTransform(transform.NewTransform().Ptr().SetSize(mgl32.Vec3{1, 1, 1}).Val()).Val(),
			)
			ecs.SaveComponent(world.Components(), background, mesh.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), background, texture.NewTexture(gameassets.MountainTileTextureID))
			ecs.SaveComponent(world.Components(), background, genericrenderersys.PipelineComponent{})

			buttonArea := world.NewEntity()
			ecs.SaveComponent(world.Components(), buttonArea, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{500, 200, 1}).Val())
			ecs.SaveComponent(world.Components(), buttonArea, anchor.NewParentAnchor(cameraEntity).Ptr().
				SetPivotPoint(mgl32.Vec3{.5, .5, .5}).Val())

			draggable := world.NewEntity()

			ecs.SaveComponent(world.Components(), draggable, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{50, 50, 1}).Val())

			ecs.SaveComponent(world.Components(), draggable, mesh.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), draggable, texture.NewTexture(gameassets.WaterTileTextureID))
			ecs.SaveComponent(world.Components(), draggable, genericrenderersys.PipelineComponent{})

			// ecs.SaveComponent(world.Components(), draggable, mouse.NewMouseEvents().
			// 	AddLeftClickEvents(scenessys.NewChangeSceneEvent(gamescenes.MenuID)))
			ecs.SaveComponent(world.Components(), draggable, collider.NewCollider(gameassets.SquareColliderID))
			type Hehe struct{}
			ecs.SaveComponent(world.Components(), draggable, mouse.NewDragEvents(Hehe{}))

			ecs.SaveComponent(world.Components(), draggable, text.Text{Text: strings.ToUpper("return to menu")})
			ecs.SaveComponent(world.Components(), draggable, text.TextAlign{Vertical: .5, Horizontal: .5})
			ecs.SaveComponent(world.Components(), draggable, text.FontSize{FontSize: 32})

			button := world.NewEntity()
			ecs.SaveComponent(world.Components(), button, transform.NewTransform().Ptr().
				SetSize(mgl32.Vec3{500, 50, 1}).Val())
			ecs.SaveComponent(world.Components(), button, anchor.NewParentAnchor(buttonArea).Ptr().
				SetPivotPoint(mgl32.Vec3{.5, 0, .5}).
				Val())

			ecs.SaveComponent(world.Components(), button, mesh.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), button, texture.NewTexture(gameassets.WaterTileTextureID))
			ecs.SaveComponent(world.Components(), button, genericrenderersys.PipelineComponent{})

			ecs.SaveComponent(world.Components(), button, mouse.NewMouseEvents().
				AddLeftClickEvents(scenessys.NewChangeSceneEvent(gamescenes.MenuID)))
			ecs.SaveComponent(world.Components(), button, collider.NewCollider(gameassets.SquareColliderID))

			ecs.SaveComponent(world.Components(), button, text.Text{Text: strings.ToUpper("drag me")})
			ecs.SaveComponent(world.Components(), button, text.TextAlign{Vertical: .5, Horizontal: .5})
			ecs.SaveComponent(world.Components(), button, text.FontSize{FontSize: 32})
		})

		return b
	})
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.CreditsBuilder { return scenes.NewSceneBuilder() })
	gamescenes.AddDefaults[gamescenes.CreditsBuilder](b)

	pkg.LoadObjects(b)
}
