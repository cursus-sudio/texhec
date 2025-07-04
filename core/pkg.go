package main

import (
	"core/ping"
	"core/tacticalmap"

	"github.com/ogiusek/ioc/v2"
)

type ServerPkg struct{}

func ServerPackage() ServerPkg {
	return ServerPkg{}
}

func (ServerPkg) Register(b ioc.Builder) {
	for _, pkg := range []ioc.Pkg{
		exBackendModPkg{},
		ping.ServerPackage(),
		tacticalmap.ServerPackage(),
	} {
		pkg.Register(b)
	}
}

type ClientPkg struct{}

func ClientPackage() ClientPkg {
	return ClientPkg{}
}

func (ClientPkg) Register(b ioc.Builder) {
	for _, pkg := range []ioc.Pkg{
		ping.ClientPackage(),
		tacticalmap.ClientPackage(),
	} {
		pkg.Register(b)
	}
}
