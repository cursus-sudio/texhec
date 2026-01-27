package menuscene

import (
	gameassets "core/assets"
	gamescenes "core/scenes"
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/genericrenderer"
	"engine/modules/inputs"
	"engine/modules/layout"
	"engine/modules/render"
	"engine/modules/scene"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/modules/transition"
	"engine/services/assets"
	"engine/services/ecs"
	"strings"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) gamescenes.MenuBuilder {
		gameAssets := ioc.Get[gameassets.GameAssets](c)
		assetsService := ioc.Get[assets.Assets](c)
		return func(sceneParent ecs.EntityID) {
			world := ioc.GetServices[gamescenes.World](c)
			cameraEntity := world.NewEntity()
			world.Hierarchy.SetParent(cameraEntity, sceneParent)
			world.Camera.Ortho().Set(cameraEntity, camera.NewOrtho(-1000, 1000))

			signature := world.NewEntity()
			world.Hierarchy.SetParent(signature, cameraEntity)
			world.Transform.Pos().Set(signature, transform.NewPos(5, 5, 0))
			world.Transform.Size().Set(signature, transform.NewSize(100, 50, 1))
			world.Transform.PivotPoint().Set(signature, transform.NewPivotPoint(0, .5, .5))
			world.Transform.Parent().Set(signature, transform.NewParent(transform.RelativePos))
			world.Transform.ParentPivotPoint().Set(signature, transform.NewParentPivotPoint(0, 0, .5))

			world.Text.Content().Set(signature, text.TextComponent{Text: "menu"})
			world.Text.FontSize().Set(signature, text.FontSizeComponent{FontSize: 32})
			world.Text.Break().Set(signature, text.BreakComponent{Break: text.BreakNone})

			background := world.NewEntity()
			world.Hierarchy.SetParent(background, cameraEntity)
			world.Transform.Parent().Set(background, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))
			world.Transform.Pos().Set(background, transform.NewPos(0, 0, 1))
			world.Transform.PivotPoint().Set(background, transform.NewPivotPoint(.5, .5, 0))
			world.Transform.ParentPivotPoint().Set(background, transform.NewParentPivotPoint(.5, .5, 0))
			world.Render.Mesh().Set(background, render.NewMesh(gameAssets.SquareMesh))
			world.Render.Texture().Set(background, render.NewTexture(gameAssets.Hud.Background))
			world.GenericRenderer.Pipeline().Set(background, genericrenderer.PipelineComponent{})
			world.Transition.Easing().Set(background, transition.NewEasing(gameassets.MyEasingFunction))
			events.Emit(world.Events, transition.NewTransitionEvent(
				background,
				render.NewColor(mgl32.Vec4{1, 0, 1, 1}),
				render.NewColor(mgl32.Vec4{1, 1, 1, 1}),
				time.Second,
			))

			buttonArea := world.NewEntity()
			world.Hierarchy.SetParent(buttonArea, cameraEntity)
			world.Transform.Parent().Set(buttonArea, transform.NewParent(transform.RelativePos))

			world.Layout.Order().Set(buttonArea, layout.NewOrder(layout.OrderVectical))
			world.Layout.Align().Set(buttonArea, layout.NewAlign(.5, .5))
			world.Layout.Gap().Set(buttonArea, layout.NewGap(10))

			type Button struct {
				Text    string
				OnClick any
			}
			buttons := []Button{
				{Text: "play", OnClick: scene.NewChangeSceneEvent(gamescenes.GameID)},
				{Text: "connect", OnClick: scene.NewChangeSceneEvent(gamescenes.GameClientID)},
				{Text: "settings", OnClick: scene.NewChangeSceneEvent(gamescenes.SettingsID)},
				{Text: "credits", OnClick: scene.NewChangeSceneEvent(gamescenes.CreditsID)},
				{Text: "exit", OnClick: inputs.QuitEvent{}},
			}

			btnAsset, err := assets.GetAsset[render.TextureAsset](assetsService, gameAssets.Hud.Btn)
			if err != nil {
				world.Logger.Warn(err)
				return
			}
			btnAspectRatio := btnAsset.AspectRatio()

			for _, button := range buttons {
				btn := world.NewEntity()
				world.Hierarchy.SetParent(btn, buttonArea)
				world.Transform.Size().Set(btn, transform.NewSize(150, 50, 1))
				world.Transform.AspectRatio().Set(btn, transform.NewAspectRatio(float32(btnAspectRatio.Dx()), float32(btnAspectRatio.Dy()), 0, transform.PrimaryAxisX))
				world.Transform.Parent().Set(btn, transform.NewParent(transform.RelativePos))

				world.Render.Mesh().Set(btn, render.NewMesh(gameAssets.SquareMesh))
				world.Render.Texture().Set(btn, render.NewTexture(gameAssets.Hud.Btn))
				world.Render.TextureFrame().Set(btn, render.NewTextureFrameComponent(1))
				world.GenericRenderer.Pipeline().Set(btn, genericrenderer.PipelineComponent{})

				world.Inputs.LeftClick().Set(btn, inputs.NewLeftClick(button.OnClick))
				world.Collider.Component().Set(btn, collider.NewCollider(gameAssets.SquareCollider))
				world.Inputs.KeepSelected().Set(btn, inputs.KeepSelectedComponent{})

				world.Text.Content().Set(btn, text.TextComponent{Text: strings.ToUpper(button.Text)})
				world.Text.Align().Set(btn, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
				world.Text.FontSize().Set(btn, text.FontSizeComponent{FontSize: 24})
			}
		}
	})
}
