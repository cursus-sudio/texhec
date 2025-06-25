package tacticalmap

import (
	"fmt"
	"frontend/services/console"
	"shared/utils/endpoint"

	"github.com/ogiusek/ioc/v2"
)

type CreatedMessage struct {
	endpoint.Message
	Added []Tile
}

type createdMessageEndpoint struct {
	Console console.Console `inject:"1"`
}

func (endpoint createdMessageEndpoint) Handle(req CreatedMessage) {
	text := fmt.Sprintf("created tiles %v\n", req)
	endpoint.Console.LogPermanentlyToConsole(text)
}

type DestroyedMessage struct {
	endpoint.Message
	Destroyed []Tile
}

type destroyedMessageEndpoint struct {
	Console console.Console `inject:"1"`
}

func (endpoint destroyedMessageEndpoint) Handle(req DestroyedMessage) {
	text := fmt.Sprintf("\ndestroyed tiles %v\n", req)
	endpoint.Console.LogPermanentlyToConsole(text)
}

//

type ClientPkg struct{}

func ClientPackage() ClientPkg {
	return ClientPkg{}
}

func (ClientPkg) Register(b ioc.Builder) {
	Package().Register(b)
	endpoint.MessageRegister[createdMessageEndpoint](b)
	endpoint.MessageRegister[destroyedMessageEndpoint](b)
}
