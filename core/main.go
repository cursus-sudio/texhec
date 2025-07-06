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
	"runtime"
	appruntime "shared/services/runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

// TODO
// set different events for scene and for rest

func main() {
	print("started\n")
	runtime.LockOSThread()

	isServer := false
	for _, arg := range os.Args {
		if arg == "server" {
			isServer = true
			break
		}
	}

	sharedPkg := SharedPackage()

	backendC := backendDic(sharedPkg)

	if isServer {
		backendRuntime := ioc.Get[appruntime.Runtime](backendC)
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
	{
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	}

	frontendRuntime := ioc.Get[appruntime.Runtime](c)
	// go func() {
	// 	time.Sleep(time.Second / 10)
	// 	frontendRuntime.Stop()
	// }()
	frontendRuntime.Run()
}
