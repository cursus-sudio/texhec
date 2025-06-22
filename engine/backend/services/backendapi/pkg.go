package backendapi

import (
	"backend/services/scopes"
	"backend/utils/endpoint"
	"backend/utils/httperrors"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterTransient(b, func(c ioc.Dic) Builder {
		b := NewBuilder()
		r := b.Relay()
		userSession := c.Scope(scopes.UserSession)

		r.DefaultHandler(func(ctx relay.AnyContext) { ctx.SetErr(httperrors.Err404) })
		r.RegisterMiddleware(func(ctx relay.AnyContext, next func()) {
			req := ctx.Req().(endpoint.AnyRequest)
			request := userSession.Scope(scopes.Request)
			req.UseC(request)
			next()
			requestService := ioc.Get[scopes.RequestEnd](request)
			requestService.Clean(scopes.NewRequestEndArgs(ctx.Err()))
		})
		return b
	})
	ioc.RegisterTransient(b, func(c ioc.Dic) Backend {
		return ioc.Get[Builder](c).Build()
	})
}
