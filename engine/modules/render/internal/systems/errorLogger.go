package systems

import (
	"engine/modules/render"
	"engine/modules/render/internal/service"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"
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
	if glErr := gl.GetError(); glErr != gl.NO_ERROR {
		logger.Logger.Warn(fmt.Errorf("opengl error: %x %s", glErr, service.GlErrorStrings[glErr]))
	}
}
