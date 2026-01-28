package internal

import (
	gameassets "core/assets"
	"core/modules/settings"
	"core/modules/ui"
	gamescenes "core/scenes"
	"engine"
	"engine/modules/audio"
	"engine/modules/collider"
	"engine/modules/groups"
	"engine/modules/inputs"
	"engine/modules/render"
	"engine/modules/scene"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/assets"
	"engine/services/ecs"
	"engine/services/frames"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

// 1. settings text
// 2. quit button

type system struct {
	GameAssets gameassets.GameAssets `inject:"1"`

	engine.World `inject:"1"`
	Ui           ui.Service `inject:"1"`
}

type temporaryToggleColorComponent struct{}

func NewSystem(c ioc.Dic) settings.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[system](c)

		events.ListenE(s.EventsBuilder, func(event settings.EnterSettingsForParentEvent) error {
			return s.ListenRender(event.Parent)
		})
		events.Listen(s.EventsBuilder, s.ListenOnTick)
		events.Listen(s.EventsBuilder, func(settings.EnterSettingsEvent) {
			event := settings.EnterSettingsForParentEvent{
				Parent: s.Ui.Show(),
			}
			events.Emit(s.Events, event)
		})

		return nil
	})
}

func (s *system) ListenOnTick(frames.TickEvent) {
	toggleArray := ecs.GetComponentsArray[temporaryToggleColorComponent](s)
	for _, entity := range toggleArray.GetEntities() {
		color, ok := s.Render.Color().Get(entity)
		if !ok {
			color.Color = mgl32.Vec4{1, 1, 1, 1}
		}

		color.Color[1] = 1 - color.Color[1]
		color.Color[2] = 1 - color.Color[2]

		s.Render.Color().Set(entity, color)
	}

}

func (s *system) ListenRender(parent ecs.EntityID) error {
	// render
	// collider
	// click

	children := []ecs.EntityID{}

	// changes
	labelEntity := s.NewEntity()
	children = append(children, labelEntity)
	s.Groups.Inherit().Set(labelEntity, groups.InheritGroupsComponent{})

	s.Transform.Parent().Set(labelEntity, transform.NewParent(transform.RelativePos|transform.RelativeSizeX))
	s.Transform.Size().Set(labelEntity, transform.NewSize(1, 50, 1))

	s.Text.Content().Set(labelEntity, text.TextComponent{Text: "SETTINGS"})
	s.Text.FontSize().Set(labelEntity, text.FontSizeComponent{FontSize: 25})
	s.Text.Align().Set(labelEntity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	//

	type Button struct {
		text  string
		event any
	}
	btns := []Button{
		{"SHOT", audio.NewPlayEvent(gamescenes.EffectChannel, s.GameAssets.ExampleAudio)},
		{"SHOT2", audio.NewPlayEvent(gamescenes.EffectChannel, s.GameAssets.ExampleAudio)},
		{"SHOT3", audio.NewPlayEvent(gamescenes.EffectChannel, s.GameAssets.ExampleAudio)},
		{"QUIT", scene.NewChangeSceneEvent(gamescenes.MenuID)},
	}

	btnAsset, err := assets.GetAsset[render.TextureAsset](s.Assets, s.GameAssets.Hud.Btn)
	if err != nil {
		return err
	}
	btnAspectRatio := btnAsset.AspectRatio()

	for _, btn := range btns {
		btnEntity := s.NewEntity()
		children = append(children, btnEntity)
		s.Groups.Inherit().Set(btnEntity, groups.InheritGroupsComponent{})

		ecs.GetComponentsArray[temporaryToggleColorComponent](s).Set(btnEntity, temporaryToggleColorComponent{})

		s.Transform.AspectRatio().Set(btnEntity, transform.NewAspectRatio(float32(btnAspectRatio.Dx()), float32(btnAspectRatio.Dy()), 0, transform.PrimaryAxisX))
		s.Transform.Parent().Set(btnEntity, transform.NewParent(transform.RelativePos|transform.RelativeSizeX))
		s.Transform.MaxSize().Set(btnEntity, transform.NewMaxSize(0, 50, 0))
		s.Transform.Size().Set(btnEntity, transform.NewSize(1, 50, 1))

		s.Render.Mesh().Set(btnEntity, render.NewMesh(s.GameAssets.SquareMesh))
		s.Render.Texture().Set(btnEntity, render.NewTexture(s.GameAssets.Hud.Btn))
		s.Render.Render(btnEntity)

		s.Text.Content().Set(btnEntity, text.TextComponent{Text: btn.text})
		s.Text.FontSize().Set(btnEntity, text.FontSizeComponent{FontSize: 25})
		s.Text.Align().Set(btnEntity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})

		s.Inputs.LeftClick().Set(btnEntity, inputs.NewLeftClick(btn.event))
		s.Inputs.KeepSelected().Set(btnEntity, inputs.KeepSelectedComponent{})
		s.Collider.Component().Set(btnEntity, collider.NewCollider(s.GameAssets.SquareCollider))
	}

	s.Hierarchy.SetChildren(parent, children...)

	return nil
}
