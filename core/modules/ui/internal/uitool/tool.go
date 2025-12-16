package uitool

import (
	gameassets "core/assets"
	"core/modules/tile"
	"core/modules/ui"
	"engine/modules/animation"
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/genericrenderer"
	"engine/modules/groups"
	"engine/modules/hierarchy"
	"engine/modules/inputs"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
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
	showAnimation     animation.AnimationID
	hideAnimation     animation.AnimationID

	world         ecs.World
	gameAssets    gameassets.GameAssets
	logger        logger.Logger
	cameraTool    camera.Interface
	transformTool transform.Interface
	tileTool      tile.Interface
	textTool      text.Interface
	renderTool    render.Interface
	hierarchyTool hierarchy.Interface

	uiCameraArray ecs.ComponentsArray[ui.UiCameraComponent]

	pipelineArray     ecs.ComponentsArray[genericrenderer.PipelineComponent]
	groupInheritArray ecs.ComponentsArray[groups.InheritGroupsComponent]
	leftClickArray    ecs.ComponentsArray[inputs.MouseLeftClickComponent]
	keepSelectedArray ecs.ComponentsArray[inputs.KeepSelectedComponent]
	colliderArray     ecs.ComponentsArray[collider.ColliderComponent]
	animationArray    ecs.ComponentsArray[animation.AnimationComponent]
}

func NewTool(
	animationDuration time.Duration,
	showAnimation animation.AnimationID,
	hideAnimation animation.AnimationID,
	world ecs.World,
	gameAssets gameassets.GameAssets,
	logger logger.Logger,
	cameraToolFactory ecs.ToolFactory[camera.Camera],
	transformToolFactory ecs.ToolFactory[transform.Transform],
	tileToolFactory ecs.ToolFactory[tile.Tile],
	textToolFactory ecs.ToolFactory[text.Text],
	renderToolFactory ecs.ToolFactory[render.Render],
	hierarchyToolFactory ecs.ToolFactory[hierarchy.Hierarchy],
) tool {
	t := tool{
		changing: &changing{},

		animationDuration: animationDuration,
		showAnimation:     showAnimation,
		hideAnimation:     hideAnimation,

		world:         world,
		gameAssets:    gameAssets,
		logger:        logger,
		cameraTool:    cameraToolFactory.Build(world).Camera(),
		transformTool: transformToolFactory.Build(world).Transform(),
		tileTool:      tileToolFactory.Build(world).Tile(),
		textTool:      textToolFactory.Build(world).Text(),
		renderTool:    renderToolFactory.Build(world).Render(),
		hierarchyTool: hierarchyToolFactory.Build(world).Hierarchy(),

		uiCameraArray: ecs.GetComponentsArray[ui.UiCameraComponent](world),

		pipelineArray:     ecs.GetComponentsArray[genericrenderer.PipelineComponent](world),
		groupInheritArray: ecs.GetComponentsArray[groups.InheritGroupsComponent](world),
		leftClickArray:    ecs.GetComponentsArray[inputs.MouseLeftClickComponent](world),
		keepSelectedArray: ecs.GetComponentsArray[inputs.KeepSelectedComponent](world),
		colliderArray:     ecs.GetComponentsArray[collider.ColliderComponent](world),
		animationArray:    ecs.GetComponentsArray[animation.AnimationComponent](world),
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
	menu := t.world.NewEntity()
	t.hierarchyTool.SetParent(menu, camera)
	t.transformTool.Parent().SaveComponent(menu, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))
	t.transformTool.ParentPivotPoint().SaveComponent(menu, transform.NewParentPivotPoint(1, 1, .5))
	t.transformTool.Pos().SaveComponent(menu, transform.NewPos(0, 0, 1))
	t.transformTool.Size().SaveComponent(menu, transform.NewSize(.2, 1, 1))
	t.transformTool.PivotPoint().SaveComponent(menu, transform.NewPivotPoint(0, 1, .5))

	t.renderTool.Color().SaveComponent(menu, render.NewColor(mgl32.Vec4{1, 1, 1, .5}))
	t.renderTool.Mesh().SaveComponent(menu, render.NewMesh(t.gameAssets.SquareMesh))
	t.renderTool.Texture().SaveComponent(menu, render.NewTexture(t.gameAssets.Tiles.Water))
	t.pipelineArray.SaveComponent(menu, genericrenderer.PipelineComponent{})

	t.groupInheritArray.SaveComponent(menu, groups.InheritGroupsComponent{})
	t.colliderArray.SaveComponent(menu, collider.NewCollider(t.gameAssets.SquareCollider))
	t.keepSelectedArray.SaveComponent(menu, inputs.KeepSelectedComponent{})

	// quit btn
	quit := t.world.NewEntity()

	t.hierarchyTool.SetParent(quit, menu)
	t.groupInheritArray.SaveComponent(quit, groups.InheritGroupsComponent{})

	t.transformTool.Parent().SaveComponent(quit, transform.NewParent(transform.RelativePos))
	t.transformTool.ParentPivotPoint().SaveComponent(quit, transform.NewParentPivotPoint(1, 1, 1))
	t.transformTool.Size().SaveComponent(quit, transform.NewSize(25, 25, 1))
	t.transformTool.PivotPoint().SaveComponent(quit, transform.NewPivotPoint(1, 1, 0))

	t.textTool.TextContent().SaveComponent(quit, text.TextComponent{Text: "X"})
	t.textTool.FontSize().SaveComponent(quit, text.FontSizeComponent{FontSize: 25})
	t.textTool.TextAlign().SaveComponent(quit, text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	t.renderTool.Color().SaveComponent(quit, render.NewColor(mgl32.Vec4{1, 0, 0, 1}))
	t.renderTool.Mesh().SaveComponent(quit, render.NewMesh(t.gameAssets.SquareMesh))
	t.renderTool.Texture().SaveComponent(quit, render.NewTexture(t.gameAssets.Tiles.Water))
	t.pipelineArray.SaveComponent(quit, genericrenderer.PipelineComponent{})

	t.leftClickArray.SaveComponent(quit, inputs.NewMouseLeftClick(ui.HideUiEvent{}))
	t.keepSelectedArray.SaveComponent(quit, inputs.KeepSelectedComponent{})
	t.colliderArray.SaveComponent(quit, collider.NewCollider(t.gameAssets.SquareCollider))

	// child wrapper
	childWrapper := t.world.NewEntity()
	t.hierarchyTool.SetParent(childWrapper, menu)
	t.groupInheritArray.SaveComponent(childWrapper, groups.InheritGroupsComponent{})
	t.transformTool.Parent().SaveComponent(childWrapper, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))

	t.menu = menu
	t.childWrapper = childWrapper
	t.isInitialized = true

	events.Listen(t.world.EventsBuilder(), func(e ui.HideUiEvent) {
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

	for _, child := range t.hierarchyTool.Children(t.childWrapper).GetIndices() {
		t.world.RemoveEntity(child)
	}
	return nil
}

func (t tool) Show() ecs.EntityID {
	t.logger.Warn(t.ResetChildWrapper())
	if !t.active {
		t.active = true
		t.animationArray.SaveComponent(t.menu, animation.NewAnimationComponent(t.showAnimation, t.animationDuration))
	}
	return t.childWrapper
}

func (t tool) Hide() {
	if t.active {
		t.active = false
		t.animationArray.SaveComponent(t.menu, animation.NewAnimationComponent(t.hideAnimation, t.animationDuration))
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
