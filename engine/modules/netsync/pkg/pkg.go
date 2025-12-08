package netsyncpkg

import (
	"engine/modules/netsync"
	"engine/modules/netsync/internal/client"
	"engine/modules/netsync/internal/clienttypes"
	"engine/modules/netsync/internal/server"
	"engine/modules/netsync/internal/servertypes"
	"engine/modules/netsync/internal/state"
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
			// client
			Register(clienttypes.PredictedEvent{}).
			Register(clienttypes.FetchStateDTO{}).
			Register(clienttypes.EmitEventDTO{}).
			Register(clienttypes.TransparentEventDTO{}).
			// server
			Register(servertypes.SendStateDTO{}).
			Register(servertypes.SendChangeDTO{}).
			Register(servertypes.TransparentEventDTO{})
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
	ioc.RegisterSingleton(b, func(c ioc.Dic) netsync.StartSystem {
		clientToolFactory := ioc.Get[ecs.ToolFactory[client.Tool]](c)
		serverToolFactory := ioc.Get[ecs.ToolFactory[server.Tool]](c)
		return ecs.NewSystemRegister(func(w ecs.World) error {
			if pkg.config.config.IsClient {
				tool := clientToolFactory.Build(w)
				for _, listen := range tool.ListenToEvents {
					listen(w.EventsBuilder(), tool.BeforeEvent)
				}
				for _, listen := range tool.ListenToTransparentEvents {
					listen(w.EventsBuilder(), tool.OnTransparentEvent)
				}
			} else {
				tool := serverToolFactory.Build(w)
				for _, listen := range tool.ListenToEvents {
					listen(w.EventsBuilder(), tool.BeforeEvent)
				}
				for _, listen := range tool.ListenToTransparentEvents {
					listen(w.EventsBuilder(), tool.OnTransparentEvent)
				}
			}
			return nil
		})
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) netsync.StopSystem {
		clientToolFactory := ioc.Get[ecs.ToolFactory[client.Tool]](c)
		serverToolFactory := ioc.Get[ecs.ToolFactory[server.Tool]](c)
		return ecs.NewSystemRegister(func(w ecs.World) error {
			if pkg.config.config.IsClient {
				tool := clientToolFactory.Build(w)
				for _, listen := range tool.ListenToEvents {
					listen(w.EventsBuilder(), tool.AfterEvent)
				}
			} else {
				tool := serverToolFactory.Build(w)
				for _, listen := range tool.ListenToEvents {
					listen(w.EventsBuilder(), tool.AfterEvent)
				}
			}
			return nil
		})
	})
}
