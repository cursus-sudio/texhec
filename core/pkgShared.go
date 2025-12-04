package main

import (
	"engine/services/clock"
	"engine/services/codec"
	"engine/services/eventspkg"
	"engine/services/runtime"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type SharedPkg struct {
	pkgs []ioc.Pkg
}

func SharedPackage() SharedPkg {
	return SharedPkg{
		pkgs: []ioc.Pkg{
			clock.Package(time.RFC3339Nano),
			eventspkg.Package(),
			codec.Package(),
			runtime.Package(),
		},
	}
}

func (pkg SharedPkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
