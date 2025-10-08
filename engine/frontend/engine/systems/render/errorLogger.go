package rendersys

import (
	"fmt"
	"frontend/services/frames"
	"shared/services/logger"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type ErrorLogger struct {
	logger logger.Logger
}

func NewErrorLogger() ErrorLogger {
	return ErrorLogger{}
}

var glErrorStrings = map[uint32]string{
	gl.NO_ERROR:                      "GL_NO_ERROR",
	gl.INVALID_ENUM:                  "GL_INVALID_ENUM",
	gl.INVALID_VALUE:                 "GL_INVALID_VALUE",
	gl.INVALID_OPERATION:             "GL_INVALID_OPERATION",
	gl.STACK_OVERFLOW:                "GL_STACK_OVERFLOW",
	gl.STACK_UNDERFLOW:               "GL_STACK_UNDERFLOW",
	gl.OUT_OF_MEMORY:                 "GL_OUT_OF_MEMORY",
	gl.INVALID_FRAMEBUFFER_OPERATION: "GL_INVALID_FRAMEBUFFER_OPERATION",
	gl.CONTEXT_LOST:                  "GL_CONTEXT_LOST",
	// gl.TABLE_TOO_LARGE:               "GL_TABLE_TOO_LARGE", // Less common in modern GL
}

func (logger *ErrorLogger) Update(args frames.FrameEvent) {
	if glErr := gl.GetError(); glErr != gl.NO_ERROR {
		logger.logger.Error(fmt.Errorf("opengl error: %x %s\n", glErr, glErrorStrings[glErr]))
	}
}
