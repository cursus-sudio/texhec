package internal

import (
	"engine/modules/render"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
)

type errorLogger struct {
	logger logger.Logger
	render render.Service
}

func NewErrorLogger(
	logger logger.Logger,
	renderService render.Service,
	eventsBuilder events.Builder,
) render.System {
	return ecs.NewSystemRegister(func() error {
		s := &errorLogger{logger, renderService}
		events.Listen(eventsBuilder, s.Listen)
		return nil
	})
}

func (logger *errorLogger) Listen(args frames.FrameEvent) {
	if glErr := gl.GetError(); glErr != gl.NO_ERROR {
		logger.logger.Warn(fmt.Errorf("opengl error: %x %s", glErr, glErrorStrings[glErr]))
	}
}
