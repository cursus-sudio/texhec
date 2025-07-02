package main

import (
	"backend/services/clients"
	"backend/services/scopes"
	"fmt"
	"frontend"
	frontendapi "frontend/services/api"
	frontendtcp "frontend/services/api/tcp"
	"frontend/services/backendconnection"
	"frontend/services/backendconnection/localconnector"
	"frontend/services/console"
	"frontend/services/media"
	windowapi "frontend/services/media/window"
	"shared"
	"shared/services/api"
	"shared/services/api/netconnection"
	"shared/services/clock"
	"shared/services/logger"
	"shared/services/uuid"
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

func frontendDic(
	backendC ioc.Dic,
	netconnectionPkg netconnection.Pkg,
	clockPkg clock.Pkg,
) ioc.Dic {
	// defer sdl.Quit()
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(fmt.Errorf("Failed to initialize SDL: %s", err))
	}

	window, err := sdl.CreateWindow(
		"texhec",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600,
		sdl.WINDOW_SHOWN,
		// sdl.WINDOW_FULLSCREEN,
	)
	if err != nil {
		panic(fmt.Errorf("Failed to create window: %s", err))
	}
	// defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(fmt.Errorf("Failed to create renderer: %s", err))
	}
	// defer renderer.Destroy()
	renderer.Clear()
	renderer.SetDrawColor(0, 0, 255, 255)
	renderer.DrawRect(&sdl.Rect{X: 50, Y: 60, W: 100, H: 200})
	renderer.Present()

	var pkg frontend.Pkg = frontend.Package(
		shared.Package(
			api.Package(
				netconnectionPkg,
				func(c ioc.Dic) ioc.Dic { return c },
			),
			clockPkg,
			logger.Package(true, func(c ioc.Dic, message string) {
				ioc.Get[console.Console](c).LogPermanentlyToConsole(message)
			}),
		),
		frontendapi.Package(
			frontendtcp.Package("tcp"),
		),
		backendconnection.Package(
			localconnector.Package(func(clientCon connection.Connection) connection.Connection {
				backendC := backendC.Scope(scopes.UserSession)
				client := clients.NewClient(
					clients.ClientID(ioc.Get[uuid.Factory](backendC).NewUUID().String()),
					clientCon,
				)
				sClient := ioc.Get[clients.SessionClient](backendC)
				sClient.UseClient(client)

				return ioc.Get[connection.Connection](backendC)
			}),
		),
		media.Package(
			windowapi.Package(window, renderer),
		),
		[]ioc.Pkg{
			ClientPackage(),
		},
	)
	b := ioc.NewBuilder()
	pkg.Register(b)
	return b.Build()
}
