package main

import (
	"core/ping"
	"core/tacticalmap"
	"fmt"
	frontendtcp "frontend/services/api/tcp"
	"frontend/services/backendconnection"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/scenes"
	frontendscopes "frontend/services/scopes"
	"os"
	"shared/services/api/netconnection"
	"shared/services/clock"
	"shared/services/runtime"
	"time"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

// TODO
// set different events for scene and for rest

func main() {
	print("started\n")
	isServer := false
	for _, arg := range os.Args {
		if arg == "server" {
			isServer = true
			break
		}
	}

	sharedPkg := SharedPackage(
		netconnection.Package(time.Second),
		clock.Package(time.RFC3339Nano),
	)

	backendC := backendDic(sharedPkg)

	if isServer {
		backendRuntime := ioc.Get[runtime.Runtime](backendC)
		backendRuntime.Run()
		return
	}

	c := frontendDic(
		backendC,
		sharedPkg,
	)

	{ // connect
		tcpConnect := ioc.Get[frontendtcp.Connect](c)
		err := tcpConnect.Connect("localhost:8080")
		if err != nil {
			panic(err)
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
		mainScene := newMainScene(sceneId, events.NewBuilder().Build(), world)
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
		mainScene := newMainScene(sceneId, events.NewBuilder().Build(), world)
		sceneManager.AddScene(mainScene)
	}

	frontendRuntime := ioc.Get[runtime.Runtime](c)
	frontendRuntime.Run()
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
