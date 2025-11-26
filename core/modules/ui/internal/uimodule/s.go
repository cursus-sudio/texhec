package uimodule

import (
	gameassets "core/assets"
	"core/modules/tile"
	"core/modules/ui"
	"engine/modules/animation"
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/genericrenderer"
	"engine/modules/groups"
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

	animationArray    ecs.ComponentsArray[animation.AnimationComponent]
	uiCameraArray     ecs.ComponentsArray[ui.UiCameraComponent]
	groupsArray       ecs.ComponentsArray[groups.GroupsComponent]
	tilePosArray      ecs.ComponentsArray[tile.PosComponent]
	pipelineArray     ecs.ComponentsArray[genericrenderer.PipelineComponent]
	leftClickArray    ecs.ComponentsArray[inputs.MouseLeftClickComponent]
	keepSelectedArray ecs.ComponentsArray[inputs.KeepSelectedComponent]
	colliderArray     ecs.ComponentsArray[collider.ColliderComponent]

	transformTransaction transform.Transaction
	textTransaction      text.Transaction
	renderTransaction    render.Transaction

	animationTransaction    ecs.ComponentsArrayTransaction[animation.AnimationComponent]
	uiCameraTransaction     ecs.ComponentsArrayTransaction[ui.UiCameraComponent]
	groupsTransaction       ecs.ComponentsArrayTransaction[groups.GroupsComponent]
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
			animationArray:    ecs.GetComponentsArray[animation.AnimationComponent](world),
			uiCameraArray:     ecs.GetComponentsArray[ui.UiCameraComponent](world),
			groupsArray:       ecs.GetComponentsArray[groups.GroupsComponent](world),
			tilePosArray:      ecs.GetComponentsArray[tile.PosComponent](world),
			pipelineArray:     ecs.GetComponentsArray[genericrenderer.PipelineComponent](world),
			leftClickArray:    ecs.GetComponentsArray[inputs.MouseLeftClickComponent](world),
			keepSelectedArray: ecs.GetComponentsArray[inputs.KeepSelectedComponent](world),
			colliderArray:     ecs.GetComponentsArray[collider.ColliderComponent](world),
			animationDuration: time.Millisecond * 200,
			maxTileDepth:      maxTileDepth,
		}
		events.ListenE(world.EventsBuilder(), s.Listen1)
		events.ListenE(world.EventsBuilder(), s.Listen2)
		events.ListenE(world.EventsBuilder(), s.Listen3)
		return nil
	})
}

func (s *uiSys) EnsureInit() {
	if s.menu != nil {
		return
	}
	// start transactions
	s.transformTransaction = s.transformTool.Transaction()
	s.renderTransaction = s.renderTool.Transaction()
	s.textTransaction = s.textTool.Transaction()

	s.animationTransaction = s.animationArray.Transaction()
	s.groupsTransaction = s.groupsArray.Transaction()
	s.pipelineTransaction = s.pipelineArray.Transaction()
	s.colliderTransaction = s.colliderArray.Transaction()
	s.leftClickTransaction = s.leftClickArray.Transaction()
	s.keepSelectedTransaction = s.keepSelectedArray.Transaction()

	// initialize main entities
	menu := s.world.NewEntity()
	menuTransform := s.transformTransaction.GetObject(menu)
	menuTransform.ParentPivotPoint().Set(transform.NewParentPivotPoint(1, 1, .5))
	menuTransform.Size().Set(transform.NewSize(.2, 1, 1))
	menuTransform.PivotPoint().Set(transform.NewPivotPoint(1, 1, .5))

	menuRender := s.renderTransaction.GetObject(menu)
	menuRender.Color().Set(render.NewColor(mgl32.Vec4{1, 1, 1, .5}))
	menuRender.Mesh().Set(render.NewMesh(gameassets.SquareMesh))
	menuRender.Texture().Set(render.NewTexture(gameassets.WaterTileTextureID))
	s.pipelineTransaction.SaveComponent(menu, genericrenderer.PipelineComponent{})

	s.colliderTransaction.SaveComponent(menu, collider.NewCollider(gameassets.SquareColliderID))

	menuText := s.textTransaction.GetObject(menu)
	menuText.FontSize().Set(text.FontSizeComponent{FontSize: 32})
	menuText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	childrenContainer := s.world.NewEntity()
	childrenContainerTransform := s.transformTransaction.GetObject(childrenContainer)
	childrenContainerTransform.Parent().Set(transform.NewParent(menu, transform.RelativePos|transform.RelativeSize))

	quit := s.world.NewEntity()
	quitTransform := s.transformTransaction.GetObject(quit)
	quitTransform.Parent().Set(transform.NewParent(menu, transform.RelativePos))
	quitTransform.ParentPivotPoint().Set(transform.NewParentPivotPoint(1, 1, .5))
	quitTransform.Size().Set(transform.NewSize(25, 25, 2))
	quitTransform.PivotPoint().Set(transform.NewPivotPoint(1, 1, .5))

	quitText := s.textTransaction.GetObject(quit)
	quitText.Text().Set(text.TextComponent{Text: "X"})
	quitText.FontSize().Set(text.FontSizeComponent{FontSize: 25})
	quitText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	quitRender := s.renderTransaction.GetObject(quit)
	quitRender.Color().Set(render.NewColor(mgl32.Vec4{1, 0, 0, 1}))
	quitRender.Mesh().Set(render.NewMesh(gameassets.SquareMesh))
	quitRender.Texture().Set(render.NewTexture(gameassets.WaterTileTextureID))
	s.pipelineTransaction.SaveComponent(quit, genericrenderer.PipelineComponent{})

	s.leftClickTransaction.SaveComponent(quit, inputs.NewMouseLeftClick(ui.UnselectEvent{}))
	s.keepSelectedTransaction.SaveComponent(quit, inputs.KeepSelectedComponent{})
	s.colliderTransaction.SaveComponent(quit, collider.NewCollider(gameassets.SquareColliderID))

	s.menu = &menuData{
		menu:              menu,
		childrenContainer: childrenContainer,
		visible:           false,
	}
}

func (s *uiSys) Flush() error {
	transactions := []ecs.AnyComponentsArrayTransaction{
		s.groupsTransaction,
		s.pipelineTransaction,
		s.leftClickTransaction,
		s.keepSelectedTransaction,
		s.colliderTransaction,
	}
	transactions = append(transactions, s.transformTransaction.Transactions()...)
	transactions = append(transactions, s.renderTransaction.Transactions()...)
	transactions = append(transactions, s.textTransaction.Transactions()...)
	if err := ecs.FlushMany(transactions...); err != nil {
		return err
	}
	return nil
}

func (s *uiSys) Render() error {
	s.EnsureInit()
	for _, entity := range s.transformTransaction.
		GetObject(s.menu.childrenContainer).
		FlatChildren().GetIndices() {
		s.world.RemoveEntity(entity)
	}

	if s.currentState == nil {
		if s.menu.visible {
			s.animationTransaction.SaveComponent(s.menu.menu, animation.NewAnimationComponent(gameassets.HideMenuAnimation, s.animationDuration))
			s.menu.visible = false
		}
		// mark parrent as hidden
		return s.Flush()
	}

	state := *s.currentState

	cameras := s.uiCameraArray.GetEntities()
	if len(cameras) != 1 {
		return errors.New("expected one ui camera")
	}
	camera := cameras[0]
	groups, err := s.groupsArray.GetComponent(camera)
	if err != nil {
		return err
	}

	s.menu.visible = true
	s.animationTransaction.SaveComponent(s.menu.menu, animation.NewAnimationComponent(gameassets.ShowMenuAnimation, s.animationDuration))

	menu := s.menu.menu
	s.groupsTransaction.SaveComponent(menu, groups)
	s.transformTransaction.GetObject(menu).Parent().Set(transform.NewParent(camera, transform.RelativePos|transform.RelativeSize))

	menuText := s.textTransaction.GetObject(menu)
	menuText.Text().Set(text.TextComponent{Text: fmt.Sprintf("pos is %v", state.Tile)})
	menuText.FontSize().Set(text.FontSizeComponent{FontSize: 32})
	menuText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	return s.Flush()
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
