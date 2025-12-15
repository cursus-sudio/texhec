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
	"engine/services/assets"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/events"
)

// 1. settings text
// 2. quit button

type system struct {
	assets     assets.Assets
	gameAssets gameassets.GameAssets

	logger logger.Logger
	world  ecs.World

	uiTool        ui.Tool
	transformTool transform.Interface
	renderTool    render.Interface
	textTool      text.Interface

	hierarchyArray     ecs.ComponentsArray[hierarchy.Component]
	groupsArray        ecs.ComponentsArray[groups.GroupsComponent]
	inheritGroupsArray ecs.ComponentsArray[groups.InheritGroupsComponent]

	colliderArray     ecs.ComponentsArray[collider.ColliderComponent]
	leftClickArray    ecs.ComponentsArray[inputs.MouseLeftClickComponent]
	keepSelectedArray ecs.ComponentsArray[inputs.KeepSelectedComponent]

	pipelineArray ecs.ComponentsArray[genericrenderer.PipelineComponent]
}

func NewSystem(
	assets assets.Assets,
	logger logger.Logger,
	gameAssets gameassets.GameAssets,
	transformToolFactory ecs.ToolFactory[transform.Transform],
	renderToolFactory ecs.ToolFactory[render.Render],
	uiToolFactory ecs.ToolFactory[ui.Tool],
	textToolFactory ecs.ToolFactory[text.Text],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(world ecs.World) error {
		s := system{
			assets,
			gameAssets,

			logger,
			world,
			uiToolFactory.Build(world),
			transformToolFactory.Build(world).Transform(),
			renderToolFactory.Build(world).Render(),
			textToolFactory.Build(world).Text(),

			ecs.GetComponentsArray[hierarchy.Component](world),
			ecs.GetComponentsArray[groups.GroupsComponent](world),
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
	// render
	// collider
	// click

	// changes
	labelEntity := s.world.NewEntity()
	s.hierarchyArray.SaveComponent(labelEntity, hierarchy.NewParent(parent))
	s.inheritGroupsArray.SaveComponent(labelEntity, groups.InheritGroupsComponent{})

	s.transformTool.Parent().SaveComponent(labelEntity, transform.NewParent(transform.RelativePos|transform.RelativeSizeX))
	s.transformTool.ParentPivotPoint().SaveComponent(labelEntity, transform.NewParentPivotPoint(.5, .5, 1))
	s.transformTool.PivotPoint().SaveComponent(labelEntity, transform.NewPivotPoint(.5, .5, 0))
	s.transformTool.Pos().SaveComponent(labelEntity, transform.NewPos(0, 25, 0))
	s.transformTool.Size().SaveComponent(labelEntity, transform.NewSize(1, 50, 1))

	s.textTool.TextContent().SaveComponent(labelEntity, text.TextComponent{Text: "SETTINGS"})
	s.textTool.FontSize().SaveComponent(labelEntity, text.FontSizeComponent{FontSize: 25})
	s.textTool.TextAlign().SaveComponent(labelEntity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	//

	type Button struct {
		text  string
		event any
	}
	btns := []Button{
		{"SHOT", audio.NewPlayEvent(gamescenes.EffectChannel, s.gameAssets.ExampleAudio)},
		{"QUIT", scenes.NewChangeSceneEvent(gamescenes.MenuID)},
	}

	btnAsset, err := assets.GetAsset[render.TextureAsset](s.assets, s.gameAssets.Hud.Btn)
	if err != nil {
		return err
	}
	btnAspectRatio := btnAsset.AspectRatio()

	for i, btn := range btns {
		var height float32 = 50
		var margin float32 = 10
		btnEntity := s.world.NewEntity()
		s.hierarchyArray.SaveComponent(btnEntity, hierarchy.NewParent(parent))
		s.inheritGroupsArray.SaveComponent(btnEntity, groups.InheritGroupsComponent{})

		// btnTransform := transformTransaction.GetObject(btnEntity)
		s.transformTool.AspectRatio().SaveComponent(btnEntity, transform.NewAspectRatio(float32(btnAspectRatio.Dx()), float32(btnAspectRatio.Dy()), 0, transform.PrimaryAxisX))
		s.transformTool.Parent().SaveComponent(btnEntity, transform.NewParent(transform.RelativePos|transform.RelativeSizeX))
		s.transformTool.ParentPivotPoint().SaveComponent(btnEntity, transform.NewParentPivotPoint(.5, .5, 1))
		s.transformTool.PivotPoint().SaveComponent(btnEntity, transform.NewPivotPoint(.5, .5, 0))
		s.transformTool.MaxSize().SaveComponent(btnEntity, transform.NewMaxSize(0, height+margin-1, 0))
		s.transformTool.Pos().SaveComponent(btnEntity, transform.NewPos(0, float32(-i)*(height+margin)-25, 0))
		s.transformTool.Size().SaveComponent(btnEntity, transform.NewSize(1, height, 1))

		s.renderTool.Mesh().SaveComponent(btnEntity, render.NewMesh(s.gameAssets.SquareMesh))
		s.renderTool.Texture().SaveComponent(btnEntity, render.NewTexture(s.gameAssets.Hud.Btn))
		s.pipelineArray.SaveComponent(btnEntity, genericrenderer.PipelineComponent{})

		s.textTool.TextContent().SaveComponent(btnEntity, text.TextComponent{Text: btn.text})
		s.textTool.FontSize().SaveComponent(btnEntity, text.FontSizeComponent{FontSize: 25})
		s.textTool.TextAlign().SaveComponent(btnEntity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})

		s.leftClickArray.SaveComponent(btnEntity, inputs.NewMouseLeftClick(btn.event))
		s.keepSelectedArray.SaveComponent(btnEntity, inputs.KeepSelectedComponent{})
		s.colliderArray.SaveComponent(btnEntity, collider.NewCollider(s.gameAssets.SquareCollider))
	}
	return nil
}
