package logger

import (
	"backend/src/utils/clock"
	"fmt"

	"github.com/ogiusek/ioc/v2"
)

type Logger interface {
	Info(message string)
	Error(err error)
}

type logger struct {
	Clock clock.Clock
}

func (logger *logger) Info(message string) {
	fmt.Printf("\033[34m[ Info ]\033[0m %s \033[34m\n%s\033[0m\n\n", logger.Clock.Now(), message)
}

func (logger *logger) Error(err error) {
	fmt.Printf("\033[31m[ Error ]\033[0m %s \033[31m\n%s\033[0m\n\n", logger.Clock.Now(), err)
}

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Logger { return &logger{Clock: ioc.Get[clock.Clock](c)} })
}
