package scopecleanup

import "github.com/ogiusek/ioc/v2"

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	// ioc.RegisterScoped(c, func(c ioc.Dic) ScopeCleanUp { return newScopeCleanUp() })
	ioc.RegisterScoped(b, func(c ioc.Dic) ScopeCleanUp { return newScopeCleanUp() })
}
