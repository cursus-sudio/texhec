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
			ecs.SaveComponent(world.Components(), cameraEntity, camera.NewDynamicOrtho(-1000, +1000, 1))

			signature := world.NewEntity()
			ecs.SaveComponent(world.Components(), signature, transform.NewPos(5, 5, 0))
			ecs.SaveComponent(world.Components(), signature, transform.NewSize(100, 50, 1))
			ecs.SaveComponent(world.Components(), signature, transform.NewPivotPoint(1, .5, .5))
			ecs.SaveComponent(world.Components(), signature, transform.NewParent(cameraEntity, transform.RelativePos))
			ecs.SaveComponent(world.Components(), signature, transform.NewParentPivotPoint(0, 0, .5))

			ecs.SaveComponent(world.Components(), signature, text.TextComponent{Text: "credits"})
			ecs.SaveComponent(world.Components(), signature, text.FontSizeComponent{FontSize: 32})
			ecs.SaveComponent(world.Components(), signature, text.BreakComponent{Break: text.BreakNone})

			background := world.NewEntity()
			ecs.SaveComponent(world.Components(), background,
				transform.NewParent(cameraEntity, transform.RelativePos|transform.RelativeSize))
			ecs.SaveComponent(world.Components(), background, transform.NewParentPivotPoint(.5, .5, .5))
			ecs.SaveComponent(world.Components(), background, render.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), background, render.NewTexture(gameassets.MountainTileTextureID))
			ecs.SaveComponent(world.Components(), background, genericrenderer.PipelineComponent{})

			buttonArea := world.NewEntity()
			ecs.SaveComponent(world.Components(), buttonArea, transform.NewSize(500, 200, 1))
			ecs.SaveComponent(world.Components(), buttonArea, transform.NewParent(cameraEntity, transform.RelativePos))

			draggable := world.NewEntity()
			ecs.SaveComponent(world.Components(), draggable, transform.NewSize(50, 50, 3))
			ecs.SaveComponent(world.Components(), draggable, render.NewColor(mgl32.Vec4{0, 1, 0, .2}))
			ecs.SaveComponent(world.Components(), draggable, render.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), draggable, render.NewTexture(gameassets.WaterTileTextureID))
			ecs.SaveComponent(world.Components(), draggable, genericrenderer.PipelineComponent{})

			ecs.SaveComponent(world.Components(), draggable, collider.NewCollider(gameassets.SquareColliderID))
			ecs.SaveComponent(world.Components(), draggable, inputs.NewMouseEvents().Ptr().
				AddDragEvents(drag.NewDraggable(draggable)).Val())

			ecs.SaveComponent(world.Components(), draggable, text.TextComponent{Text: strings.ToUpper("drag me")})
			ecs.SaveComponent(world.Components(), draggable, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
			ecs.SaveComponent(world.Components(), draggable, text.FontSizeComponent{FontSize: 15})
			ecs.SaveComponent(world.Components(), draggable, text.TextColorComponent{Color: mgl32.Vec4{.5, 0, 1, 1}})

			btn := world.NewEntity()
			ecs.SaveComponent(world.Components(), btn, transform.NewSize(500, 50, 2))
			ecs.SaveComponent(world.Components(), btn, transform.NewParent(buttonArea, transform.RelativePos))
			ecs.SaveComponent(world.Components(), btn, transform.NewParentPivotPoint(.5, 0, .5))

			ecs.SaveComponent(world.Components(), btn, render.NewMesh(gameassets.SquareMesh))
			ecs.SaveComponent(world.Components(), btn, render.NewTexture(gameassets.WaterTileTextureID))
			ecs.SaveComponent(world.Components(), btn, genericrenderer.PipelineComponent{})

			ecs.SaveComponent(world.Components(), btn, inputs.NewMouseEvents().Ptr().
				AddLeftClickEvents(scenessys.NewChangeSceneEvent(gamescenes.MenuID)).Val())
			ecs.SaveComponent(world.Components(), btn, inputs.KeepSelectedComponent{})
			ecs.SaveComponent(world.Components(), btn, collider.NewCollider(gameassets.SquareColliderID))

			ecs.SaveComponent(world.Components(), btn, text.TextComponent{Text: strings.ToUpper("return to menu")})
			ecs.SaveComponent(world.Components(), btn, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
			ecs.SaveComponent(world.Components(), btn, text.FontSizeComponent{FontSize: 32})
		})

		return b
	})
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.CreditsBuilder { return scenes.NewSceneBuilder() })
	gamescenes.AddDefaults[gamescenes.CreditsBuilder](b)

	pkg.LoadObjects(b)
}
