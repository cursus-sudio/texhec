package netsyncpkg

import (
	"engine/modules/netsync"
	"engine/modules/netsync/internal/client"
	"engine/modules/netsync/internal/clienttypes"
	"engine/modules/netsync/internal/server"
	"engine/modules/netsync/internal/servertypes"
	"engine/modules/netsync/internal/tool"
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
	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
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

	ioc.RegisterSingleton(b, func(c ioc.Dic) netsync.ToolFactory {
		return tool.NewToolFactory()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[netsync.World, *server.Tool] {
		return server.NewToolFactory(
			*pkg.config.config,
			ioc.Get[netsync.ToolFactory](c),
			ioc.Get[logger.Logger](c),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[netsync.World, *client.Tool] {
		return client.NewToolFactory(
			*pkg.config.config,
			ioc.Get[netsync.ToolFactory](c),
			ioc.Get[logger.Logger](c),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) netsync.StartSystem {
		clientToolFactory := ioc.Get[ecs.ToolFactory[netsync.World, *client.Tool]](c)
		serverToolFactory := ioc.Get[ecs.ToolFactory[netsync.World, *server.Tool]](c)
		return ecs.NewSystemRegister(func(w netsync.World) error {
			clientTool := clientToolFactory.Build(w)
			for _, listen := range clientTool.ListenToEvents {
				listen(w.EventsBuilder(), clientTool.BeforeEvent)
			}
			for _, listen := range clientTool.ListenToSimulatedEvents {
				listen(w.EventsBuilder(), clientTool.BeforeEventRecord)
			}
			for _, listen := range clientTool.ListenToTransparentEvents {
				listen(w.EventsBuilder(), clientTool.OnTransparentEvent)
			}

			serverTool := serverToolFactory.Build(w)
			for _, listen := range serverTool.ListenToEvents {
				listen(w.EventsBuilder(), serverTool.BeforeEvent)
			}
			for _, listen := range clientTool.ListenToSimulatedEvents {
				listen(w.EventsBuilder(), serverTool.BeforeEvent)
			}
			for _, listen := range serverTool.ListenToTransparentEvents {
				listen(w.EventsBuilder(), serverTool.OnTransparentEvent)
			}
			return nil
		})
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) netsync.StopSystem {
		clientToolFactory := ioc.Get[ecs.ToolFactory[netsync.World, *client.Tool]](c)
		serverToolFactory := ioc.Get[ecs.ToolFactory[netsync.World, *server.Tool]](c)
		return ecs.NewSystemRegister(func(w netsync.World) error {
			clientTool := clientToolFactory.Build(w)
			for _, listen := range clientTool.ListenToEvents {
				listen(w.EventsBuilder(), clientTool.AfterEvent)
			}
			for _, listen := range clientTool.ListenToSimulatedEvents {
				listen(w.EventsBuilder(), clientTool.AfterEvent)
			}

			serverTool := serverToolFactory.Build(w)
			for _, listen := range serverTool.ListenToEvents {
				listen(w.EventsBuilder(), serverTool.AfterEvent)
			}
			for _, listen := range serverTool.ListenToSimulatedEvents {
				listen(w.EventsBuilder(), serverTool.AfterEvent)
			}
			return nil
		})
	})
}
