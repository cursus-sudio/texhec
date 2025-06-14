package backendapi

import (
	"github.com/ogiusek/ioc"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(c ioc.Dic) {
	ioc.RegisterTransient(c, func(c ioc.Dic) Backend { return ioc.GetServices[backend](c) })
}
