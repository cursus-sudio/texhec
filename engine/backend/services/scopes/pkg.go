package scopes

import "github.com/ogiusek/ioc/v2"

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg pkg) Register(b ioc.Builder) {
	b.RegisterScope(Request)
	b.RegisterScope(UserSession)
	// ioc.RegisterScoped(c, func(c ioc.Dic) ScopeCleanUp { return newScopeCleanUp() })
	// ioc.RegisterScoped(b, func(c ioc.Dic) ScopeCleanUp { return newScopeCleanUp() })
	ioc.RegisterScoped(b, Request, func(c ioc.Dic) RequestService { return newRequestEnd() })
	ioc.RegisterScoped(b, UserSession, func(c ioc.Dic) UserSessionService { return newSessionEnd() })
}
