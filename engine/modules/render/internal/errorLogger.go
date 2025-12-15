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
	render.Render
}

func NewErrorLogger(logger logger.Logger, t render.Render) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &errorLogger{logger, t}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (logger *errorLogger) Listen(args frames.FrameEvent) {
	if glErr := gl.GetError(); glErr != gl.NO_ERROR {
		logger.logger.Warn(fmt.Errorf("opengl error: %x %s\n", glErr, glErrorStrings[glErr]))
	}
}
