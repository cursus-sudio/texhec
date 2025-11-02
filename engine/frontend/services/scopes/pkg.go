package scopes

import "github.com/ogiusek/ioc/v2"

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	b.RegisterScope(Scene)

	b.RegisterScope(Request)
	ioc.RegisterScoped(b, Request, func(c ioc.Dic) RequestService {
		return newRequestService()
	})
}
