package tacticalmap

import (
	"fmt"
	"shared/services/logger"
	"shared/utils/endpoint"

	"github.com/ogiusek/ioc/v2"
)

type CreatedMessage struct {
	endpoint.Message
	Added []Tile
}

type createdMessageEndpoint struct {
	Logger logger.Logger `inject:"1"`
}

func (endpoint createdMessageEndpoint) Handle(req CreatedMessage) {
	text := fmt.Sprintf("created tiles %v\n", req)
	endpoint.Logger.Info(text)
}

type DestroyedMessage struct {
	endpoint.Message
	Destroyed []Tile
}

type destroyedMessageEndpoint struct {
	Logger logger.Logger `inject:"1"`
}

func (endpoint destroyedMessageEndpoint) Handle(req DestroyedMessage) {
	text := fmt.Sprintf("\ndestroyed tiles %v\n", req)
	endpoint.Logger.Info(text)
}

//

type FrontendPkg struct{}

func FrontendPackage() FrontendPkg {
	return FrontendPkg{}
}

func (FrontendPkg) Register(b ioc.Builder) {
	Package().Register(b)
	endpoint.MessageRegister[createdMessageEndpoint](b)
	endpoint.MessageRegister[destroyedMessageEndpoint](b)
}
