package systems

import (
	"engine/modules/render"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"
	"fmt"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type errorLogger struct {
	Logger        logger.Logger  `inject:"1"`
	Render        render.Service `inject:"1"`
	EventsBuilder events.Builder `inject:"1"`
}

func NewErrorLogger(c ioc.Dic) render.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*errorLogger](c)
		events.Listen(s.EventsBuilder, s.Listen)
		return nil
	})
}

func (logger *errorLogger) Listen(args frames.FrameEvent) {
	if err := logger.Render.Error(); err != nil {
		logger.Logger.Warn(fmt.Errorf("opengl error: %s", err))
	}
}
