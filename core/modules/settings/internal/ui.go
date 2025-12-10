package internal

import (
	gameassets "core/assets"
	"core/modules/settings"
	"core/modules/ui"
	gamescenes "core/scenes"
	"engine/modules/audio"
	"engine/modules/collider"
	"engine/modules/genericrenderer"
	"engine/modules/groups"
	"engine/modules/hierarchy"
	"engine/modules/inputs"
	"engine/modules/render"
	"engine/modules/scenes"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/events"
)

// 1. settings text
// 2. quit button

type system struct {
	gameAssets gameassets.GameAssets

	world ecs.World

	uiTool        ui.Tool
	transformTool transform.Tool
	renderTool    render.Tool
	textTool      text.Tool

	hierarchyArray     ecs.ComponentsArray[hierarchy.ParentComponent]
	inheritGroupsArray ecs.ComponentsArray[groups.InheritGroupsComponent]

	colliderArray     ecs.ComponentsArray[collider.ColliderComponent]
	leftClickArray    ecs.ComponentsArray[inputs.MouseLeftClickComponent]
	keepSelectedArray ecs.ComponentsArray[inputs.KeepSelectedComponent]

	pipelineArray ecs.ComponentsArray[genericrenderer.PipelineComponent]
}

func NewSystem(
	logger logger.Logger,
	gameAssets gameassets.GameAssets,
	transformToolFactory ecs.ToolFactory[transform.Tool],
	renderToolFactory ecs.ToolFactory[render.Tool],
	uiToolFactory ecs.ToolFactory[ui.Tool],
	textToolFactory ecs.ToolFactory[text.Tool],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(world ecs.World) error {
		s := system{
			gameAssets,

			world,
			uiToolFactory.Build(world),
			transformToolFactory.Build(world),
			renderToolFactory.Build(world),
			textToolFactory.Build(world),

			ecs.GetComponentsArray[hierarchy.ParentComponent](world),
			ecs.GetComponentsArray[groups.InheritGroupsComponent](world),
			ecs.GetComponentsArray[collider.ColliderComponent](world),
			ecs.GetComponentsArray[inputs.MouseLeftClickComponent](world),
			ecs.GetComponentsArray[inputs.KeepSelectedComponent](world),
			ecs.GetComponentsArray[genericrenderer.PipelineComponent](world),
		}

		events.ListenE(world.EventsBuilder(), func(event settings.EnterSettingsForParentEvent) error {
			return s.Render(event.Parent)
		})
		events.Listen(world.EventsBuilder(), func(settings.EnterSettingsEvent) {
			event := settings.EnterSettingsForParentEvent{
				Parent: s.uiTool.Show(),
			}
			events.Emit(world.Events(), event)
		})

		return nil
	})
}

func (s system) Render(parent ecs.EntityID) error {
	// transactions
	hierarchyTransaction := s.hierarchyArray.Transaction()
	inheritGroupsTransaction := s.inheritGroupsArray.Transaction()

	colliderTransaction := s.colliderArray.Transaction()
	leftClickTransaction := s.leftClickArray.Transaction()
	keepSelectedTransaction := s.keepSelectedArray.Transaction()

	pipelineTransaction := s.pipelineArray.Transaction()

	transactions := []ecs.AnyComponentsArrayTransaction{
		hierarchyTransaction, inheritGroupsTransaction,
		colliderTransaction, leftClickTransaction, keepSelectedTransaction,
		pipelineTransaction,
	}

	transformTransaction := s.transformTool.Transaction()
	transactions = append(transactions, transformTransaction.Transactions()...)

	renderTransaction := s.renderTool.Transaction()
	transactions = append(transactions, renderTransaction.Transactions()...)

	textTransaction := s.textTool.Transaction()
	transactions = append(transactions, textTransaction.Transactions()...)

	// render
	// collider
	// click

	// changes
	textEntity := s.world.NewEntity()
	hierarchyTransaction.SaveComponent(textEntity, hierarchy.NewParent(parent))
	inheritGroupsTransaction.SaveComponent(textEntity, groups.InheritGroupsComponent{})

	labelTransform := transformTransaction.GetObject(textEntity)
	labelTransform.Parent().Set(transform.NewParent(transform.RelativePos | transform.RelativeSizeX))
	labelTransform.ParentPivotPoint().Set(transform.NewParentPivotPoint(.5, .5, 1))
	labelTransform.PivotPoint().Set(transform.NewPivotPoint(.5, .5, 0))
	labelTransform.Pos().Set(transform.NewPos(0, 25, 0))
	labelTransform.Size().Set(transform.NewSize(1, 50, 1))

	labelText := textTransaction.GetObject(textEntity)
	labelText.Text().Set(text.TextComponent{Text: "SETTINGS"})
	labelText.FontSize().Set(text.FontSizeComponent{FontSize: 25})
	labelText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	//

	type Button struct {
		text  string
		event any
	}
	btns := []Button{
		{"SHOT", audio.NewPlayEvent(gamescenes.EffectChannel, s.gameAssets.ExampleAudio)},
		{"QUIT", scenes.NewChangeSceneEvent(gamescenes.MenuID)},
	}

	for i, btn := range btns {
		var height float32 = 50
		var margin float32 = 10
		btnEntity := s.world.NewEntity()
		hierarchyTransaction.SaveComponent(btnEntity, hierarchy.NewParent(parent))
		inheritGroupsTransaction.SaveComponent(btnEntity, groups.InheritGroupsComponent{})

		btnTransform := transformTransaction.GetObject(btnEntity)
		btnTransform.Parent().Set(transform.NewParent(transform.RelativePos | transform.RelativeSizeX))
		btnTransform.ParentPivotPoint().Set(transform.NewParentPivotPoint(.5, .5, 1))
		btnTransform.PivotPoint().Set(transform.NewPivotPoint(.5, .5, 0))
		btnTransform.MaxSize().Set(transform.NewMaxSize(300, 0, 0))
		btnTransform.Pos().Set(transform.NewPos(0, float32(-i)*(height+margin)-25, 0))
		btnTransform.Size().Set(transform.NewSize(1, height, 1))

		btnRender := renderTransaction.GetObject(btnEntity)
		btnRender.Mesh().Set(render.NewMesh(s.gameAssets.SquareMesh))
		btnRender.Texture().Set(render.NewTexture(s.gameAssets.Hud.Btn))
		pipelineTransaction.SaveComponent(btnEntity, genericrenderer.PipelineComponent{})

		btnText := textTransaction.GetObject(btnEntity)
		btnText.Text().Set(text.TextComponent{Text: btn.text})
		btnText.FontSize().Set(text.FontSizeComponent{FontSize: 25})
		btnText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})

		leftClickTransaction.SaveComponent(btnEntity, inputs.NewMouseLeftClick(btn.event))
		keepSelectedTransaction.SaveComponent(btnEntity, inputs.KeepSelectedComponent{})
		colliderTransaction.SaveComponent(btnEntity, collider.NewCollider(s.gameAssets.SquareCollider))
	}

	// flush
	return ecs.FlushMany(transactions...)
}
