package netsyncpkg

import (
	"engine/modules/connection"
	"engine/modules/netsync"
	"engine/modules/netsync/internal/client"
	"engine/modules/netsync/internal/clienttypes"
	"engine/modules/netsync/internal/server"
	"engine/modules/netsync/internal/servertypes"
	"engine/modules/netsync/internal/service"
	"engine/modules/record"
	"engine/modules/uuid"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/events"
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

	ioc.RegisterSingleton(b, func(c ioc.Dic) netsync.Service {
		return service.NewToolFactory(
			ioc.Get[ecs.World](c),
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) *server.Service {
		return server.NewService(
			*pkg.config.config,
			ioc.Get[logger.Logger](c),
			ioc.Get[events.Builder](c),
			ioc.Get[ecs.World](c),
			ioc.Get[netsync.Service](c),
			ioc.Get[connection.Service](c),
			ioc.Get[record.Service](c),
			ioc.Get[uuid.Service](c),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) *client.Service {
		return client.NewService(
			*pkg.config.config,
			ioc.Get[logger.Logger](c),
			ioc.Get[events.Builder](c),
			ioc.Get[ecs.World](c),
			ioc.Get[netsync.Service](c),
			ioc.Get[connection.Service](c),
			ioc.Get[record.Service](c),
			ioc.Get[uuid.Service](c),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) netsync.StartSystem {
		clientService := ioc.Get[*client.Service](c)
		serverService := ioc.Get[*server.Service](c)
		eventsBuilder := ioc.Get[events.Builder](c)
		return ecs.NewSystemRegister(func() error {
			for _, listen := range clientService.ListenToEvents {
				listen(eventsBuilder, clientService.BeforeEvent)
			}
			for _, listen := range clientService.ListenToSimulatedEvents {
				listen(eventsBuilder, clientService.BeforeEventRecord)
			}
			for _, listen := range clientService.ListenToTransparentEvents {
				listen(eventsBuilder, clientService.OnTransparentEvent)
			}

			for _, listen := range serverService.ListenToEvents {
				listen(eventsBuilder, serverService.BeforeEvent)
			}
			for _, listen := range clientService.ListenToSimulatedEvents {
				listen(eventsBuilder, serverService.BeforeEvent)
			}
			for _, listen := range serverService.ListenToTransparentEvents {
				listen(eventsBuilder, serverService.OnTransparentEvent)
			}
			return nil
		})
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) netsync.StopSystem {
		clientService := ioc.Get[*client.Service](c)
		serverService := ioc.Get[*server.Service](c)
		eventsBuilder := ioc.Get[events.Builder](c)
		return ecs.NewSystemRegister(func() error {
			for _, listen := range clientService.ListenToEvents {
				listen(eventsBuilder, clientService.AfterEvent)
			}
			for _, listen := range clientService.ListenToSimulatedEvents {
				listen(eventsBuilder, clientService.AfterEvent)
			}

			for _, listen := range serverService.ListenToEvents {
				listen(eventsBuilder, serverService.AfterEvent)
			}
			for _, listen := range serverService.ListenToSimulatedEvents {
				listen(eventsBuilder, serverService.AfterEvent)
			}
			return nil
		})
	})
}
