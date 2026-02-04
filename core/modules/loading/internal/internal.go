package internal

import (
	gameassets "core/assets"
	"core/modules/loading"
	"core/modules/ui"
	"engine"
	"engine/modules/camera"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/frames"
	"fmt"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type system struct {
	World      engine.World          `inject:"1"`
	GameAssets gameassets.GameAssets `inject:"1"`
	Ui         ui.Service            `inject:"1"`

	Camera *ecs.EntityID
	Text   ecs.EntityID
}

func NewSystem(c ioc.Dic) loading.System {
	s := ioc.GetServices[*system](c)
	return ecs.NewSystemRegister(func() error {
		events.Listen(s.World.EventsBuilder, s.Listen)
		return nil
	})
}

func (s *system) Hide() {
	if s.Camera == nil {
		return
	}
	s.World.RemoveEntity(*s.Camera)
}

func (s *system) Render(message string) {
	if s.Camera != nil {
		s.World.Text.Content().Set(s.Text, text.TextComponent{Text: message})
		return
	}

	cameraEntity := s.World.NewEntity()
	s.World.Camera.Ortho().Set(cameraEntity, camera.NewOrtho(-5, 5))

	background := s.World.NewEntity()
	s.World.Hierarchy.SetParent(background, cameraEntity)
	s.World.Transform.Pos().Set(background, transform.NewPos(0, 0, 1))
	s.World.Transform.PivotPoint().Set(background, transform.NewPivotPoint(.5, .5, 0))
	s.World.Transform.ParentPivotPoint().Set(background, transform.NewParentPivotPoint(.5, .5, 0))
	s.Ui.AnimatedBackground().Set(background, ui.NewAnimatedBackground())

	textEntity := s.World.NewEntity()
	s.World.Hierarchy.SetParent(textEntity, cameraEntity)
	s.World.Transform.Pos().Set(textEntity, transform.NewPos(0, 0, 2))
	s.World.Transform.Parent().Set(textEntity, transform.NewParent(transform.RelativePos))

	s.World.Text.Content().Set(textEntity, text.TextComponent{Text: message})
	s.World.Text.FontSize().Set(textEntity, text.FontSizeComponent{FontSize: 32})
	s.World.Text.Break().Set(textEntity, text.BreakComponent{Break: text.BreakNone})
	s.World.Text.Align().Set(textEntity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	s.Camera = &cameraEntity
	s.Text = textEntity
}

func (s *system) Listen(frames.FrameEvent) {
	progress := s.World.Batcher.Progress()
	if progress == -1 {
		s.Hide()
		return
	}

	if progress < 0 && progress != -1 {
		panic(progress)
	}
	message := fmt.Sprintf("Loading... %6.2f%%", progress*100)
	s.Render(message)
}
