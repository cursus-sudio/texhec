package main

import (
	"backend/src/backendapi"
	"backend/src/backendapi/tacticalmapapi"
	"backend/src/modules/saves"
	"encoding/json"
	"errors"
	"fmt"
	"frontend/src/engine/console"
	"frontend/src/engine/ecs"
	"frontend/src/engine/scenes"
	"reflect"
	"sync"
	"time"

	"github.com/ogiusek/ioc"
	"github.com/ogiusek/relay"
)

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

	tacticalMap, _ := relay.Handle(system.Backend.Relay(), tacticalmapapi.NewGetReq())
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
	return err == nil
}

func (savableIntRepo *intRepo) TakeSnapshot() saves.RepoSnapshot {
	bytes, err := json.Marshal(savableIntRepo)
	if err != nil {
		panic(errors.Join(errors.New("error taking snapshot"), err))
	}
	savableIntRepo.Changed = false
	return saves.NewRepoSnapshot(bytes)
}

func (savableIntRepo *intRepo) LoadSnapshot(snapshot saves.RepoSnapshot) error {
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

type exBackendModPkg struct{}

func (exBackendModPkg) Register(c ioc.Dic) {
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

		builder := ioc.Get[saves.ListSavesQueryBuilder](c)
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
	}
}
