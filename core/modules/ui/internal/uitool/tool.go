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
	cameraTool    camera.Tool
	transformTool transform.Tool
	tileTool      tile.Tool
	textTool      text.Tool
	renderTool    render.Tool
	hierarchyTool hierarchy.Tool

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
	cameraToolFactory ecs.ToolFactory[camera.Tool],
	transformToolFactory ecs.ToolFactory[transform.Tool],
	tileToolFactory ecs.ToolFactory[tile.Tool],
	textToolFactory ecs.ToolFactory[text.Tool],
	renderToolFactory ecs.ToolFactory[render.Tool],
	hierarchyToolFactory ecs.ToolFactory[hierarchy.Tool],
) tool {
	t := tool{
		changing: &changing{},

		animationDuration: animationDuration,
		showAnimation:     showAnimation,
		hideAnimation:     hideAnimation,

		world:         world,
		gameAssets:    gameAssets,
		logger:        logger,
		cameraTool:    cameraToolFactory.Build(world),
		transformTool: transformToolFactory.Build(world),
		tileTool:      tileToolFactory.Build(world),
		textTool:      textToolFactory.Build(world),
		renderTool:    renderToolFactory.Build(world),
		hierarchyTool: hierarchyToolFactory.Build(world),

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
	transactions := []ecs.AnyComponentsArrayTransaction{}

	transformTransaction := t.transformTool.Transaction()
	transactions = append(transactions, transformTransaction.Transactions()...)

	renderTransaction := t.renderTool.Transaction()
	transactions = append(transactions, renderTransaction.Transactions()...)

	textTransaction := t.textTool.Transaction()
	transactions = append(transactions, textTransaction.Transactions()...)

	hierarchyTransaction := t.hierarchyTool.Transaction()
	transactions = append(transactions, hierarchyTransaction.Transactions()...)

	pipelineTransaction := t.pipelineArray.Transaction()
	groupInheritTransaction := t.groupInheritArray.Transaction()

	leftClickTransaction := t.leftClickArray.Transaction()
	keepSelectedTransaction := t.keepSelectedArray.Transaction()
	colliderTransaction := t.colliderArray.Transaction()

	transactions = append(transactions,
		pipelineTransaction, groupInheritTransaction,
		leftClickTransaction, keepSelectedTransaction, colliderTransaction,
	)

	// objects
	// menu
	menu := t.world.NewEntity()
	hierarchyTransaction.GetObject(menu).Parent().Set(hierarchy.NewParent(camera))
	menuTransform := transformTransaction.GetObject(menu)
	menuTransform.Parent().Set(transform.NewParent(transform.RelativePos | transform.RelativeSizeXY))
	menuTransform.ParentPivotPoint().Set(transform.NewParentPivotPoint(1, 1, .5))
	menuTransform.Pos().Set(transform.NewPos(0, 0, 1))
	menuTransform.Size().Set(transform.NewSize(.2, 1, 1))
	menuTransform.PivotPoint().Set(transform.NewPivotPoint(0, 1, .5))

	menuRender := renderTransaction.GetObject(menu)
	menuRender.Color().Set(render.NewColor(mgl32.Vec4{1, 1, 1, .5}))
	menuRender.Mesh().Set(render.NewMesh(t.gameAssets.SquareMesh))
	menuRender.Texture().Set(render.NewTexture(t.gameAssets.Tiles.Water))
	pipelineTransaction.SaveComponent(menu, genericrenderer.PipelineComponent{})

	groupInheritTransaction.SaveComponent(menu, groups.InheritGroupsComponent{})
	colliderTransaction.SaveComponent(menu, collider.NewCollider(t.gameAssets.SquareCollider))
	keepSelectedTransaction.SaveComponent(menu, inputs.KeepSelectedComponent{})

	// quit btn
	quit := t.world.NewEntity()

	hierarchyTransaction.GetObject(quit).Parent().Set(hierarchy.NewParent(menu))
	groupInheritTransaction.SaveComponent(quit, groups.InheritGroupsComponent{})

	quitTransform := transformTransaction.GetObject(quit)
	quitTransform.Parent().Set(transform.NewParent(transform.RelativePos))
	quitTransform.ParentPivotPoint().Set(transform.NewParentPivotPoint(1, 1, 1))
	quitTransform.Size().Set(transform.NewSize(25, 25, 1))
	quitTransform.PivotPoint().Set(transform.NewPivotPoint(1, 1, 0))

	quitText := textTransaction.GetObject(quit)
	quitText.Text().Set(text.TextComponent{Text: "X"})
	quitText.FontSize().Set(text.FontSizeComponent{FontSize: 25})
	quitText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	quitRender := renderTransaction.GetObject(quit)
	quitRender.Color().Set(render.NewColor(mgl32.Vec4{1, 0, 0, 1}))
	quitRender.Mesh().Set(render.NewMesh(t.gameAssets.SquareMesh))
	quitRender.Texture().Set(render.NewTexture(t.gameAssets.Tiles.Water))
	pipelineTransaction.SaveComponent(quit, genericrenderer.PipelineComponent{})

	leftClickTransaction.SaveComponent(quit, inputs.NewMouseLeftClick(ui.HideUiEvent{}))
	keepSelectedTransaction.SaveComponent(quit, inputs.KeepSelectedComponent{})
	colliderTransaction.SaveComponent(quit, collider.NewCollider(t.gameAssets.SquareCollider))

	// child wrapper
	childWrapper := t.world.NewEntity()
	hierarchyTransaction.GetObject(childWrapper).Parent().Set(hierarchy.NewParent(menu))
	childrenContainerTransform := transformTransaction.GetObject(childWrapper)
	childrenContainerTransform.Parent().Set(transform.NewParent(transform.RelativePos | transform.RelativeSizeXY))

	t.menu = menu
	t.childWrapper = childWrapper
	t.isInitialized = true

	events.Listen(t.world.EventsBuilder(), func(e ui.HideUiEvent) {
		t.Hide()
	})
	return ecs.FlushMany(transactions...)
}

func (t tool) ResetChildWrapper() error {
	if !t.isInitialized {
		err := t.Init()
		if err != nil {
			return err
		}
	}
	transactions := []ecs.AnyComponentsArrayTransaction{}

	groupsInheritTransaction := t.groupInheritArray.Transaction()
	transactions = append(transactions, groupsInheritTransaction)

	transformTransaction := t.transformTool.Transaction()
	transactions = append(transactions, transformTransaction.Transactions()...)

	hierarchyTransaction := t.hierarchyTool.Transaction()
	transactions = append(transactions, hierarchyTransaction.Transactions()...)

	t.world.RemoveEntity(t.childWrapper)

	childrenContainer := t.world.NewEntity()
	hierarchyTransaction.GetObject(childrenContainer).Parent().Set(hierarchy.NewParent(t.menu))
	childrenContainerTransform := transformTransaction.GetObject(childrenContainer)
	childrenContainerTransform.Parent().Set(transform.NewParent(transform.RelativePos | transform.RelativeSizeXY))
	groupsInheritTransaction.SaveComponent(childrenContainer, groups.InheritGroupsComponent{})

	t.childWrapper = childrenContainer

	return ecs.FlushMany(transactions...)
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
