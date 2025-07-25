package scopes

import "github.com/ogiusek/ioc/v2"

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {

	b.RegisterScope(Request)
	ioc.RegisterScoped(b, Request, func(c ioc.Dic) RequestService {
		return newRequestService()
	})
}
