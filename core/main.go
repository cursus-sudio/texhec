package main

import (
	"backend"
	backendapi "backend/services/api"
	backendtcp "backend/services/api/tcp"
	"backend/services/clients"
	"backend/services/db"
	"backend/services/files"
	"backend/services/logger"
	backendscopes "backend/services/scopes"
	"core/ping"
	"core/tacticalmap"
	"errors"
	"fmt"
	"frontend"
	frontendapi "frontend/services/api"
	frontendtcp "frontend/services/api/tcp"
	"frontend/services/backendconnection"
	"frontend/services/backendconnection/localconnector"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/scenes"
	frontendscopes "frontend/services/scopes"
	"os"
	"path/filepath"
	"shared"
	"shared/services/api"
	"shared/services/clock"
	"shared/services/runtime"
	"shared/services/uuid"
	"shared/utils/connection"
	"time"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/null"
	"github.com/ogiusek/relay/v2"
)

func main() {
	print("started\n")
	isServer := false
	for _, arg := range os.Args {
		if arg == "server" {
			isServer = true
			break
		}
	}

	// defer sdl.Quit()
	//
	// if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
	// 	log.Fatalf("Failed to initialize SDL: %s", err)
	// }
	//
	// window, err := sdl.CreateWindow(
	// 	"SDL2 Go Example",
	// 	sdl.WINDOWPOS_UNDEFINED,
	// 	sdl.WINDOWPOS_UNDEFINED,
	// 	800,
	// 	600,
	// 	sdl.WINDOW_SHOWN,
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to create window: %s", err)
	// }
	// defer window.Destroy()
	//
	// // window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
	//
	// renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	// if err != nil {
	// 	log.Fatalf("Failed to create renderer: %s", err)
	// }
	// defer renderer.Destroy()
	// renderer.SetDrawColor(0, 0, 255, 255)
	// renderer.Clear()
	// renderer.Present()
	// time.Sleep(3 * time.Second)
	//
	// // ---------------------------------------------------------------------------------
	// if true {
	// 	return
	// }

	engineDir, err := os.Getwd()
	if err != nil {
		panic(errors.Join(errors.New("current wordking direcotry"), err))
	}
	// parent of both /backend and /frontend directory
	engineDir = filepath.Dir(engineDir)
	userStorage := filepath.Join(engineDir, "user_storage")

	var clockPkg = clock.Package(time.RFC3339Nano)

	var backendPkg backend.Pkg = backend.Package(
		shared.Package(
			api.Package(func(c ioc.Dic) ioc.Dic { return c.Scope(backendscopes.Request) }),
			clockPkg,
		),
		backendapi.Package(
			backendtcp.Package(
				"0.0.0.0",
				"8080",
				"tcp",
			),
		),
		db.Package(
			fmt.Sprintf("%s/db.sql", userStorage),
			null.New(fmt.Sprintf("%s/engine/backend/services/db/migrations", engineDir)),
		),
		files.Package(fmt.Sprintf("%s/files", userStorage)),
		logger.Package(true),
		[]ioc.Pkg{
			exBackendModPkg{},
			ServerPackage(),
		},
	)

	sCB := ioc.NewBuilder()
	backendPkg.Register(sCB)

	sC := sCB.Build()

	var pkg frontend.Pkg = frontend.Package(
		shared.Package(
			api.Package(func(c ioc.Dic) ioc.Dic { return c }),
			clockPkg,
		),
		frontendapi.Package(
			frontendtcp.Package("tcp"),
		),
		backendconnection.Package(
			localconnector.Package(func(clientCon connection.Connection) connection.Connection {
				sC := sC.Scope(backendscopes.UserSession)
				client := clients.NewClient(
					clients.ClientID(ioc.Get[uuid.Factory](sC).NewUUID().String()),
					clientCon,
				)
				sClient := ioc.Get[clients.SessionClient](sC)
				sClient.UseClient(client)

				return ioc.Get[connection.Connection](sC)
			}),
		),
		[]ioc.Pkg{
			ClientPackage(),
		},
	)

	b := ioc.NewBuilder()
	pkg.Register(b)
	c := b.Build()

	{ // connect
		if !isServer {
			tcpConnect := ioc.Get[frontendtcp.Connect](c)
			err := tcpConnect.Connect("localhost:8080")
			if err != nil {
				panic(err)
			}
		}
	}
	{ // pinging backend
		backend := ioc.Get[backendconnection.Backend](c).Connection()
		r := backend.Relay()
		res, err := relay.Handle(r, ping.PingReq{ID: 2077})
		fmt.Printf("client recieved ping res is %v\nerr is %s\n", res, err)
	}
	{
		r := ioc.Get[backendconnection.Backend](c).Connection().Relay()
		res, err := relay.Handle(r, tacticalmap.NewCreateReq(
			tacticalmap.CreateArgs{
				Tiles: []tacticalmap.Tile{
					{Pos: tacticalmap.Pos{X: 7, Y: 13}},
				},
			},
		))
		fmt.Printf("create res is %v\nerr is %s\n", res, err)
	}

	{ // adding scene 1
		c := c.Scope(frontendscopes.Scene)
		sceneManager := ioc.Get[scenes.SceneManager](c)

		world := ioc.Get[ecs.WorldFactory](c)()

		for i := 0; i < 1; i++ {
			entity := world.NewEntity()
			world.SaveComponent(entity, newSomeComponent())
		}

		someSystem := NewSomeSystem(
			sceneManager,
			world,
			ioc.Get[backendconnection.Backend](c).Connection(),
			ioc.Get[console.Console](c),
		)
		world.LoadSystem(&someSystem, ecs.DrawSystem)

		toggleSystem := NewToggledSystem(sceneManager, world, scenes.NewSceneId("main scene 2"))
		world.LoadSystem(&toggleSystem, ecs.UpdateSystem)

		sceneId := scenes.NewSceneId("main scene")
		mainScene := newMainScene(sceneId, ioc.Get[scenes.SceneEvents](c), world)
		sceneManager.AddScene(mainScene)
		sceneManager.LoadScene(mainScene.Id())
	}
	{ // adding scene 2
		sceneManager := ioc.Get[scenes.SceneManager](c)

		world := ioc.Get[ecs.WorldFactory](c)()

		for i := 0; i < 2; i++ {
			entity := world.NewEntity()
			world.SaveComponent(entity, newSomeComponent())
		}

		someSystem := NewSomeSystem(
			sceneManager,
			world,
			ioc.Get[backendconnection.Backend](c).Connection(),
			ioc.Get[console.Console](c),
		)
		world.LoadSystem(&someSystem, ecs.DrawSystem)

		sceneId := scenes.NewSceneId("main scene 2")
		mainScene := newMainScene(sceneId, ioc.Get[scenes.SceneEvents](c), world)
		sceneManager.AddScene(mainScene)
	}

	closeChan := make(chan struct{})
	backendRuntime := ioc.Get[runtime.Runtime](sC)
	frontendRuntime := ioc.Get[runtime.Runtime](c)
	go func() {
		if isServer {
			backendRuntime.Run()
			closeChan <- struct{}{}
		}
	}()
	go func() {
		if !isServer {
			frontendRuntime.Run()
			closeChan <- struct{}{}
		}
	}()
	<-closeChan
	backendRuntime.Stop()
	frontendRuntime.Stop()
}
