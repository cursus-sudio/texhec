package internal

import (
	"fmt"
	"frontend/engine/render"
	"frontend/services/frames"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
)

type errorLogger struct {
	logger logger.Logger
	render.RenderTool
}

func NewErrorLogger(logger logger.Logger, t render.RenderTool) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &errorLogger{logger, t}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (logger *errorLogger) Listen(args frames.FrameEvent) {
	if glErr := gl.GetError(); glErr != gl.NO_ERROR {
		logger.logger.Error(fmt.Errorf("opengl error: %x %s\n", glErr, glErrorStrings[glErr]))
	}
}
