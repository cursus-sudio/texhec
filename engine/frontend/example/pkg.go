package main

import (
	"frontend/example/ping"
	"frontend/example/tacticalmap"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	for _, pkg := range []ioc.Pkg{
		ping.Package(),
		tacticalmap.Package(),
	} {
		pkg.Register(b)
	}

}
