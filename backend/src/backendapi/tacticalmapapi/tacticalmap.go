package tacticalmapapi

import (
	"backend/src/modules/tacticalmap"
	"backend/src/utils/endpoint"

	"github.com/ogiusek/ioc"
	"github.com/ogiusek/relay"
)

// OnCreate(CreateListener)
// OnDestroy(DestroyListener)

type Uh struct{}

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
	c ioc.Dic
}

func Package(c ioc.Dic) Pkg {
	return Pkg{c: c}
}

func (pkg Pkg) Register(r relay.Relay) {
	endpoint.Register[createEndpoint](pkg.c, r)
	endpoint.Register[destroyEndpoint](pkg.c, r)
	endpoint.Register[getEndpoint](pkg.c, r)
}
