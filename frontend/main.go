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
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/ogiusek/ioc"
)

// game loop

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(errors.Join(errors.New("current wordking direcotry"), err))
	}
	// parent of both /backend and /frontend directory
	currentDir = filepath.Dir(currentDir)
	userStorage := filepath.Join(currentDir, "user_storage")

	var backendPkgs []ioc.Pkg = []ioc.Pkg{
		backendsrc.Package(
			utils.Package(
				clock.Package(time.RFC3339Nano),
				db.Package(
					fmt.Sprintf("%s/db.sql", userStorage),
					// fmt.Sprintf("%s/migrations", userStorage),
					fmt.Sprintf("%s/backend/db/migrations", currentDir),
				),
				files.Package(fmt.Sprintf("%s/files", userStorage)),
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
	}

	var frontendPkgs []ioc.Pkg = []ioc.Pkg{
		frontendsrc.Package(
			engine.Package(
				inputs.Package(),
				window.Package(),
			),
		),
	}

	bC := ioc.NewContainer()
	fC := ioc.NewContainer()

	for _, pkg := range backendPkgs {
		pkg.Register(bC)
		pkg.Register(fC) // temporary until mediator isn't created
	}

	for _, pkg := range frontendPkgs {
		pkg.Register(fC)
	}

	{
		repository := NewIntRepository(
			0,
			false,
			ioc.Get[saves.StateCodecRWMutex](bC).RWMutex().RLocker(),
		)
		var intRepo IntRepo = nil
		repoId := reflect.TypeOf(&intRepo).Elem().String()
		print(repoId)
		repositories := ioc.Get[saves.SavableRepositories](bC)
		repositories.AddRepo(saves.RepoId(repoId), repository)

		ioc.RegisterSingleton(bC, func(c ioc.Dic) IntRepo { return repository })
	}

	{
		fmt.Print("saving\n")
		saveMetaFactory := ioc.Get[saves.SaveMetaFactory](bC)
		savesService := ioc.Get[saves.Saves](bC)

		builder := ioc.Get[saves.ListSavesQueryBuilder](bC)
		{
			metas, err := savesService.ListSaves(
				builder.SavesPerPage(100).Build(),
			)
			if err != nil {
				panic(fmt.Sprintf("queried saves to delete them %s", err.Error()))
			} else {
				for _, meta := range metas {
					savesService.Delete(meta.Id)
				}
			}
		}

		repo := ioc.Get[IntRepo](bC)
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

		metas, err := savesService.ListSaves(
			builder.SavesPerPage(100).Build(),
		)
		if err != nil {
			print(err.Error())
			print("\nthis were error in case you didn't notice\n")
		} else {
			for i, meta := range metas {
				fmt.Printf("meta %d: %s\n", i, meta)
			}
		}

		scopeCleanUp := ioc.Get[scopecleanup.ScopeCleanUp](bC)
		scopeCleanUp.Clean(scopecleanup.NewCleanUpArgs(nil))
	}

	{ // adding scene 1
		sceneManager := ioc.Get[scenes.SceneManager](fC)

		world := ioc.Get[ecs.WorldFactory](fC)()

		for i := 0; i < 1; i++ {
			entity := world.NewEntity()
			world.SaveComponent(entity, newSomeComponent())
		}

		someSystem := NewSomeSystem(
			sceneManager,
			world,
			ioc.Get[backendapi.Backend](fC),
			ioc.Get[console.Console](fC),
		)
		world.LoadSystem(&someSystem, ecs.DrawSystem)

		toggleSystem := NewToggledSystem(sceneManager, world, scenes.NewSceneId("main scene 2"))
		world.LoadSystem(&toggleSystem, ecs.UpdateSystem)

		sceneId := scenes.NewSceneId("main scene")
		mainScene := newMainScene(sceneId, world)
		sceneManager.AddScene(mainScene)
	}
	{ // adding scene 2
		sceneManager := ioc.Get[scenes.SceneManager](fC)

		world := ioc.Get[ecs.WorldFactory](fC)()

		for i := 0; i < 2; i++ {
			entity := world.NewEntity()
			world.SaveComponent(entity, newSomeComponent())
		}

		someSystem := NewSomeSystem(
			sceneManager,
			world,
			ioc.Get[backendapi.Backend](fC),
			ioc.Get[console.Console](fC),
		)
		world.LoadSystem(&someSystem, ecs.DrawSystem)

		sceneId := scenes.NewSceneId("main scene 2")
		mainScene := newMainScene(sceneId, world)
		sceneManager.AddScene(mainScene)
	}

	var previousFrame time.Time
	previousFrame = time.Now()

	for { // runnning game loop
		world := ioc.Get[ecs.World](fC)

		now := time.Now()
		deltaTime := ecsargs.NewDeltaTime(now.Sub(previousFrame))

		args := ecs.NewArgs(deltaTime)
		world.Update(args)

		previousFrame = now
		time.Sleep(previousFrame.Add(time.Millisecond * 16).Sub(now))
	}
}
