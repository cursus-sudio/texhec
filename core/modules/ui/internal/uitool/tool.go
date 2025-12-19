package uitool

import (
	gameassets "core/assets"
	"core/modules/ui"
	"engine/modules/collider"
	"engine/modules/genericrenderer"
	"engine/modules/groups"
	"engine/modules/inputs"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/modules/transition"
	"engine/services/ecs"
	"engine/services/logger"
	"errors"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

type changing struct {
	isInitialized bool
	active        bool
	menu          ecs.EntityID
	childWrapper  ecs.EntityID
}

type tool struct {
	*changing

	animationDuration time.Duration

	ui.World
	gameAssets gameassets.GameAssets
	logger     logger.Logger

	uiCameraArray ecs.ComponentsArray[ui.UiCameraComponent]
}

func NewTool(
	animationDuration time.Duration,
	world ui.World,
	gameAssets gameassets.GameAssets,
	logger logger.Logger,
) tool {
	t := tool{
		changing: &changing{},

		animationDuration: animationDuration,

		World:      world,
		gameAssets: gameAssets,
		logger:     logger,

		uiCameraArray: ecs.GetComponentsArray[ui.UiCameraComponent](world),
	}

	t.logger.Warn(t.Init())

	return t
}

func (t tool) Init() error {
	cameras := t.uiCameraArray.GetEntities()
	if len(cameras) == 0 {
		return nil
	}
	if len(cameras) != 1 {
		return errors.New("expected one camera")
	}
	camera := cameras[0]

	// transactions

	// objects
	// menu
	menu := t.NewEntity()
	t.Hierarchy().SetParent(menu, camera)
	t.Transform().Parent().Set(menu, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))
	t.Transform().ParentPivotPoint().Set(menu, transform.NewParentPivotPoint(1, 1, .5))
	t.Transform().Pos().Set(menu, transform.NewPos(0, 0, 1))
	t.Transform().Size().Set(menu, transform.NewSize(.2, 1, 1))
	t.Transform().PivotPoint().Set(menu, transform.NewPivotPoint(0, 1, .5))

	t.Render().Color().Set(menu, render.NewColor(mgl32.Vec4{1, 1, 1, .5}))
	t.Render().Mesh().Set(menu, render.NewMesh(t.gameAssets.SquareMesh))
	t.Render().Texture().Set(menu, render.NewTexture(t.gameAssets.Tiles.Water))
	t.GenericRenderer().Pipeline().Set(menu, genericrenderer.PipelineComponent{})

	t.Groups().Inherit().Set(menu, groups.InheritGroupsComponent{})
	t.Collider().Component().Set(menu, collider.NewCollider(t.gameAssets.SquareCollider))
	t.Inputs().KeepSelected().Set(menu, inputs.KeepSelectedComponent{})

	// quit btn
	quit := t.NewEntity()

	t.Hierarchy().SetParent(quit, menu)
	t.Groups().Inherit().Set(quit, groups.InheritGroupsComponent{})

	t.Transform().Parent().Set(quit, transform.NewParent(transform.RelativePos))
	t.Transform().ParentPivotPoint().Set(quit, transform.NewParentPivotPoint(1, 1, 1))
	t.Transform().Size().Set(quit, transform.NewSize(25, 25, 1))
	t.Transform().PivotPoint().Set(quit, transform.NewPivotPoint(1, 1, 0))

	t.Text().Content().Set(quit, text.TextComponent{Text: "X"})
	t.Text().FontSize().Set(quit, text.FontSizeComponent{FontSize: 25})
	t.Text().Align().Set(quit, text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	t.Render().Color().Set(quit, render.NewColor(mgl32.Vec4{1, 0, 0, 1}))
	t.Render().Mesh().Set(quit, render.NewMesh(t.gameAssets.SquareMesh))
	t.Render().Texture().Set(quit, render.NewTexture(t.gameAssets.Tiles.Water))
	t.GenericRenderer().Pipeline().Set(quit, genericrenderer.PipelineComponent{})

	t.Inputs().MouseLeft().Set(quit, inputs.NewMouseLeftClick(ui.HideUiEvent{}))
	t.Inputs().KeepSelected().Set(quit, inputs.KeepSelectedComponent{})
	t.Collider().Component().Set(quit, collider.NewCollider(t.gameAssets.SquareCollider))

	// child wrapper
	childWrapper := t.NewEntity()
	t.Hierarchy().SetParent(childWrapper, menu)
	t.Groups().Inherit().Set(childWrapper, groups.InheritGroupsComponent{})
	t.Transform().Parent().Set(childWrapper, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))

	t.menu = menu
	t.childWrapper = childWrapper
	t.isInitialized = true

	events.Listen(t.EventsBuilder(), func(e ui.HideUiEvent) {
		t.Hide()
	})
	return nil
}

func (t tool) ResetChildWrapper() error {
	if !t.isInitialized {
		err := t.Init()
		if err != nil {
			return err
		}
	}

	for _, child := range t.Hierarchy().Children(t.childWrapper).GetIndices() {
		t.RemoveEntity(child)
	}
	return nil
}

func (t tool) Ui() ui.Interface { return t }

func (t tool) UiCamera() ecs.ComponentsArray[ui.UiCameraComponent] { return t.uiCameraArray }

func (t tool) Show() ecs.EntityID {
	t.logger.Warn(t.ResetChildWrapper())
	if !t.active {
		t.active = true
		events.Emit(t.Events(), transition.NewTransitionEvent(
			t.menu,
			transform.NewPivotPoint(0, 1, .5),
			transform.NewPivotPoint(1, 1, .5),
			t.animationDuration,
		))
	}
	return t.childWrapper
}

func (t tool) Hide() {
	if t.active {
		t.active = false
		events.Emit(t.Events(), transition.NewTransitionEvent(
			t.menu,
			transform.NewPivotPoint(1, 1, .5),
			transform.NewPivotPoint(0, 1, .5),
			t.animationDuration,
		))
	}
}

// func (t tool) Buttons(buttons ...ui.Button) []ecs.EntityID {
// 	entities := []ecs.EntityID{}
//
// 	btnAsset, err := assets.GetAsset[render.TextureAsset](assetsService, gameAssets.Hud.Btn)
// 	if err != nil {
// 		t.logger.Warn(err)
// 		return
// 	}
// 	btnAspectRatio := btnAsset.AspectRatio()
//
// 	for i, button := range buttons {
// 		btn := t.world.NewEntity()
// 		normalizedIndex := float32(i) / (float32(len(buttons)) - 1)
// 		ecs.SaveComponent(world, btn, transform.NewSize(150, 50, 1))
// 		ecs.SaveComponent(world, btn, transform.NewAspectRatio(float32(btnAspectRatio.Dx()), float32(btnAspectRatio.Dy()), 0, transform.PrimaryAxisX))
// 		ecs.SaveComponent(world, btn, hierarchy.NewParent(buttonArea))
// 		ecs.SaveComponent(world, btn, transform.NewParent(transform.RelativePos))
// 		ecs.SaveComponent(world, btn, transform.NewParentPivotPoint(.5, normalizedIndex, .5))
//
// 		ecs.SaveComponent(world, btn, render.NewMesh(gameAssets.SquareMesh))
// 		ecs.SaveComponent(world, btn, render.NewTexture(gameAssets.Hud.Btn))
// 		ecs.SaveComponent(world, btn, render.NewTextureFrameComponent(1))
// 		ecs.SaveComponent(world, btn, genericrenderer.PipelineComponent{})
//
// 		ecs.SaveComponent(world, btn, inputs.NewMouseLeftClick(button.OnClick))
// 		ecs.SaveComponent(world, btn, collider.NewCollider(gameAssets.SquareCollider))
// 		ecs.SaveComponent(world, btn, inputs.KeepSelectedComponent{})
//
// 		ecs.SaveComponent(world, btn, text.TextComponent{Text: strings.ToUpper(button.Text)})
// 		ecs.SaveComponent(world, btn, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
// 		// ecs.SaveComponent(world, btn, text.FontSizeComponent{FontSize: 32})
// 		ecs.SaveComponent(world, btn, text.FontSizeComponent{FontSize: 24})
// 	}
// 	return entities
// }
