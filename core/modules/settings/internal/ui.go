package internal

import (
	gameassets "core/assets"
	"core/modules/settings"
	gamescenes "core/scenes"
	"engine/modules/audio"
	"engine/modules/collider"
	"engine/modules/genericrenderer"
	"engine/modules/groups"
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
	settings.World
}

func NewSystem(
	assets assets.Assets,
	logger logger.Logger,
	gameAssets gameassets.GameAssets,
) settings.System {
	return ecs.NewSystemRegister(func(world settings.World) error {
		s := system{
			assets,
			gameAssets,

			logger,
			world,
		}

		events.ListenE(world.EventsBuilder(), func(event settings.EnterSettingsForParentEvent) error {
			return s.Render(event.Parent)
		})
		events.Listen(world.EventsBuilder(), func(settings.EnterSettingsEvent) {
			event := settings.EnterSettingsForParentEvent{
				Parent: s.Ui().Show(),
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
	labelEntity := s.NewEntity()
	s.Hierarchy().SetParent(labelEntity, parent)
	s.Groups().Inherit().Set(labelEntity, groups.InheritGroupsComponent{})

	s.Transform().Parent().Set(labelEntity, transform.NewParent(transform.RelativePos|transform.RelativeSizeX))
	s.Transform().ParentPivotPoint().Set(labelEntity, transform.NewParentPivotPoint(.5, .5, 1))
	s.Transform().PivotPoint().Set(labelEntity, transform.NewPivotPoint(.5, .5, 0))
	s.Transform().Pos().Set(labelEntity, transform.NewPos(0, 25, 0))
	s.Transform().Size().Set(labelEntity, transform.NewSize(1, 50, 1))

	s.Text().Content().Set(labelEntity, text.TextComponent{Text: "SETTINGS"})
	s.Text().FontSize().Set(labelEntity, text.FontSizeComponent{FontSize: 25})
	s.Text().Align().Set(labelEntity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})

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
		btnEntity := s.NewEntity()
		s.Hierarchy().SetParent(btnEntity, parent)
		s.Groups().Inherit().Set(btnEntity, groups.InheritGroupsComponent{})

		// btnTransform := transformTransaction.GetObject(btnEntity)
		s.Transform().AspectRatio().Set(btnEntity, transform.NewAspectRatio(float32(btnAspectRatio.Dx()), float32(btnAspectRatio.Dy()), 0, transform.PrimaryAxisX))
		s.Transform().Parent().Set(btnEntity, transform.NewParent(transform.RelativePos|transform.RelativeSizeX))
		s.Transform().ParentPivotPoint().Set(btnEntity, transform.NewParentPivotPoint(.5, .5, 1))
		s.Transform().PivotPoint().Set(btnEntity, transform.NewPivotPoint(.5, .5, 0))
		s.Transform().MaxSize().Set(btnEntity, transform.NewMaxSize(0, height+margin-1, 0))
		s.Transform().Pos().Set(btnEntity, transform.NewPos(0, float32(-i)*(height+margin)-25, 0))
		s.Transform().Size().Set(btnEntity, transform.NewSize(1, height, 1))

		s.World.Render().Mesh().Set(btnEntity, render.NewMesh(s.gameAssets.SquareMesh))
		s.World.Render().Texture().Set(btnEntity, render.NewTexture(s.gameAssets.Hud.Btn))
		s.GenericRenderer().Pipeline().Set(btnEntity, genericrenderer.PipelineComponent{})

		s.Text().Content().Set(btnEntity, text.TextComponent{Text: btn.text})
		s.Text().FontSize().Set(btnEntity, text.FontSizeComponent{FontSize: 25})
		s.Text().Align().Set(btnEntity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})

		s.Inputs().MouseLeft().Set(btnEntity, inputs.NewMouseLeftClick(btn.event))
		s.Inputs().KeepSelected().Set(btnEntity, inputs.KeepSelectedComponent{})
		s.Collider().Component().Set(btnEntity, collider.NewCollider(s.gameAssets.SquareCollider))
	}
	return nil
}
