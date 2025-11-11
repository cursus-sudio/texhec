package logger

import (
	"fmt"
	"shared/services/clock"

	"github.com/ogiusek/ioc/v2"
)

// doesn't warn or is fatal when there are no not-nil errors
type Logger interface {
	Info(message string)
	Warn(err ...error)
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

func (logger *logger) FilterNils(allErrors ...error) []error {
	notNilErrors := []error{}
	for _, err := range allErrors {
		if err == nil {
			continue
		}
		notNilErrors = append(notNilErrors, err)
	}
	return notNilErrors
}

func (logger *logger) Warn(errors ...error) {
	if errors = logger.FilterNils(errors...); len(errors) == 0 {
		return
	}
	message := fmt.Sprintf("\033[31m[ Error ]\033[0m %s \033[31m\n%s\033[0m\n", logger.Clock.Now(), errors)
	if logger.PanicOnError {
		logger.Panic(message)
	} else {
		logger.Print(message)
	}
}

func (logger *logger) Fatal(errors ...error) {
	if errors = logger.FilterNils(errors...); len(errors) == 0 {
		return
	}
	message := fmt.Sprintf("\033[31m[ Error ]\033[0m %s \033[31m\n%s\033[0m\n", logger.Clock.Now(), errors)
	logger.Panic(message)
}

type pkg struct {
	panicOnWarn bool
	print       func(c ioc.Dic, message string)
}

func Package(
	panicOnWarn bool,
	print func(c ioc.Dic, message string),
) ioc.Pkg {
	return pkg{
		panicOnWarn: panicOnWarn,
		print:       print,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Logger {
		return &logger{
			PanicOnError: pkg.panicOnWarn,
			Clock:        ioc.Get[clock.Clock](c),
			Print:        func(s string) { pkg.print(c, s) },
			Panic:        func(s string) { panic(s) },
		}
	})
	ioc.RegisterDependency[Logger, clock.Clock](b)
}
