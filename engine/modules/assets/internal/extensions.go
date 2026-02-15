package internal

import (
	"engine/modules/assets"
	"engine/services/logger"
	"fmt"
	"strings"

	"github.com/ogiusek/ioc/v2"
)

type extensions struct {
	Logger     logger.Logger `inject:"1"`
	extensions map[string]func(assets.Path) (any, error)
}

func NewExtensions(c ioc.Dic) assets.Extensions {
	e := ioc.GetServices[*extensions](c)
	e.extensions = make(map[string]func(assets.Path) (any, error))
	return e
}

func (s *extensions) Register(
	/* shouldn't have dots and be after dots in asset */ extension string,
	dispatcher func(path assets.Path) (any, error),
) {
	extension = strings.Trim(extension, ".")
	if _, ok := s.extensions[extension]; ok {
		s.Logger.Warn(fmt.Errorf("extension \"%v\" is already taken", extension))
		return
	}
	s.extensions[extension] = dispatcher
}

func (s *extensions) PathExntesion(path assets.Path) string {
	parts := strings.Split(string(path), ".")
	return parts[len(parts)-1]
}

func (s *extensions) ExtensionDispatcher(extension string) (func(assets.Path) (any, error), bool) {
	extension = strings.Trim(extension, ".")
	dispatcher, ok := s.extensions[extension]
	return dispatcher, ok
}
