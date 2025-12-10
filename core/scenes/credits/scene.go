package creditsscene

import (
	gameassets "core/assets"
	gamescenes "core/scenes"
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/drag"
	"engine/modules/genericrenderer"
	"engine/modules/hierarchy"
	"engine/modules/inputs"
	"engine/modules/render"
	scenessys "engine/modules/scenes"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/scenes"
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
		gameAssets := ioc.Get[gameassets.GameAssets](c)
		b.OnLoad(func(world ecs.World) {
			cameraEntity := world.NewEntity()
			ecs.SaveComponent(world, cameraEntity, camera.NewOrtho(-1000, +1000))

			signature := world.NewEntity()
			ecs.SaveComponent(world, signature, hierarchy.NewParent(cameraEntity))
			ecs.SaveComponent(world, signature, transform.NewPos(5, 5, 0))
			ecs.SaveComponent(world, signature, transform.NewSize(100, 50, 1))
			ecs.SaveComponent(world, signature, transform.NewPivotPoint(0, .5, .5))
			ecs.SaveComponent(world, signature, transform.NewParent(transform.RelativePos))
			ecs.SaveComponent(world, signature, transform.NewParentPivotPoint(0, 0, .5))

			ecs.SaveComponent(world, signature, text.TextComponent{Text: "credits"})
			ecs.SaveComponent(world, signature, text.FontSizeComponent{FontSize: 32})
			ecs.SaveComponent(world, signature, text.BreakComponent{Break: text.BreakNone})

			background := world.NewEntity()
			ecs.SaveComponent(world, background, hierarchy.NewParent(cameraEntity))
			ecs.SaveComponent(world, background, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))
			ecs.SaveComponent(world, background, transform.NewParentPivotPoint(.5, .5, .5))
			ecs.SaveComponent(world, background, render.NewMesh(gameAssets.SquareMesh))
			ecs.SaveComponent(world, background, render.NewTexture(gameAssets.Tiles.Mountain))
			ecs.SaveComponent(world, background, genericrenderer.PipelineComponent{})

			buttonArea := world.NewEntity()
			ecs.SaveComponent(world, buttonArea, transform.NewSize(500, 200, 1))
			ecs.SaveComponent(world, buttonArea, hierarchy.NewParent(cameraEntity))
			ecs.SaveComponent(world, buttonArea, transform.NewParent(transform.RelativePos))

			draggable := world.NewEntity()
			ecs.SaveComponent(world, draggable, transform.NewSize(50, 50, 3))
			ecs.SaveComponent(world, draggable, render.NewColor(mgl32.Vec4{0, 1, 0, .2}))
			ecs.SaveComponent(world, draggable, render.NewMesh(gameAssets.SquareMesh))
			ecs.SaveComponent(world, draggable, render.NewTexture(gameAssets.Tiles.Water))
			ecs.SaveComponent(world, draggable, genericrenderer.PipelineComponent{})

			ecs.SaveComponent(world, draggable, collider.NewCollider(gameAssets.SquareCollider))
			ecs.SaveComponent(world, draggable, inputs.NewMouseDragComponent(drag.NewDraggable(draggable)))

			ecs.SaveComponent(world, draggable, text.TextComponent{Text: strings.ToUpper("drag me")})
			ecs.SaveComponent(world, draggable, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
			ecs.SaveComponent(world, draggable, text.FontSizeComponent{FontSize: 15})
			ecs.SaveComponent(world, draggable, text.TextColorComponent{Color: mgl32.Vec4{.5, 0, 1, 1}})

			btn := world.NewEntity()
			ecs.SaveComponent(world, btn, transform.NewSize(500, 50, 2))
			ecs.SaveComponent(world, btn, hierarchy.NewParent(buttonArea))
			ecs.SaveComponent(world, btn, transform.NewParent(transform.RelativePos))
			ecs.SaveComponent(world, btn, transform.NewParentPivotPoint(.5, 0, .5))

			ecs.SaveComponent(world, btn, render.NewMesh(gameAssets.SquareMesh))
			ecs.SaveComponent(world, btn, render.NewTexture(gameAssets.Tiles.Water))
			ecs.SaveComponent(world, btn, genericrenderer.PipelineComponent{})

			ecs.SaveComponent(world, btn, inputs.NewMouseLeftClick(scenessys.NewChangeSceneEvent(gamescenes.MenuID)))
			ecs.SaveComponent(world, btn, inputs.KeepSelectedComponent{})
			ecs.SaveComponent(world, btn, collider.NewCollider(gameAssets.SquareCollider))

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
