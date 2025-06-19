package tacticalmap

import (
	"github.com/ogiusek/relay/v2"
)

type CreateReq struct {
	relay.Req[CreateRes]
	CreateArgs
}

func NewCreateReq(args CreateArgs) CreateReq { return CreateReq{CreateArgs: args} }

type CreateRes struct{}
type createEndpoint struct {
	TacticalMap TacticalMap `inject:"1"`
}

func (endpoint createEndpoint) Handle(req CreateReq) (CreateRes, error) {
	err := endpoint.TacticalMap.Create(req.CreateArgs)
	return CreateRes{}, err
}

//

type DestroyReq struct {
	relay.Req[DestroyRes]
	DestroyArgs
}

func NewDestroyReq(args DestroyArgs) DestroyReq { return DestroyReq{DestroyArgs: args} }

type DestroyRes struct{}
type destroyEndpoint struct {
	TacticalMap TacticalMap `inject:"1"`
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
	Tiles []Tile
}
type getEndpoint struct {
	TacticalMap TacticalMap `inject:"1"`
}

func (endpoint getEndpoint) Handle(req GetReq) (GetRes, error) {
	tiles, err := endpoint.TacticalMap.GetMap()
	return GetRes{Tiles: tiles}, err
}

// TODO

// OnCreate(CreateListener)
// OnDestroy(DestroyListener)
