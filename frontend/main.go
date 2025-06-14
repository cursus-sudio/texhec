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
	"encoding/json"
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
	"reflect"
	"sync"
	"time"

	"github.com/ogiusek/ioc"
)

var pkgs []ioc.Pkg = []ioc.Pkg{
	backendsrc.Package(
		utils.Package(
			clock.Package(time.RFC3339Nano),
			db.Package("/home/ogius/Documents/code/fullstack/texhec_7/frontend/data/db.sql"),
			files.Package("/home/ogius/Documents/code/fullstack/texhec_7/frontend/data/saves"),
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

// component

type someComponent struct {
	StartedAt time.Time
	Frame     int
	PastTime  time.Duration
}

func newSomeComponent() someComponent {
	return someComponent{
		StartedAt: time.Now(),
		Frame:     0,
		PastTime:  time.Duration(0),
	}
}

// system

type toggleSystem struct {
	SceneManager scenes.SceneManager
	World        ecs.World
	Toggled      bool
	Scene2       scenes.SceneId
}

func NewToggledSystem(
	sceneManager scenes.SceneManager,
	world ecs.World,
	scene2 scenes.SceneId,
) toggleSystem {
	return toggleSystem{
		SceneManager: sceneManager,
		World:        world,
		Toggled:      false,
		Scene2:       scene2,
	}
}

func (system *toggleSystem) Update(args ecs.Args) {
	if system.Toggled {
		return
	}

	for _, entity := range system.World.GetEntitiesWithComponents(
		ecs.GetComponentType(someComponent{}),
	) {
		component := someComponent{}
		if err := system.World.GetComponent(entity, &component); err != nil {
			continue
		}
		if component.PastTime.Seconds() < 5 {
			return
		}
	}

	system.SceneManager.LoadScene(system.Scene2)
	system.Toggled = true
}

type someSystem struct {
	SceneManager scenes.SceneManager
	World        ecs.World
	Backend      backendapi.Backend
	Console      console.Console
}

func NewSomeSystem(
	sceneMagener scenes.SceneManager,
	world ecs.World,
	backend backendapi.Backend,
	console console.Console,
) someSystem {
	return someSystem{
		SceneManager: sceneMagener,
		World:        world,
		Backend:      backend,
		Console:      console,
	}
}

func (system *someSystem) Update(args ecs.Args) {
	format := "02-01-2006 15:04:05"
	text := ""
	text += fmt.Sprintf("now %s\n", time.Now().Format(format))

	for _, entity := range system.World.GetEntitiesWithComponents(ecs.GetComponentType(someComponent{})) {
		component := someComponent{}
		if err := system.World.GetComponent(entity, &component); err != nil {
			continue
		}
		component.PastTime += args.DeltaTime.Duration()
		component.Frame += 1
		system.World.SaveComponent(entity, component)

		text += "\n"
		text += fmt.Sprintf("first frame %s\n", component.StartedAt.Format(format))
		text += fmt.Sprintf("current frame %d\n", component.Frame)
		// text += fmt.Sprintf("time in game %d\n", int(component.PastTime.Seconds()))
		text += fmt.Sprintf("time in game %f\n", component.PastTime.Seconds())
	}
	text += "\n"

	tacticalMap, _ := system.Backend.TacticalMap().GetMap()
	text += fmt.Sprintf("found shit %v\n", tacticalMap)

	text += "entities found: {\n"
	for _, entity := range system.World.GetEntities() {
		text += fmt.Sprintf(" - %s\n", entity)
	}
	text += "}\n"
	text += "\n"

	system.Console.ClearAndLogToConsole(text)
}

// scene

type mainScene struct {
	sceneId scenes.SceneId
	world   ecs.World
}

func newMainScene(sceneId scenes.SceneId, world ecs.World) scenes.Scene {
	return &mainScene{sceneId: sceneId, world: world}
}

func (mainScene *mainScene) Id() scenes.SceneId { return mainScene.sceneId }
func (mainScene *mainScene) Load()              {}
func (mainScene *mainScene) Unload()            {}
func (mainScene *mainScene) World() ecs.World   { return mainScene.world }

// repo

type IntRepo interface {
	Increment()
	GetCount() int
}

type intRepo struct {
	Count   int         `json:"count"`
	Changed bool        `json:"-"`
	Mutex   sync.Locker `json:"-"`
}

func NewIntRepository(
	count int,
	changed bool,
	mutex sync.Locker,
) *intRepo {
	return &intRepo{
		Count:   count,
		Changed: changed,
		Mutex:   mutex,
	}
}

func (savableIntRepo *intRepo) IsValidSnapshot(snapshot saves.RepoSnapshot) bool {
	repo := intRepo{}
	err := json.Unmarshal(snapshot.Bytes(), &repo)
	return err != nil
}

func (savableIntRepo *intRepo) TakeSnapshot() saves.RepoSnapshot {
	savableIntRepo.Mutex.Lock()
	defer savableIntRepo.Mutex.Unlock()
	bytes, err := json.Marshal(savableIntRepo)
	if err != nil {
		panic(errors.Join(errors.New("error taking snapshot"), err))
	}
	savableIntRepo.Changed = false
	return saves.NewRepoSnapshot(bytes)
}

func (savableIntRepo *intRepo) LoadSnapshot(snapshot saves.RepoSnapshot) error {
	savableIntRepo.Mutex.Lock()
	defer savableIntRepo.Mutex.Unlock()
	repo := intRepo{}
	err := json.Unmarshal(snapshot.Bytes(), &repo)
	if err != nil {
		return err
	}
	savableIntRepo.Count = repo.Count
	savableIntRepo.Changed = false
	return nil
}

func (savableIntRepo *intRepo) HasChanges() bool {
	return savableIntRepo.Changed
}

func (savableIntRepo *intRepo) Increment() {
	savableIntRepo.Mutex.Lock()
	defer savableIntRepo.Mutex.Unlock()
	savableIntRepo.Count += 1
	savableIntRepo.Changed = true
}

func (savableIntRepo *intRepo) GetCount() int {
	return savableIntRepo.Count
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
		fmt.Print("saving")
		saveMetaFactory := ioc.Get[saves.SaveMetaFactory](c)
		savesService := ioc.Get[saves.Saves](c)
		repo := ioc.Get[IntRepo](c)
		fmt.Printf("count s1 is %d", repo.GetCount())
		s1 := saveMetaFactory.New(saves.SaveName("s1"))
		if err := savesService.NewSave(s1); err != nil {
			panic(err)
		}
		repo.Increment()
		fmt.Printf("count s2 is %d", repo.GetCount())
		s2 := saveMetaFactory.New(saves.SaveName("s2"))
		if err := savesService.NewSave(s2); err != nil {
			panic(err)
		}
		fmt.Print("loading")
		if err := savesService.Load(s1.Id); err != nil {
			panic(err)
		}
		fmt.Printf("count s1 is %d", repo.GetCount())
		if err := savesService.Load(s2.Id); err != nil {
			panic(err)
		}
		fmt.Printf("count s2 is %d", repo.GetCount())
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
