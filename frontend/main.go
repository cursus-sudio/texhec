package main

import (
	backendsrc "backend/src"
	"backend/src/backendapi"
	"backend/src/backendapi/backendapipkg"
	"backend/src/backendapi/ping"
	"backend/src/backendapi/tacticalmapapi"
	backendmodules "backend/src/modules"
	"backend/src/modules/saves"
	"backend/src/modules/tacticalmap"
	"backend/src/utils"
	"backend/src/utils/clock"
	"backend/src/utils/db"
	"backend/src/utils/files"
	"backend/src/utils/logger"
	"backend/src/utils/services"
	"backend/src/utils/services/scopecleanup"
	"backend/src/utils/uuid"
	"errors"
	"fmt"
	frontendsrc "frontend/src"
	"frontend/src/engine"
	"frontend/src/engine/console"
	"frontend/src/engine/ecs"
	"frontend/src/engine/ecs/ecsargs"
	"frontend/src/engine/inputs"
	"frontend/src/engine/scenes"
	"frontend/src/engine/window"
	frontendmodules "frontend/src/modules"
	"frontend/src/modules/backendconnector"
	"frontend/src/modules/backendconnector/localconnector"
	"os"
	"path/filepath"
	"time"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type X struct {
	Shi []int
}

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(errors.Join(errors.New("current wordking direcotry"), err))
	}
	// parent of both /backend and /frontend directory
	currentDir = filepath.Dir(currentDir)
	userStorage := filepath.Join(currentDir, "user_storage")

	var backendPkg backendsrc.Pkg = backendsrc.Package(
		utils.Package(
			clock.Package(time.RFC3339Nano),
			db.Package(
				fmt.Sprintf("%s/db.sql", userStorage),
				fmt.Sprintf("%s/backend/db/migrations", currentDir),
			),
			files.Package(fmt.Sprintf("%s/files", userStorage)),
			logger.Package(),
			services.Package(
				scopecleanup.Package(),
			),
			uuid.Package(),
		),
		backendmodules.Package(
			saves.Package(),
			tacticalmap.Package(),
		),
		backendapipkg.Package(),
		[]ioc.Pkg{
			exBackendModPkg{},
		},
	)

	var pkg frontendsrc.Pkg = frontendsrc.Package(
		engine.Package(
			inputs.Package(),
			window.Package(),
		),
		frontendmodules.Package(
			backendconnector.Package(
				localconnector.Package(backendPkg),
			),
		),
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
		res, err := relay.Handle(r, tacticalmapapi.NewCreateReq(
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
