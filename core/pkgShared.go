package main

import (
	"shared/services/api/netconnection"
	"shared/services/clock"
	"shared/services/codec"
	"shared/services/eventspkg"
	"shared/services/runtime"
	"shared/services/uuid"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type SharedPkg struct {
	pkgs []ioc.Pkg
}

func SharedPackage() SharedPkg {
	return SharedPkg{
		pkgs: []ioc.Pkg{
			netconnection.Package(time.Second),
			clock.Package(time.RFC3339Nano),
			// netconnectionPkg,
			// clockPkg,
			eventspkg.Package(),
			codec.Package(),
			runtime.Package(),
			uuid.Package(),
		},
	}
}

func (pkg SharedPkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
