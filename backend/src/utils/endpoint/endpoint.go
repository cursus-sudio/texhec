package endpoint

import (
	"backend/src/utils/services/scopecleanup"

	"github.com/ogiusek/ioc"
	"github.com/ogiusek/relay"
)

type endpoint[Req relay.Req[Res], Res relay.Res] interface {
	Handle(Req) (Res, error)
}

func Register[Endpoint endpoint[Req, Res], Req relay.Req[Res], Res relay.Res](c ioc.Dic, r relay.Relay) {
	relay.Register(r, func(req Req) (Res, error) {
		c := c.Scope()
		endpoint := ioc.GetServices[Endpoint](c)
		res, err := endpoint.Handle(req)
		cleanUpArgs := scopecleanup.NewCleanUpArgs(err)
		ioc.Get[scopecleanup.ScopeCleanUp](c).Clean(cleanUpArgs)
		return res, err
	})
}
