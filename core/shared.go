package main

import (
	"shared/services/api/netconnection"
	"shared/services/clock"
	"shared/services/codec"
	"shared/services/events"
	"shared/services/runtime"
	"shared/services/uuid"

	"github.com/ogiusek/ioc/v2"
)

type SharedPkg struct {
	pkgs []ioc.Pkg
}

func SharedPackage(
	netconnectionPkg netconnection.Pkg,
	clockPkg clock.Pkg,
) SharedPkg {
	return SharedPkg{
		pkgs: []ioc.Pkg{
			netconnectionPkg,
			clockPkg,
			events.Package(),
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
