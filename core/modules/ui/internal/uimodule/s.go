package uimodule

import (
	gameassets "core/assets"
	"core/modules/tile"
	"core/modules/ui"
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
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

type UiElementComponent struct{}

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

type uiSys struct {
	world  ecs.World
	logger logger.Logger

	cameraTool    camera.Tool
	transformTool transform.Tool
	tileTool      tile.Tool
	textTool      text.Tool

	uiCameraArray     ecs.ComponentsArray[ui.UiCameraComponent]
	groupsArray       ecs.ComponentsArray[groups.GroupsComponent]
	uiElementArray    ecs.ComponentsArray[UiElementComponent]
	tilePosArray      ecs.ComponentsArray[tile.PosComponent]
	colorArray        ecs.ComponentsArray[render.ColorComponent]
	meshArray         ecs.ComponentsArray[render.MeshComponent]
	textureArray      ecs.ComponentsArray[render.TextureComponent]
	pipelineArray     ecs.ComponentsArray[genericrenderer.PipelineComponent]
	leftClickArray    ecs.ComponentsArray[inputs.MouseLeftClickComponent]
	keepSelectedArray ecs.ComponentsArray[inputs.KeepSelectedComponent]
	colliderArray     ecs.ComponentsArray[collider.ColliderComponent]

	maxTileDepth tile.Layer
	currentState *SelectedTile
}

func NewSystem(
	logger logger.Logger,
	cameraToolFactory ecs.ToolFactory[camera.Tool],
	transformToolFactory ecs.ToolFactory[transform.Tool],
	tileToolFactory ecs.ToolFactory[tile.Tool],
	textToolFactory ecs.ToolFactory[text.Tool],
	maxTileDepth tile.Layer,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(world ecs.World) error {
		s := &uiSys{
			world,
			logger,
			cameraToolFactory.Build(world),
			transformToolFactory.Build(world),
			tileToolFactory.Build(world),
			textToolFactory.Build(world),
			ecs.GetComponentsArray[ui.UiCameraComponent](world),
			ecs.GetComponentsArray[groups.GroupsComponent](world),
			ecs.GetComponentsArray[UiElementComponent](world),
			ecs.GetComponentsArray[tile.PosComponent](world),
			ecs.GetComponentsArray[render.ColorComponent](world),
			ecs.GetComponentsArray[render.MeshComponent](world),
			ecs.GetComponentsArray[render.TextureComponent](world),
			ecs.GetComponentsArray[genericrenderer.PipelineComponent](world),
			ecs.GetComponentsArray[inputs.MouseLeftClickComponent](world),
			ecs.GetComponentsArray[inputs.KeepSelectedComponent](world),
			ecs.GetComponentsArray[collider.ColliderComponent](world),
			maxTileDepth,
			nil,
		}
		events.ListenE(world.EventsBuilder(), s.Listen1)
		events.ListenE(world.EventsBuilder(), s.Listen2)
		events.ListenE(world.EventsBuilder(), s.Listen3)
		return nil
	})
}

func (s *uiSys) Render() error {
	for _, entity := range s.uiElementArray.GetEntities() {
		s.world.RemoveEntity(entity)
	}

	if s.currentState == nil {
		return nil
	}

	state := *s.currentState

	transformTransaction := s.transformTool.Transaction()
	textTransaction := s.textTool.Transaction()

	groupsTransaction := s.groupsArray.Transaction()
	uiElementsTransaction := s.uiElementArray.Transaction()
	colorTransaction := s.colorArray.Transaction()
	meshTransaction := s.meshArray.Transaction()
	textureTransaction := s.textureArray.Transaction()
	genericTransaction := s.pipelineArray.Transaction()
	leftClickTransaction := s.leftClickArray.Transaction()
	keepSelectedTransaction := s.keepSelectedArray.Transaction()
	colliderTransaction := s.colliderArray.Transaction()
	for _, camera := range s.uiCameraArray.GetEntities() {
		groups, err := s.groupsArray.GetComponent(camera)
		if err != nil {
			s.logger.Warn(err)
			continue
		}
		menu := s.world.NewEntity()
		menuTransform := transformTransaction.GetObject(menu)
		menuTransform.Parent().Set(transform.NewParent(camera, transform.RelativePos))
		menuTransform.ParentPivotPoint().Set(transform.NewParentPivotPoint(.5, 0, .5))
		menuTransform.Size().Set(transform.NewSize(500, 100, 1))
		menuTransform.PivotPoint().Set(transform.NewPivotPoint(.5, 0, .5))

		uiElementsTransaction.SaveComponent(menu, UiElementComponent{})
		groupsTransaction.SaveComponent(menu, groups)

		menuText := textTransaction.GetObject(menu)
		menuText.Text().Set(text.TextComponent{Text: fmt.Sprintf("pos is %v", state.Tile)})
		menuText.FontSize().Set(text.FontSizeComponent{FontSize: 32})
		menuText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})

		colorTransaction.SaveComponent(menu, render.NewColor(mgl32.Vec4{1, 1, 1, .5}))
		meshTransaction.SaveComponent(menu, render.NewMesh(gameassets.SquareMesh))
		textureTransaction.SaveComponent(menu, render.NewTexture(gameassets.WaterTileTextureID))
		genericTransaction.SaveComponent(menu, genericrenderer.PipelineComponent{})

		quit := s.world.NewEntity()
		quitTransform := transformTransaction.GetObject(quit)
		quitTransform.Parent().Set(transform.NewParent(menu, transform.RelativePos))
		quitTransform.ParentPivotPoint().Set(transform.NewParentPivotPoint(1, 1, .5))
		quitTransform.Size().Set(transform.NewSize(25, 25, 2))
		quitTransform.PivotPoint().Set(transform.NewPivotPoint(1, 1, .5))

		uiElementsTransaction.SaveComponent(quit, UiElementComponent{})
		groupsTransaction.SaveComponent(quit, groups)

		quitText := textTransaction.GetObject(quit)
		quitText.Text().Set(text.TextComponent{Text: "X"})
		quitText.FontSize().Set(text.FontSizeComponent{FontSize: 25})
		quitText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})

		colorTransaction.SaveComponent(quit, render.NewColor(mgl32.Vec4{1, 0, 0, 1}))
		meshTransaction.SaveComponent(quit, render.NewMesh(gameassets.SquareMesh))
		textureTransaction.SaveComponent(quit, render.NewTexture(gameassets.WaterTileTextureID))
		genericTransaction.SaveComponent(quit, genericrenderer.PipelineComponent{})

		leftClickTransaction.SaveComponent(quit, inputs.NewMouseLeftClick(ui.UnselectEvent{}))
		keepSelectedTransaction.SaveComponent(quit, inputs.KeepSelectedComponent{})
		colliderTransaction.SaveComponent(quit, collider.NewCollider(gameassets.SquareColliderID))
	}
	transactions := []ecs.AnyComponentsArrayTransaction{
		groupsTransaction,
		uiElementsTransaction,
		colorTransaction,
		meshTransaction,
		textureTransaction,
		genericTransaction,
		leftClickTransaction,
		keepSelectedTransaction,
		colliderTransaction,
	}
	transactions = append(transactions, transformTransaction.Transactions()...)
	transactions = append(transactions, textTransaction.Transactions()...)
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
