package api

import (
	"fmt"
	"frontend/services/console"
	"reflect"
	"shared/utils/connection"
	"shared/utils/endpoint"
	"shared/utils/httperrors"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) connection.Definition {
		return connection.NewDefinition()
	})

	ioc.WrapService(b, func(c ioc.Dic, con connection.Definition) connection.Definition {
		console := ioc.Get[console.Console](c)
		con.MessageListener().Relay().RegisterMiddleware(func(ctx relay.AnyContext, next func()) {
			rawReq := ctx.Req()
			req, ok := rawReq.(endpoint.AnyRequest)
			if ok {
				req.UseC(c)
			} else {
				console.LogPermanentlyToConsole(
					fmt.Sprintf(
						"request of type `%s` doesn't implement `endpoint.AnyRequest` interface",
						reflect.TypeOf(rawReq),
					),
				)
			}
			next()
		})
		con.MessageListener().Relay().RegisterMessageMiddleware(func(ctx relay.AnyMessageCtx, next func()) {
			rawReq := ctx.Message()
			req, ok := rawReq.(endpoint.Message)
			if ok {
				req.UseC(c)
			} else {
				console.LogPermanentlyToConsole(
					fmt.Sprintf(
						"request of type `%s` doesn't implement `endpoint.AnyRequest` interface",
						reflect.TypeOf(rawReq),
					),
				)
			}
			go next()
		})
		return con
	})

	ioc.WrapService(b, func(c ioc.Dic, con connection.Definition) connection.Definition {
		con.MessageListener().Relay().DefaultHandler(func(ctx relay.AnyContext) {
			ctx.SetErr(httperrors.Err404)
		})
		return con
	})

	ioc.RegisterTransient(b, func(c ioc.Dic) connection.Connection {
		return ioc.Get[connection.Definition](c).Build()
	})
}
