package logger

import (
	"fmt"
	"shared/services/clock"

	"github.com/ogiusek/ioc/v2"
)

type Logger interface {
	Info(message string)
	Error(err error)
}

type logger struct {
	PanicOnError bool
	Clock        clock.Clock
}

func (logger *logger) Info(message string) {
	fmt.Printf("\033[34m[ Info ]\033[0m %s \033[34m\n%s\033[0m\n", logger.Clock.Now(), message)
}

func (logger *logger) Error(err error) {
	message := fmt.Sprintf("\033[31m[ Error ]\033[0m %s \033[31m\n%s\033[0m\n", logger.Clock.Now(), err)
	if logger.PanicOnError {
		panic(message)
	} else {
		print(message)
	}
}

type Pkg struct {
	panicOnError bool
}

func Package(panicOnError bool) Pkg {
	return Pkg{
		panicOnError: panicOnError,
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Logger {
		return &logger{PanicOnError: pkg.panicOnError, Clock: ioc.Get[clock.Clock](c)}
	})
}
