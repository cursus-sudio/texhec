package internal

import (
	"core/modules/loading"
	"core/modules/registry"
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

type CamComp struct{}
type TextComp struct{}

type system struct {
	World      engine.World    `inject:"1"`
	GameAssets registry.Assets `inject:"1"`
	Ui         ui.Service      `inject:"1"`

	CamArr  ecs.ComponentsArray[CamComp]
	TextArr ecs.ComponentsArray[TextComp]
}

func NewSystem(c ioc.Dic) loading.System {
	s := ioc.GetServices[*system](c)
	s.CamArr = ecs.GetComponentsArray[CamComp](s.World)
	s.TextArr = ecs.GetComponentsArray[TextComp](s.World)
	return ecs.NewSystemRegister(func() error {
		events.Listen(s.World.EventsBuilder, s.Listen)
		return nil
	})
}

func (s *system) Hide() {
	for _, e := range s.CamArr.GetEntities() {
		s.World.RemoveEntity(e)
	}
	for _, e := range s.TextArr.GetEntities() {
		s.World.RemoveEntity(e)
	}
}

func (s *system) Render(message string) {
	if len(s.TextArr.GetEntities()) == 1 {
		textEntity := s.TextArr.GetEntities()[0]
		s.World.Text.Content().Set(textEntity, text.TextComponent{Text: message})
		return
	}

	cameraEntity := s.World.NewEntity()
	s.World.Camera.Ortho().Set(cameraEntity, camera.NewOrtho(-5, 5))
	s.CamArr.Set(cameraEntity, CamComp{})

	background := s.World.NewEntity()
	s.World.Hierarchy.SetParent(background, cameraEntity)
	s.World.Transform.Pos().Set(background, transform.NewPos(0, 0, 1))
	s.World.Transform.PivotPoint().Set(background, transform.NewPivotPoint(.5, .5, 0))
	s.World.Transform.ParentPivotPoint().Set(background, transform.NewParentPivotPoint(.5, .5, 0))
	s.Ui.AnimatedBackground().Set(background, ui.AnimatedBackgroundComponent{})

	textEntity := s.World.NewEntity()
	s.TextArr.Set(textEntity, TextComp{})
	s.World.Hierarchy.SetParent(textEntity, cameraEntity)
	s.World.Transform.Pos().Set(textEntity, transform.NewPos(0, 0, 2))
	s.World.Transform.Parent().Set(textEntity, transform.NewParent(transform.RelativePos))

	s.World.Text.Content().Set(textEntity, text.TextComponent{Text: message})
	s.World.Text.FontSize().Set(textEntity, text.FontSizeComponent{FontSize: 32})
	s.World.Text.Break().Set(textEntity, text.BreakComponent{Break: text.BreakNone})
	s.World.Text.Align().Set(textEntity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
}

func (s *system) Listen(frames.FrameEvent) {
	progress := s.World.Batcher.Progress()
	if progress == -1 {
		s.Hide()
		return
	}

	message := fmt.Sprintf("Loading... %6.2f%%", progress*100)
	s.Render(message)
}
