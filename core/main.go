package main

import (
	"backend"
	"backend/services/clients"
	"backend/services/db"
	"backend/services/files"
	"backend/services/logger"
	"backend/services/scopes"
	"core/ping"
	"core/tacticalmap"
	"errors"
	"fmt"
	"frontend"
	"frontend/services/backendconnection"
	"frontend/services/backendconnection/localconnector"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/ecs/ecsargs"
	"frontend/services/inputs"
	"frontend/services/scenes"
	"frontend/services/window"
	"os"
	"path/filepath"
	"shared"
	"shared/services/clock"
	"shared/services/uuid"
	"shared/utils/connection"
	"shared/utils/endpoint"
	"time"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/null"
	"github.com/ogiusek/relay/v2"
)

func main() {
	engineDir, err := os.Getwd()
	if err != nil {
		panic(errors.Join(errors.New("current wordking direcotry"), err))
	}
	// parent of both /backend and /frontend directory
	engineDir = filepath.Dir(engineDir)
	userStorage := filepath.Join(engineDir, "user_storage")

	var sharedPkg shared.Pkg = shared.Package(
		clock.Package(time.RFC3339Nano),
	)

	var backendPkg backend.Pkg = backend.Package(
		sharedPkg,
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
		sharedPkg,
		backendconnection.Package(
			localconnector.Package(func(clientCon connection.Connection) connection.Connection {
				sC := sC.Scope(scopes.UserSession)
				client := clients.NewClient(
					clients.ClientID(ioc.Get[uuid.Factory](sC).NewUUID().String()),
					clientCon,
				)
				sClient := ioc.Get[clients.SessionClient](sC)
				sClient.UseClient(client)

				return ioc.Get[connection.Connection](sC)
			}),
		),
		inputs.Package(),
		window.Package(),
		[]ioc.Pkg{
			ClientPackage(),
		},
	)

	b := ioc.NewBuilder()
	pkg.Register(b)
	c := b.Build()

	{ // pinging backend
		backend := ioc.Get[backendconnection.Backend](c).Connection()
		r := backend.Relay()
		res, err := relay.Handle(r, ping.PingReq{Request: endpoint.NewRequest[ping.PingRes](), ID: 2077})
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
			ioc.Get[backendconnection.Backend](c).Connection(),
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
