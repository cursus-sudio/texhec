package main

import (
	"backend"
	"backend/services/backendapi"
	"backend/services/clock"
	"backend/services/db"
	"backend/services/files"
	"errors"
	"fmt"
	"frontend"
	"frontend/example/ping"
	"frontend/example/tacticalmap"
	"frontend/services/backendconnector"
	"frontend/services/backendconnector/localconnector"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/ecs/ecsargs"
	"frontend/services/inputs"
	"frontend/services/scenes"
	"frontend/services/window"
	"os"
	"path/filepath"
	"time"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/null"
	"github.com/ogiusek/relay/v2"
)

type X struct {
	Shi []int
}

func main() {
	engineDir, err := os.Getwd()
	if err != nil {
		panic(errors.Join(errors.New("current wordking direcotry"), err))
	}
	// parent of both /backend and /frontend directory
	engineDir = filepath.Dir(filepath.Dir(engineDir))
	userStorage := filepath.Join(engineDir, "user_storage")

	var backendPkg backend.Pkg = backend.Package(
		clock.Package(time.RFC3339Nano),
		db.Package(
			fmt.Sprintf("%s/db.sql", userStorage),
			null.New(fmt.Sprintf("%s/backend/services/db/migrations", engineDir)),
		),
		files.Package(fmt.Sprintf("%s/files", userStorage)),
		[]ioc.Pkg{
			exBackendModPkg{},
			Package(),
		},
	)

	var pkg frontend.Pkg = frontend.Package(
		backendconnector.Package(localconnector.Package(backendPkg)),
		inputs.Package(),
		window.Package(),
	)

	b := ioc.NewBuilder()
	pkg.Register(b)
	c := b.Build()

	{ // pinging backend
		backend := ioc.Get[backendapi.Backend](c)
		r := backend.Relay()
		res, err := relay.Handle(r, ping.PingReq{ID: 2077})
		fmt.Printf("ping res is %v\nerr is %s\n", res, err)
	}
	{
		r := ioc.Get[backendapi.Backend](c).Relay()
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
		sceneManager := ioc.Get[scenes.SceneManager](c)

		world := ioc.Get[ecs.WorldFactory](c)()

		for i := 0; i < 1; i++ {
			entity := world.NewEntity()
			world.SaveComponent(entity, newSomeComponent())
		}

		someSystem := NewSomeSystem(
			sceneManager,
			world,
			ioc.Get[backendapi.Backend](c),
			ioc.Get[console.Console](c),
		)
		world.LoadSystem(&someSystem, ecs.DrawSystem)

		toggleSystem := NewToggledSystem(sceneManager, world, scenes.NewSceneId("main scene 2"))
		world.LoadSystem(&toggleSystem, ecs.UpdateSystem)

		sceneId := scenes.NewSceneId("main scene")
		mainScene := newMainScene(sceneId, world)
		sceneManager.AddScene(mainScene)
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
			ioc.Get[backendapi.Backend](c),
			ioc.Get[console.Console](c),
		)
		world.LoadSystem(&someSystem, ecs.DrawSystem)

		sceneId := scenes.NewSceneId("main scene 2")
		mainScene := newMainScene(sceneId, world)
		sceneManager.AddScene(mainScene)
	}

	var previousFrame time.Time
	previousFrame = time.Now()

	for { // runnning game loop
		world := ioc.Get[ecs.World](c)

		now := time.Now()
		deltaTime := ecsargs.NewDeltaTime(now.Sub(previousFrame))

		args := ecs.NewArgs(deltaTime)
		world.Update(args)

		previousFrame = now
		time.Sleep(previousFrame.Add(time.Millisecond * 16).Sub(now))
	}
}
