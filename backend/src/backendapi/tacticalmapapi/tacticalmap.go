package tacticalmapapi

import (
	"backend/src/backendapi"
	"backend/src/modules/tacticalmap"
	"backend/src/utils/endpoint"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type CreateReq struct {
	relay.Req[CreateRes]
	tacticalmap.CreateArgs
}

func NewCreateReq(args tacticalmap.CreateArgs) CreateReq { return CreateReq{CreateArgs: args} }

type CreateRes struct{}
type createEndpoint struct {
	TacticalMap tacticalmap.TacticalMap `inject:"1"`
}

func (endpoint createEndpoint) Handle(req CreateReq) (CreateRes, error) {
	err := endpoint.TacticalMap.Create(req.CreateArgs)
	return CreateRes{}, err
}

//

type DestroyReq struct {
	relay.Req[DestroyRes]
	tacticalmap.DestroyArgs
}

func NewDestroyReq(args tacticalmap.DestroyArgs) DestroyReq { return DestroyReq{DestroyArgs: args} }

type DestroyRes struct{}
type destroyEndpoint struct {
	TacticalMap tacticalmap.TacticalMap `inject:"1"`
}

func (endpoint destroyEndpoint) Handle(req DestroyReq) (DestroyRes, error) {
	err := endpoint.TacticalMap.Destroy(req.DestroyArgs)
	return DestroyRes{}, err
}

//

type GetReq struct {
	relay.Req[GetRes]
}

func NewGetReq() GetReq { return GetReq{} }

type GetRes struct {
	Tiles []tacticalmap.Tile
}
type getEndpoint struct {
	TacticalMap tacticalmap.TacticalMap `inject:"1"`
}

func (endpoint getEndpoint) Handle(req GetReq) (GetRes, error) {
	tiles, err := endpoint.TacticalMap.GetMap()
	return GetRes{Tiles: tiles}, err
}

//

type Pkg struct {
}

func Package() Pkg {
	return Pkg{}
}

//

// OnCreate(CreateListener)
// OnDestroy(DestroyListener)

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.WrapService[backendapi.Builder](b, func(c ioc.Dic, s backendapi.Builder) backendapi.Builder {
		r := s.Relay()
		endpoint.Register[createEndpoint](c, r)
		endpoint.Register[destroyEndpoint](c, r)
		endpoint.Register[getEndpoint](c, r)
		return s
	})
}
