package logger

import (
	"fmt"
	"shared/services/clock"
)

// nil errors are not warned or fatal
type Logger interface {
	Info(message string)
	Warn(err error)
	Fatal(err error)
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

func (logger *logger) Warn(err error) {
	if err == nil {
		return
	}
	message := fmt.Sprintf("\033[31m[ Error ]\033[0m %s \033[31m\n%s\033[0m\n", logger.Clock.Now(), err)
	if logger.PanicOnError {
		logger.Panic(message)
	} else {
		logger.Print(message)
	}
}

func (logger *logger) Fatal(err error) {
	if err == nil {
		return
	}
	message := fmt.Sprintf("\033[31m[ Error ]\033[0m %s \033[31m\n%s\033[0m\n", logger.Clock.Now(), err)
	logger.Panic(message)
}
