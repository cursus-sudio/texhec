package syncpkg

import (
	"engine/modules/sync"
	"engine/modules/sync/internal/client"
	"engine/modules/sync/internal/clienttypes"
	"engine/modules/sync/internal/server"
	"engine/modules/sync/internal/servertypes"
	"engine/modules/sync/internal/state"
	"engine/modules/uuid"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	config Config
}

func Package(config Config) ioc.Pkg {
	return pkg{
		config,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		return b.
			Register(clienttypes.PredictedEvent{}).
			Register(clienttypes.FetchStateDTO{}).
			Register(clienttypes.EmitEventDTO{}).
			Register(servertypes.SendStateDTO{}).
			Register(servertypes.SendChangeDTO{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[state.Tool] {
		return state.NewToolFactory(
			*pkg.config.config,
			ioc.Get[ecs.ToolFactory[uuid.Tool]](c),
			ioc.Get[logger.Logger](c),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[server.Tool] {
		return server.NewToolFactory(
			*pkg.config.config,
			ioc.Get[ecs.ToolFactory[state.Tool]](c),
			ioc.Get[ecs.ToolFactory[uuid.Tool]](c),
			ioc.Get[logger.Logger](c),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[client.Tool] {
		return client.NewToolFactory(
			*pkg.config.config,
			ioc.Get[ecs.ToolFactory[state.Tool]](c),
			ioc.Get[ecs.ToolFactory[uuid.Tool]](c),
			ioc.Get[logger.Logger](c),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) sync.StartSystem {
		clientToolFactory := ioc.Get[ecs.ToolFactory[client.Tool]](c)
		serverToolFactory := ioc.Get[ecs.ToolFactory[server.Tool]](c)
		return ecs.NewSystemRegister(func(w ecs.World) error {
			if pkg.config.config.IsClient {
				tool := clientToolFactory.Build(w)
				for _, listen := range tool.ListenToEvents {
					listen(w.EventsBuilder(), tool.BeforeInternalEvent)
				}
			} else {
				tool := serverToolFactory.Build(w)
				for _, listen := range tool.ListenToEvents {
					listen(w.EventsBuilder(), tool.BeforeInternalEvent)
				}
			}
			return nil
		})
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) sync.StopSystem {
		clientToolFactory := ioc.Get[ecs.ToolFactory[client.Tool]](c)
		serverToolFactory := ioc.Get[ecs.ToolFactory[server.Tool]](c)
		return ecs.NewSystemRegister(func(w ecs.World) error {
			if pkg.config.config.IsClient {
				tool := clientToolFactory.Build(w)
				for _, listen := range tool.ListenToEvents {
					listen(w.EventsBuilder(), tool.AfterInternalEvent)
				}
			} else {
				tool := serverToolFactory.Build(w)
				for _, listen := range tool.ListenToEvents {
					listen(w.EventsBuilder(), tool.AfterInternalEvent)
				}
			}
			return nil
		})
	})
}
