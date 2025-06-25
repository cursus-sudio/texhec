package tacticalmap

import (
	"backend/services/clients"
	"backend/services/logger"
	"shared/utils/endpoint"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
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
	TacticalMap TacticalMap           `inject:"1"`
	Logger      logger.Logger         `inject:"1"`
	Client      clients.SessionClient `inject:"1"`
}

func (e createEndpoint) Handle(req CreateReq) (CreateRes, error) {
	err := e.TacticalMap.Create(req.CreateArgs)
	e.Logger.Info("fml\n\n\n\n\n\n")
	client, _ := e.Client.Client()
	relay.HandleMessage(client.Connection.Relay(), CreatedMessage{
		Message: endpoint.NewMessage(),
		Added:   req.CreateArgs.Tiles,
	})
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
	TacticalMap TacticalMap           `inject:"1"`
	Client      clients.SessionClient `inject:"1"`
}

func (e destroyEndpoint) Handle(req DestroyReq) (DestroyRes, error) {
	err := e.TacticalMap.Destroy(req.DestroyArgs)
	client, _ := e.Client.Client()
	relay.HandleMessage(client.Connection.Relay(), DestroyedMessage{
		Message:   endpoint.NewMessage(),
		Destroyed: req.Tiles,
	})
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

//

type ServerPkg struct{}

func ServerPackage() ServerPkg {
	return ServerPkg{}
}

func (pkg ServerPkg) Register(b ioc.Builder) {
	Package().Register(b)
	ioc.RegisterSingleton(b, func(c ioc.Dic) TacticalMap { return newTacticalMap() })

	endpoint.Register[createEndpoint](b)
	endpoint.Register[destroyEndpoint](b)
	endpoint.Register[getEndpoint](b)
}
