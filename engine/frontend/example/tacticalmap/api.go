package tacticalmap

import (
	"backend/utils/endpoint"
)

type CreateReq struct {
	endpoint.Request[CreateRes]
	CreateArgs
}

func NewCreateReq(args CreateArgs) CreateReq {
	return CreateReq{Request: endpoint.NewRequest[CreateRes](), CreateArgs: args}
}

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
	endpoint.Request[DestroyRes]
	DestroyArgs
}

func NewDestroyReq(args DestroyArgs) DestroyReq {
	return DestroyReq{Request: endpoint.NewRequest[DestroyRes](), DestroyArgs: args}
}

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
	endpoint.Request[GetRes]
}

func NewGetReq() GetReq { return GetReq{Request: endpoint.NewRequest[GetRes]()} }

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
