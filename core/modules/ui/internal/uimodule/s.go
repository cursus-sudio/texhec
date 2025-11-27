package uimodule

import (
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
	"time"

	"github.com/ogiusek/events"
)

//

type SelectedTile struct {
	Tile tile.PosComponent
}

// type state any

//

type Option struct{}
type SelectOptionEvent struct {
	Entity ecs.EntityID
}

//

// what is need:
// - state history
// - current state

// current state:
// - all buttons and text on the screen

type menuData struct {
	menu              ecs.EntityID
	quit              ecs.EntityID
	childrenContainer ecs.EntityID
	visible           bool
}

type uiSys struct {
	// scenessys.NewChangeSceneEvent(gamescenes.MenuID)
	world  ecs.World
	logger logger.Logger

	cameraTool    camera.Tool
	transformTool transform.Tool
	tileTool      tile.Tool
	textTool      text.Tool
	renderTool    render.Tool
	hierarchyTool hierarchy.Tool

	animationArray    ecs.ComponentsArray[animation.AnimationComponent]
	uiCameraArray     ecs.ComponentsArray[ui.UiCameraComponent]
	groupArray        ecs.ComponentsArray[groups.GroupsComponent]
	groupInheritArray ecs.ComponentsArray[groups.InheritGroupsComponent]
	tilePosArray      ecs.ComponentsArray[tile.PosComponent]
	pipelineArray     ecs.ComponentsArray[genericrenderer.PipelineComponent]
	leftClickArray    ecs.ComponentsArray[inputs.MouseLeftClickComponent]
	keepSelectedArray ecs.ComponentsArray[inputs.KeepSelectedComponent]
	colliderArray     ecs.ComponentsArray[collider.ColliderComponent]

	transformTransaction transform.Transaction
	textTransaction      text.Transaction
	renderTransaction    render.Transaction
	hierarchyTransaction hierarchy.Transaction

	animationTransaction    ecs.ComponentsArrayTransaction[animation.AnimationComponent]
	uiCameraTransaction     ecs.ComponentsArrayTransaction[ui.UiCameraComponent]
	groupInheritTransaction ecs.ComponentsArrayTransaction[groups.InheritGroupsComponent]
	tilePosTransaction      ecs.ComponentsArrayTransaction[tile.PosComponent]
	pipelineTransaction     ecs.ComponentsArrayTransaction[genericrenderer.PipelineComponent]
	leftClickTransaction    ecs.ComponentsArrayTransaction[inputs.MouseLeftClickComponent]
	keepSelectedTransaction ecs.ComponentsArrayTransaction[inputs.KeepSelectedComponent]
	colliderTransaction     ecs.ComponentsArrayTransaction[collider.ColliderComponent]

	maxTileDepth      tile.Layer
	animationDuration time.Duration
	menu              *menuData
	currentState      *SelectedTile
}

func NewSystem(
	logger logger.Logger,
	cameraToolFactory ecs.ToolFactory[camera.Tool],
	transformToolFactory ecs.ToolFactory[transform.Tool],
	tileToolFactory ecs.ToolFactory[tile.Tool],
	textToolFactory ecs.ToolFactory[text.Tool],
	renderToolFactory ecs.ToolFactory[render.Tool],
	hierarchyToolFactory ecs.ToolFactory[hierarchy.Tool],
	maxTileDepth tile.Layer,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(world ecs.World) error {
		s := &uiSys{
			world:             world,
			logger:            logger,
			cameraTool:        cameraToolFactory.Build(world),
			transformTool:     transformToolFactory.Build(world),
			tileTool:          tileToolFactory.Build(world),
			textTool:          textToolFactory.Build(world),
			renderTool:        renderToolFactory.Build(world),
			hierarchyTool:     hierarchyToolFactory.Build(world),
			animationArray:    ecs.GetComponentsArray[animation.AnimationComponent](world),
			uiCameraArray:     ecs.GetComponentsArray[ui.UiCameraComponent](world),
			groupArray:        ecs.GetComponentsArray[groups.GroupsComponent](world),
			groupInheritArray: ecs.GetComponentsArray[groups.InheritGroupsComponent](world),
			tilePosArray:      ecs.GetComponentsArray[tile.PosComponent](world),
			pipelineArray:     ecs.GetComponentsArray[genericrenderer.PipelineComponent](world),
			leftClickArray:    ecs.GetComponentsArray[inputs.MouseLeftClickComponent](world),
			keepSelectedArray: ecs.GetComponentsArray[inputs.KeepSelectedComponent](world),
			colliderArray:     ecs.GetComponentsArray[collider.ColliderComponent](world),
			animationDuration: time.Millisecond * 100,
			maxTileDepth:      maxTileDepth,
		}
		events.ListenE(world.EventsBuilder(), s.Listen1)
		events.ListenE(world.EventsBuilder(), s.Listen2)
		events.ListenE(world.EventsBuilder(), s.Listen3)
		events.ListenE(world.EventsBuilder(), s.Listen4)
		return nil
	})
}

func (s *uiSys) Flush() error {
	transactions := []ecs.AnyComponentsArrayTransaction{}
	transactions = append(transactions,
		s.animationTransaction,
		s.groupInheritTransaction,
		s.pipelineTransaction,
		s.leftClickTransaction,
		s.keepSelectedTransaction,
		s.colliderTransaction,
	)
	transactions = append(transactions, s.transformTransaction.Transactions()...)
	transactions = append(transactions, s.renderTransaction.Transactions()...)
	transactions = append(transactions, s.textTransaction.Transactions()...)
	transactions = append(transactions, s.hierarchyTransaction.Transactions()...)
	if err := ecs.FlushMany(transactions...); err != nil {
		return err
	}
	return nil
}

func (ui *uiSys) Listen1(e ui.UnselectEvent) error {
	// regress tile or set selected to nil
	ui.currentState = nil
	return ui.Render()
}

func (ui *uiSys) Listen2(e SelectOptionEvent) error {
	// TODO define options
	// posComponent, err := ui.tilePosArray.GetComponent(e.Entity)
	// if err != nil {
	// 	return err
	// }
	// ui.currentState = &SelectedTile{posComponent}
	// ui.Render()
	return nil
}

func (ui *uiSys) Listen3(e tile.TileClickEvent) error {
	posComponent, err := ui.tilePosArray.GetComponent(e.Tile)
	if err != nil {
		return err
	}
	if ui.currentState != nil &&
		ui.currentState.Tile.X == posComponent.X &&
		ui.currentState.Tile.Y == posComponent.Y {
		for i := ui.currentState.Tile.Layer + 1; i < ui.currentState.Tile.Layer+ui.maxTileDepth; i++ {
			layer := i % ui.maxTileDepth
			pos := tile.NewPos(posComponent.X, posComponent.Y, layer)
			_, ok := ui.tileTool.TilePos().Get(pos)
			if ok {
				posComponent = pos
				break
			}
		}
	}
	ui.currentState = &SelectedTile{posComponent}
	return ui.Render()
}

func (ui *uiSys) Listen4(e ui.SettingsEvent) error {
	ui.logger.Info("settings")
	return nil
}
