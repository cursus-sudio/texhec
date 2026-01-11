package logger

import (
	"engine/services/clock"

	"github.com/ogiusek/ioc/v2"
)

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
}
