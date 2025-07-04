package main

import (
	"backend/services/clients"
	backendscopes "backend/services/scopes"
	"fmt"
	frontendapi "frontend/services/api"
	frontendtcp "frontend/services/api/tcp"
	"frontend/services/backendconnection"
	"frontend/services/backendconnection/localconnector"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/frames"
	"frontend/services/media"
	windowapi "frontend/services/media/window"
	"frontend/services/scenes"
	frontendscopes "frontend/services/scopes"
	"shared/services/api"
	"shared/services/logger"
	"shared/services/uuid"
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

func frontendDic(
	backendC ioc.Dic,
	sharedPkg SharedPkg,
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

	pkgs := []ioc.Pkg{
		sharedPkg,
		api.Package(func(c ioc.Dic) ioc.Dic { return c }),
		logger.Package(true, func(c ioc.Dic, message string) {
			ioc.Get[console.Console](c).LogPermanentlyToConsole(message)
		}),
		frontendtcp.Package("tcp"),
		frontendapi.Package(),
		localconnector.Package(func(clientCon connection.Connection) connection.Connection {
			backendC := backendC.Scope(backendscopes.UserSession)
			client := clients.NewClient(
				clients.ClientID(ioc.Get[uuid.Factory](backendC).NewUUID().String()),
				clientCon,
			)
			sClient := ioc.Get[clients.SessionClient](backendC)
			sClient.UseClient(client)

			return ioc.Get[connection.Connection](backendC)
		}),
		backendconnection.Package(func(c ioc.Dic) connection.Connection {
			return ioc.Get[localconnector.Connector](c).Connect()
		}),
		console.Package(),
		media.Package(
			windowapi.Package(window, renderer),
		),
		ecs.Package(),
		frames.Package(),
		scenes.Package(),
		frontendscopes.Package(),

		// mods
		ClientPackage(),
	}

	b := ioc.NewBuilder()
	for _, pkg := range pkgs {
		pkg.Register(b)
	}
	// pkg.Register(b)
	return b.Build()
}
