package creditsscene

import (
	gameassets "core/assets"
	gamescenes "core/scenes"
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/drag"
	"engine/modules/inputs"
	"engine/modules/render"
	"engine/modules/scene"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/assets"
	"engine/services/ecs"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.CreditsBuilder {
		gameAssets := ioc.Get[gameassets.GameAssets](c)
		assetsService := ioc.Get[assets.Assets](c)
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

			world.Text.Content().Set(signature, text.TextComponent{Text: "credits"})
			world.Text.FontSize().Set(signature, text.FontSizeComponent{FontSize: 32})
			world.Text.Break().Set(signature, text.BreakComponent{Break: text.BreakNone})

			background := world.NewEntity()
			world.Hierarchy.SetParent(background, cameraEntity)
			world.Transform.Parent().Set(background, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))
			world.Transform.ParentPivotPoint().Set(background, transform.NewParentPivotPoint(.5, .5, .5))
			world.Render.Mesh().Set(background, render.NewMesh(gameAssets.SquareMesh))
			world.Render.Texture().Set(background, render.NewTexture(gameAssets.Hud.Background))
			world.Render.Render(background)

			buttonArea := world.NewEntity()
			world.Hierarchy.SetParent(buttonArea, cameraEntity)
			world.Transform.Pos().Set(buttonArea, transform.NewPos(0, 0, 1))
			world.Transform.Size().Set(buttonArea, transform.NewSize(500, 200, 1))
			world.Transform.Parent().Set(buttonArea, transform.NewParent(transform.RelativePos))

			draggable := world.NewEntity()
			world.Hierarchy.SetParent(draggable, sceneParent)
			world.Transform.Pos().Set(draggable, transform.NewPos(0, 0, 2))
			world.Transform.Size().Set(draggable, transform.NewSize(50, 50, 1))
			world.Render.Color().Set(draggable, render.NewColor(mgl32.Vec4{0, 1, 0, .2}))
			world.Render.Mesh().Set(draggable, render.NewMesh(gameAssets.SquareMesh))
			world.Render.Texture().Set(draggable, render.NewTexture(gameAssets.Hud.Btn))
			world.Render.Render(draggable)

			world.Collider.Component().Set(draggable, collider.NewCollider(gameAssets.SquareCollider))
			world.Inputs.Drag().Set(draggable, inputs.NewDragComponent(drag.NewDraggable(draggable)))

			world.Text.Content().Set(draggable, text.TextComponent{Text: strings.ToUpper("drag me")})
			world.Text.Align().Set(draggable, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
			world.Text.FontSize().Set(draggable, text.FontSizeComponent{FontSize: 15})
			world.Text.Color().Set(draggable, text.TextColorComponent{Color: mgl32.Vec4{.5, 0, 1, 1}})

			btnAsset, err := assets.GetAsset[render.TextureAsset](assetsService, gameAssets.Hud.Btn)
			if err != nil {
				world.Logger.Warn(err)
				return
			}
			btnAspectRatio := btnAsset.AspectRatio()

			btn := world.NewEntity()
			world.Hierarchy.SetParent(btn, buttonArea)
			world.Transform.Size().Set(btn, transform.NewSize(500, 100, 1))
			world.Transform.Parent().Set(btn, transform.NewParent(transform.RelativePos))
			world.Transform.ParentPivotPoint().Set(btn, transform.NewParentPivotPoint(.5, 0, .5))
			world.Transform.AspectRatio().Set(btn, transform.NewAspectRatio(float32(btnAspectRatio.Dx()), float32(btnAspectRatio.Dy()), 0, transform.PrimaryAxisY))

			world.Render.Mesh().Set(btn, render.NewMesh(gameAssets.SquareMesh))
			world.Render.Texture().Set(btn, render.NewTexture(gameAssets.Hud.Btn))
			world.Render.Render(btn)

			world.Inputs.LeftClick().Set(btn, inputs.NewLeftClick(scene.NewChangeSceneEvent(gamescenes.MenuID)))
			world.Inputs.KeepSelected().Set(btn, inputs.KeepSelectedComponent{})
			world.Collider.Component().Set(btn, collider.NewCollider(gameAssets.SquareCollider))

			world.Text.Content().Set(btn, text.TextComponent{Text: strings.ToUpper("return to menu")})
			world.Text.Align().Set(btn, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
			world.Text.FontSize().Set(btn, text.FontSizeComponent{FontSize: 32})
		}
	})
}
