package api

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
	ioc.RegisterTransient(b, func(c ioc.Dic) ServerBuilder {
		b := NewServerBuilder()
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
	ioc.RegisterTransient(b, func(c ioc.Dic) Server {
		return ioc.Get[ServerBuilder](c).Build(func() {
			userSessionService := ioc.Get[scopes.UserSessionEnd](c)
			userSessionService.Clean(scopes.NewUserSessionEndArgs())
		})
	})
}
