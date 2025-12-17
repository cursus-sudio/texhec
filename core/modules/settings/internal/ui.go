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

	uiTool        ui.Interface
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
	transformToolFactory ecs.ToolFactory[transform.TransformTool],
	renderToolFactory ecs.ToolFactory[render.RenderTool],
	uiToolFactory ecs.ToolFactory[ui.UiTool],
	textToolFactory ecs.ToolFactory[text.TextTool],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(world ecs.World) error {
		s := system{
			assets,
			gameAssets,

			logger,
			world,
			uiToolFactory.Build(world).Ui(),
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
	s.hierarchyArray.Set(labelEntity, hierarchy.NewParent(parent))
	s.inheritGroupsArray.Set(labelEntity, groups.InheritGroupsComponent{})

	s.transformTool.Parent().Set(labelEntity, transform.NewParent(transform.RelativePos|transform.RelativeSizeX))
	s.transformTool.ParentPivotPoint().Set(labelEntity, transform.NewParentPivotPoint(.5, .5, 1))
	s.transformTool.PivotPoint().Set(labelEntity, transform.NewPivotPoint(.5, .5, 0))
	s.transformTool.Pos().Set(labelEntity, transform.NewPos(0, 25, 0))
	s.transformTool.Size().Set(labelEntity, transform.NewSize(1, 50, 1))

	s.textTool.Content().Set(labelEntity, text.TextComponent{Text: "SETTINGS"})
	s.textTool.FontSize().Set(labelEntity, text.FontSizeComponent{FontSize: 25})
	s.textTool.Align().Set(labelEntity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})

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
		s.hierarchyArray.Set(btnEntity, hierarchy.NewParent(parent))
		s.inheritGroupsArray.Set(btnEntity, groups.InheritGroupsComponent{})

		// btnTransform := transformTransaction.GetObject(btnEntity)
		s.transformTool.AspectRatio().Set(btnEntity, transform.NewAspectRatio(float32(btnAspectRatio.Dx()), float32(btnAspectRatio.Dy()), 0, transform.PrimaryAxisX))
		s.transformTool.Parent().Set(btnEntity, transform.NewParent(transform.RelativePos|transform.RelativeSizeX))
		s.transformTool.ParentPivotPoint().Set(btnEntity, transform.NewParentPivotPoint(.5, .5, 1))
		s.transformTool.PivotPoint().Set(btnEntity, transform.NewPivotPoint(.5, .5, 0))
		s.transformTool.MaxSize().Set(btnEntity, transform.NewMaxSize(0, height+margin-1, 0))
		s.transformTool.Pos().Set(btnEntity, transform.NewPos(0, float32(-i)*(height+margin)-25, 0))
		s.transformTool.Size().Set(btnEntity, transform.NewSize(1, height, 1))

		s.renderTool.Mesh().Set(btnEntity, render.NewMesh(s.gameAssets.SquareMesh))
		s.renderTool.Texture().Set(btnEntity, render.NewTexture(s.gameAssets.Hud.Btn))
		s.pipelineArray.Set(btnEntity, genericrenderer.PipelineComponent{})

		s.textTool.Content().Set(btnEntity, text.TextComponent{Text: btn.text})
		s.textTool.FontSize().Set(btnEntity, text.FontSizeComponent{FontSize: 25})
		s.textTool.Align().Set(btnEntity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})

		s.leftClickArray.Set(btnEntity, inputs.NewMouseLeftClick(btn.event))
		s.keepSelectedArray.Set(btnEntity, inputs.KeepSelectedComponent{})
		s.colliderArray.Set(btnEntity, collider.NewCollider(s.gameAssets.SquareCollider))
	}
	return nil
}
