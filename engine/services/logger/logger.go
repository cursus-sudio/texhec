package logger

import (
	"engine/services/clock"
	"fmt"
)

// nil errors are not warned or fatal
type Logger interface {
	Info(format string, a ...any)
	Warn(err error)
	Fatal(err error)
}

type logger struct {
	PanicOnError bool
	Clock        clock.Clock
	Print        func(string)
	Panic        func(string)
}

func (logger *logger) Info(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	msg := fmt.Sprintf("\033[34m[ Info ]\033[0m %s \033[34m\n%s\033[0m\n", logger.Clock.Now(), message)
	logger.Print(msg)
}

func TestWarning(logger Logger) {
	// print("hihi %v")
	logger.Info("hihi %v")
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

// func renderColor(b *strings.Builder, color mgl64.Vec3) {
// 	fmt.Fprintf(b,
// 		"\033[38;2;%d;%d;%dm█",
// 		uint8(color[0]*255),
// 		uint8(color[1]*255),
// 		uint8(color[2]*255),
// 	)
// }
