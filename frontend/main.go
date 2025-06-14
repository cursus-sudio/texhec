package main

import (
	backendsrc "backend/src"
	"backend/src/backendapi"
	"backend/src/modules"
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
	"fmt"
	frontendsrc "frontend/src"
	"frontend/src/engine"
	"frontend/src/engine/console"
	"frontend/src/engine/ecs"
	"frontend/src/engine/ecs/ecsargs"
	"frontend/src/engine/inputs"
	"frontend/src/engine/scenes"
	"frontend/src/engine/window"
	"reflect"
	"time"

	"github.com/ogiusek/ioc"
)

var pkgs []ioc.Pkg = []ioc.Pkg{
	backendsrc.Package(
		utils.Package(
			clock.Package(time.RFC3339Nano),
			db.Package(
				"/home/ogius/Documents/code/fullstack/texhec_7/backend/db/db.sql",
				"/home/ogius/Documents/code/fullstack/texhec_7/backend/db/migrations",
			),
			files.Package("/home/ogius/Documents/code/fullstack/texhec_7/backend/db/files"),
			logger.Package(),
			services.Package(
				scopecleanup.Package(),
			),
			uuid.Package(),
		),
		modules.Package(
			saves.Package(),
			tacticalmap.Package(),
		),
		backendapi.Package(),
	),
	frontendsrc.Package(
		engine.Package(
			inputs.Package(),
			window.Package(),
		),
	),
}

// game loop

func main() {
	c := ioc.NewContainer()

	for _, pkg := range pkgs {
		pkg.Register(c)
	}

	{
		repository := NewIntRepository(
			0,
			false,
			ioc.Get[saves.StateCodecRWMutex](c).RWMutex().RLocker(),
		)
		var intRepo IntRepo = nil
		repoId := reflect.TypeOf(&intRepo).Elem().String()
		print(repoId)
		repositories := ioc.Get[saves.SavableRepositories](c)
		repositories.AddRepo(saves.RepoId(repoId), repository)

		ioc.RegisterSingleton(c, func(c ioc.Dic) IntRepo { return repository })
	}

	{
		fmt.Print("saving\n")
		saveMetaFactory := ioc.Get[saves.SaveMetaFactory](c)
		savesService := ioc.Get[saves.Saves](c)
		repo := ioc.Get[IntRepo](c)
		for i := 0; i < 10; i++ {
			repo.Increment()
		}
		fmt.Printf("count s1 is %d\n", repo.GetCount())
		s1 := saveMetaFactory.New(saves.SaveName("s1"))
		if err := savesService.NewSave(s1); err != nil {
			panic(err)
		}
		for i := 0; i < 10; i++ {
			repo.Increment()
		}
		fmt.Printf("count s2 is %d\n", repo.GetCount())
		s2 := saveMetaFactory.New(saves.SaveName("s2"))
		if err := savesService.NewSave(s2); err != nil {
			panic(err)
		}
		fmt.Print("loading\n")
		if err := savesService.Load(s1.Id); err != nil {
			panic(err)
		}
		fmt.Printf("count s1 is %d\n", repo.GetCount())
		if err := savesService.Load(s2.Id); err != nil {
			panic(err)
		}
		fmt.Printf("count s2 is %d\n", repo.GetCount())

		builder := ioc.Get[saves.ListSavesQueryBuilder](c)
		metas, err := savesService.ListSaves(builder.Build())
		if err != nil {
			print(err)
		} else {
			for i, meta := range metas {
				fmt.Printf("meta %d: %s\n", i, meta)
			}
		}

		scopeCleanUp := ioc.Get[scopecleanup.ScopeCleanUp](c)
		scopeCleanUp.Clean(scopecleanup.NewCleanUpArgs(nil))
	}
	return

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
