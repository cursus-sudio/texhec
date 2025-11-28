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
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

type changing struct {
	active       bool
	menu         ecs.EntityID
	childWrapper ecs.EntityID
}

type tool struct {
	*changing

	animationDuration time.Duration
	showAnimation     animation.AnimationID
	hideAnimation     animation.AnimationID

	world         ecs.World
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
		t.logger.Info("hihi fm")
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
		pipelineTransaction,
		groupInheritTransaction,
		leftClickTransaction,
		keepSelectedTransaction,
		colliderTransaction,
	)

	// objects
	// menu
	menu := t.world.NewEntity()
	hierarchyTransaction.GetObject(menu).Parent().Set(hierarchy.NewParent(camera))
	menuTransform := transformTransaction.GetObject(menu)
	menuTransform.Parent().Set(transform.NewParent(transform.RelativePos | transform.RelativeSize))
	menuTransform.ParentPivotPoint().Set(transform.NewParentPivotPoint(1, 1, .5))
	menuTransform.Size().Set(transform.NewSize(.2, 1, 1))
	menuTransform.PivotPoint().Set(transform.NewPivotPoint(0, 1, .5))

	menuRender := renderTransaction.GetObject(menu)
	menuRender.Color().Set(render.NewColor(mgl32.Vec4{1, 1, 1, .5}))
	menuRender.Mesh().Set(render.NewMesh(gameassets.SquareMesh))
	menuRender.Texture().Set(render.NewTexture(gameassets.WaterTileTextureID))
	pipelineTransaction.SaveComponent(menu, genericrenderer.PipelineComponent{})

	groupInheritTransaction.SaveComponent(menu, groups.InheritGroupsComponent{})
	colliderTransaction.SaveComponent(menu, collider.NewCollider(gameassets.SquareColliderID))

	// quit btn
	quit := t.world.NewEntity()
	hierarchyTransaction.GetObject(quit).Parent().Set(hierarchy.NewParent(menu))
	quitTransform := transformTransaction.GetObject(quit)
	quitTransform.Parent().Set(transform.NewParent(transform.RelativePos))
	quitTransform.ParentPivotPoint().Set(transform.NewParentPivotPoint(1, 1, .5))
	quitTransform.Size().Set(transform.NewSize(25, 25, 2))
	quitTransform.PivotPoint().Set(transform.NewPivotPoint(1, 1, .5))

	quitText := textTransaction.GetObject(quit)
	quitText.Text().Set(text.TextComponent{Text: "X"})
	quitText.FontSize().Set(text.FontSizeComponent{FontSize: 25})
	quitText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	quitRender := renderTransaction.GetObject(quit)
	quitRender.Color().Set(render.NewColor(mgl32.Vec4{1, 0, 0, 1}))
	quitRender.Mesh().Set(render.NewMesh(gameassets.SquareMesh))
	quitRender.Texture().Set(render.NewTexture(gameassets.WaterTileTextureID))
	pipelineTransaction.SaveComponent(quit, genericrenderer.PipelineComponent{})
	groupInheritTransaction.SaveComponent(quit, groups.InheritGroupsComponent{})

	leftClickTransaction.SaveComponent(quit, inputs.NewMouseLeftClick(ui.HideUiEvent{}))
	keepSelectedTransaction.SaveComponent(quit, inputs.KeepSelectedComponent{})
	colliderTransaction.SaveComponent(quit, collider.NewCollider(gameassets.SquareColliderID))

	// child wrapper
	childWrapper := t.world.NewEntity()
	hierarchyTransaction.GetObject(childWrapper).Parent().Set(hierarchy.NewParent(menu))
	childrenContainerTransform := transformTransaction.GetObject(childWrapper)
	childrenContainerTransform.Parent().Set(transform.NewParent(transform.RelativePos | transform.RelativeSize))

	t.menu = menu
	t.childWrapper = childWrapper

	events.Listen(t.world.EventsBuilder(), func(e ui.HideUiEvent) {
		t.Hide()
	})
	events.Listen(t.world.EventsBuilder(), func(e ui.SettingsEvent) {
		t.logger.Info("settings")
		p := t.Show()
		textTransaction := t.textTool.Transaction()
		quitText := textTransaction.GetObject(p)
		quitText.Text().Set(text.TextComponent{Text: "SETTINGS"})
		quitText.FontSize().Set(text.FontSizeComponent{FontSize: 25})
		quitText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})
		textTransaction.Flush()
	})
	tilePosArray := ecs.GetComponentsArray[tile.PosComponent](t.world)
	events.Listen(t.world.EventsBuilder(), func(e tile.TileClickEvent) {
		t.logger.Info("tile click")
		pos, err := tilePosArray.GetComponent(e.Tile)
		if err != nil {
			t.logger.Warn(err)
		}
		p := t.Show()
		textTransaction := t.textTool.Transaction()
		quitText := textTransaction.GetObject(p)
		quitText.Text().Set(text.TextComponent{Text: fmt.Sprintf("TILE: %v", pos)})
		quitText.FontSize().Set(text.FontSizeComponent{FontSize: 25})
		quitText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})
		textTransaction.Flush()
	})
	events.Listen(t.world.EventsBuilder(), func(e sdl.MouseButtonEvent) {
		if e.Button != sdl.BUTTON_RIGHT || e.State != sdl.RELEASED {
			return
		}
		events.Emit(t.world.Events(), ui.HideUiEvent{})
	})
	return ecs.FlushMany(transactions...)
}

func (t tool) ResetChildWrapper() error {
	groupsInheritTransaction := t.groupInheritArray.Transaction()
	transactions := []ecs.AnyComponentsArrayTransaction{groupsInheritTransaction}

	transformTransaction := t.transformTool.Transaction()
	transactions = append(transactions, transformTransaction.Transactions()...)

	hierarchyTransaction := t.hierarchyTool.Transaction()
	transactions = append(transactions, hierarchyTransaction.Transactions()...)

	t.world.RemoveEntity(t.childWrapper)

	childrenContainer := t.world.NewEntity()
	hierarchyTransaction.GetObject(childrenContainer).Parent().Set(hierarchy.NewParent(t.menu))
	childrenContainerTransform := transformTransaction.GetObject(childrenContainer)
	childrenContainerTransform.Parent().Set(transform.NewParent(transform.RelativePos | transform.RelativeSize))
	groupsInheritTransaction.SaveComponent(childrenContainer, groups.InheritGroupsComponent{})

	t.childWrapper = childrenContainer

	return ecs.FlushMany(transactions...)
}

func (t tool) Show() ecs.EntityID {
	t.ResetChildWrapper()
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
