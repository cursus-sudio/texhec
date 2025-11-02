package logger

import (
	"fmt"
	"shared/services/clock"

	"github.com/ogiusek/ioc/v2"
)

type Logger interface {
	Info(message string)
	Error(err ...error)
	Fatal(err ...error)
}

type logger struct {
	PanicOnError bool
	Clock        clock.Clock
	Print        func(string)
	Panic        func(string)
}

func (logger *logger) Info(message string) {
	msg := fmt.Sprintf("\033[34m[ Info ]\033[0m %s \033[34m\n%s\033[0m\n", logger.Clock.Now(), message)
	logger.Print(msg)
}

func (logger *logger) Error(err ...error) {
	message := fmt.Sprintf("\033[31m[ Error ]\033[0m %s \033[31m\n%s\033[0m\n", logger.Clock.Now(), err)
	if logger.PanicOnError {
		logger.Panic(message)
	} else {
		logger.Print(message)
	}
}

func (logger *logger) Fatal(err ...error) {
	message := fmt.Sprintf("\033[31m[ Error ]\033[0m %s \033[31m\n%s\033[0m\n", logger.Clock.Now(), err)
	logger.Panic(message)
}

type pkg struct {
	panicOnError bool
	print        func(c ioc.Dic, message string)
}

func Package(
	panicOnError bool,
	print func(c ioc.Dic, message string),
) ioc.Pkg {
	return pkg{
		panicOnError: panicOnError,
		print:        print,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Logger {
		return &logger{
			PanicOnError: pkg.panicOnError,
			Clock:        ioc.Get[clock.Clock](c),
			Print:        func(s string) { pkg.print(c, s) },
			Panic:        func(s string) { panic(s) },
		}
	})
	ioc.RegisterDependency[Logger, clock.Clock](b)
}
