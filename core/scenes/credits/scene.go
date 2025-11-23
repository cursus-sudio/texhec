package creditsscene

import (
	gameassets "core/assets"
	gamescenes "core/scenes"
	"frontend/modules/camera"
	"frontend/modules/collider"
	"frontend/modules/drag"
	"frontend/modules/genericrenderer"
	"frontend/modules/inputs"
	"frontend/modules/render"
	scenessys "frontend/modules/scenes"
	"frontend/modules/text"
	"frontend/modules/transform"
	"frontend/services/scenes"
	"shared/services/ecs"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) LoadObjects(b ioc.Builder) {
	ioc.WrapService(b, scenes.LoadObjects, func(c ioc.Dic, b gamescenes.CreditsBuilder) gamescenes.CreditsBuilder {
		b.OnLoad(func(world scenes.SceneCtx) {
			cameraEntity := world.NewEntity()
			ecs.SaveComponent(world, cameraEntity, camera.NewOrtho(-1000, +1000, 1))

			signature := world.NewEntity()
			ecs.SaveComponent(world, signature, transform.NewPos(5, 5, 0))
			ecs.SaveComponent(world, signature, transform.NewSize(100, 50, 1))
			ecs.SaveComponent(world, signature, transform.NewPivotPoint(0, .5, .5))
			ecs.SaveComponent(world, signature, transform.NewParent(cameraEntity, transform.RelativePos))
			ecs.SaveComponent(world, signature, transform.NewParentPivotPoint(0, 0, .5))

			ecs.SaveComponent(world, signature, text.TextComponent{Text: "credits"})
			ecs.SaveComponent(world, signature, text.FontSizeComponent{FontSize: 32})
			ecs.SaveComponent(world, signature, text.BreakComponent{Break: text.BreakNone})

			background := world.NewEntity()
			ecs.SaveComponent(world, background,
				transform.NewParent(cameraEntity, transform.RelativePos|transform.RelativeSize))
			ecs.SaveComponent(world, background, transform.NewParentPivotPoint(.5, .5, .5))
			ecs.SaveComponent(world, background, render.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world, background, render.NewTexture(gameassets.MountainTileTextureID))
			ecs.SaveComponent(world, background, genericrenderer.PipelineComponent{})

			buttonArea := world.NewEntity()
			ecs.SaveComponent(world, buttonArea, transform.NewSize(500, 200, 1))
			ecs.SaveComponent(world, buttonArea, transform.NewParent(cameraEntity, transform.RelativePos))

			draggable := world.NewEntity()
			ecs.SaveComponent(world, draggable, transform.NewSize(50, 50, 3))
			ecs.SaveComponent(world, draggable, render.NewColor(mgl32.Vec4{0, 1, 0, .2}))
			ecs.SaveComponent(world, draggable, render.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world, draggable, render.NewTexture(gameassets.WaterTileTextureID))
			ecs.SaveComponent(world, draggable, genericrenderer.PipelineComponent{})

			ecs.SaveComponent(world, draggable, collider.NewCollider(gameassets.SquareColliderID))
			ecs.SaveComponent(world, draggable, inputs.NewMouseDragComponent(drag.NewDraggable(draggable)))

			ecs.SaveComponent(world, draggable, text.TextComponent{Text: strings.ToUpper("drag me")})
			ecs.SaveComponent(world, draggable, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
			ecs.SaveComponent(world, draggable, text.FontSizeComponent{FontSize: 15})
			ecs.SaveComponent(world, draggable, text.TextColorComponent{Color: mgl32.Vec4{.5, 0, 1, 1}})

			btn := world.NewEntity()
			ecs.SaveComponent(world, btn, transform.NewSize(500, 50, 2))
			ecs.SaveComponent(world, btn, transform.NewParent(buttonArea, transform.RelativePos))
			ecs.SaveComponent(world, btn, transform.NewParentPivotPoint(.5, 0, .5))

			ecs.SaveComponent(world, btn, render.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world, btn, render.NewTexture(gameassets.WaterTileTextureID))
			ecs.SaveComponent(world, btn, genericrenderer.PipelineComponent{})

			ecs.SaveComponent(world, btn, inputs.NewMouseLeftClick(scenessys.NewChangeSceneEvent(gamescenes.MenuID)))
			ecs.SaveComponent(world, btn, inputs.KeepSelectedComponent{})
			ecs.SaveComponent(world, btn, collider.NewCollider(gameassets.SquareColliderID))

			ecs.SaveComponent(world, btn, text.TextComponent{Text: strings.ToUpper("return to menu")})
			ecs.SaveComponent(world, btn, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
			ecs.SaveComponent(world, btn, text.FontSizeComponent{FontSize: 32})
		})

		return b
	})
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.CreditsBuilder { return scenes.NewSceneBuilder() })
	gamescenes.AddDefaults[gamescenes.CreditsBuilder](b)

	pkg.LoadObjects(b)
}
