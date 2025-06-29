package connpkg

import (
	"shared/services/connpkg/netconnection"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	netconnection.Package().Register(b)
}
