package api

import (
	"reflect"
	"shared/utils/connection"
	"shared/utils/httperrors"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type Pkg struct {
	pkgs            []ioc.Pkg
	getRequestScope func(c ioc.Dic) ioc.Dic
}

func Package(
	getRequestScope func(c ioc.Dic) ioc.Dic,
) Pkg {
	return Pkg{
		getRequestScope: getRequestScope,
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, connection.OrderAttachServices, func(c ioc.Dic, con connection.Definition) connection.Definition {
		con.MessageListener().Relay().RegisterMiddleware(func(ctx relay.AnyContext, next func()) {
			c := reflect.ValueOf(pkg.getRequestScope(c))

			reqValue := reflect.ValueOf(ctx.Req())
			newReqValue := reflect.New(reqValue.Type()).Elem()
			newReqValue.Set(reqValue)
			newReqValue.FieldByName("C").Set(c)

			ctx.SetReq(newReqValue.Interface())

			next()
		})
		con.MessageListener().Relay().RegisterMessageMiddleware(func(ctx relay.AnyMessageCtx, next func()) {
			c := reflect.ValueOf(pkg.getRequestScope(c))

			msgValue := reflect.ValueOf(ctx.Message())
			newMsgValue := reflect.New(msgValue.Type()).Elem()
			newMsgValue.Set(msgValue)
			newMsgValue.FieldByName("C").Set(c)

			ctx.SetMessage(newMsgValue.Interface())

			go next()
		})
		return con
	})

	ioc.WrapService(b, connection.OrderEndpoint, func(c ioc.Dic, con connection.Definition) connection.Definition {
		con.MessageListener().Relay().DefaultHandler(func(ctx relay.AnyContext) {
			ctx.SetErr(httperrors.Err404)
		})
		return con
	})
}
