package internal

import (
	gameassets "core/assets"
	"core/modules/settings"
	"core/modules/ui"
	gamescenes "core/scenes"
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

func NewSystem(
	logger logger.Logger,
	transformToolFactory ecs.ToolFactory[transform.Tool],
	renderToolFactory ecs.ToolFactory[render.Tool],
	uiToolFactory ecs.ToolFactory[ui.Tool],
	textToolFactory ecs.ToolFactory[text.Tool],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(world ecs.World) error {
		uiTool := uiToolFactory.Build(world)

		transformTool := transformToolFactory.Build(world)
		renderTool := renderToolFactory.Build(world)
		textTool := textToolFactory.Build(world)

		hierarchyArray := ecs.GetComponentsArray[hierarchy.ParentComponent](world)
		inheritGroupsArray := ecs.GetComponentsArray[groups.InheritGroupsComponent](world)

		colliderArray := ecs.GetComponentsArray[collider.ColliderComponent](world)
		leftClickArray := ecs.GetComponentsArray[inputs.MouseLeftClickComponent](world)
		keepSelectedArray := ecs.GetComponentsArray[inputs.KeepSelectedComponent](world)

		pipelineArray := ecs.GetComponentsArray[genericrenderer.PipelineComponent](world)

		events.Listen(world.EventsBuilder(), func(e settings.EnterSettingsEvent) {
			parent := uiTool.Show()

			// transactions
			hierarchyTransaction := hierarchyArray.Transaction()
			inheritGroupsTransaction := inheritGroupsArray.Transaction()

			colliderTransaction := colliderArray.Transaction()
			leftClickTransaction := leftClickArray.Transaction()
			keepSelectedTransaction := keepSelectedArray.Transaction()

			pipelineTransaction := pipelineArray.Transaction()

			transactions := []ecs.AnyComponentsArrayTransaction{
				hierarchyTransaction,
				inheritGroupsTransaction,
				colliderTransaction,
				leftClickTransaction,
				keepSelectedTransaction,
				pipelineTransaction,
			}

			transformTransaction := transformTool.Transaction()
			transactions = append(transactions, transformTransaction.Transactions()...)

			renderTransaction := renderTool.Transaction()
			transactions = append(transactions, renderTransaction.Transactions()...)

			textTransaction := textTool.Transaction()
			transactions = append(transactions, textTransaction.Transactions()...)

			// render
			// collider
			// click

			// changes
			textEntity := world.NewEntity()
			hierarchyTransaction.SaveComponent(textEntity, hierarchy.NewParent(parent))
			inheritGroupsTransaction.SaveComponent(textEntity, groups.InheritGroupsComponent{})

			labelTransform := transformTransaction.GetObject(textEntity)
			labelTransform.Parent().Set(transform.NewParent(transform.RelativePos))
			labelTransform.Pos().Set(transform.NewPos(0, -25, 0))
			labelTransform.Size().Set(transform.NewSize(160, 50, 1))

			labelText := textTransaction.GetObject(textEntity)
			labelText.Text().Set(text.TextComponent{Text: "SETTINGS"})
			labelText.FontSize().Set(text.FontSizeComponent{FontSize: 25})
			labelText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})

			//

			btnEntity := world.NewEntity()
			hierarchyTransaction.SaveComponent(btnEntity, hierarchy.NewParent(parent))
			inheritGroupsTransaction.SaveComponent(btnEntity, groups.InheritGroupsComponent{})

			btnTransform := transformTransaction.GetObject(btnEntity)
			btnTransform.Parent().Set(transform.NewParent(transform.RelativePos))
			btnTransform.Pos().Set(transform.NewPos(0, 25, 0))
			btnTransform.Size().Set(transform.NewSize(160, 50, 1))

			btnRender := renderTransaction.GetObject(btnEntity)
			btnRender.Mesh().Set(render.NewMesh(gameassets.SquareMesh))
			btnRender.Texture().Set(render.NewTexture(gameassets.GroundTileTextureID))
			pipelineTransaction.SaveComponent(btnEntity, genericrenderer.PipelineComponent{})

			btnText := textTransaction.GetObject(btnEntity)
			btnText.Text().Set(text.TextComponent{Text: "QUIT"})
			btnText.FontSize().Set(text.FontSizeComponent{FontSize: 25})
			btnText.TextAlign().Set(text.TextAlignComponent{Vertical: .5, Horizontal: .5})

			leftClickTransaction.SaveComponent(btnEntity, inputs.NewMouseLeftClick(scenes.NewChangeSceneEvent(gamescenes.MenuID)))
			keepSelectedTransaction.SaveComponent(btnEntity, inputs.KeepSelectedComponent{})
			colliderTransaction.SaveComponent(btnEntity, collider.NewCollider(gameassets.SquareColliderID))

			// flush
			logger.Warn(ecs.FlushMany(transactions...))
		})
		return nil
	})
}
