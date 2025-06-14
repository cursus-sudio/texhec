package scopecleanup

import "github.com/ogiusek/ioc"

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(c ioc.Dic) {
	// ioc.RegisterScoped(c, func(c ioc.Dic) ScopeCleanUp { return newScopeCleanUp() })
	ioc.RegisterScoped(c, func(c ioc.Dic) ScopeCleanUp { return newScopeCleanUp() })
}
