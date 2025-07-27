package example

import (
	"backend/services/saves"
	"core/tacticalmap"
	"encoding/json"
	"errors"
	"fmt"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/frames"
	"frontend/services/scenes"
	"shared/utils/connection"
	"sync"
	"time"

	"github.com/ogiusek/relay/v2"
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
	SceneManager   scenes.SceneManager
	World          ecs.World
	Toggled        bool
	Scene2         scenes.SceneId
	ToggleTreshold time.Duration
}

func NewToggledSystem(
	sceneManager scenes.SceneManager,
	world ecs.World,
	scene2 scenes.SceneId,
	toggleTreshold time.Duration,
) toggleSystem {
	return toggleSystem{
		SceneManager:   sceneManager,
		World:          world,
		Toggled:        false,
		Scene2:         scene2,
		ToggleTreshold: toggleTreshold,
	}
}

func (system *toggleSystem) Update(args frames.FrameEvent) {
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
		if component.PastTime < system.ToggleTreshold {
			return
		}
	}

	system.SceneManager.LoadScene(system.Scene2)
	system.Toggled = true
}

//

//

type someSystem struct {
	SceneManager scenes.SceneManager
	World        ecs.World
	Backend      connection.Connection
	Console      console.Console
}

func NewSomeSystem(
	sceneMagener scenes.SceneManager,
	world ecs.World,
	backend connection.Connection,
	console console.Console,
) someSystem {
	return someSystem{
		SceneManager: sceneMagener,
		World:        world,
		Backend:      backend,
		Console:      console,
	}
}

func (system *someSystem) Update(args frames.FrameEvent) {
	format := "02-01-2006 15:04:05"
	text := "----------------------------------------------------------------\n"
	text += fmt.Sprintf("now %s\n", time.Now().Format(format))

	for _, entity := range system.World.GetEntitiesWithComponents(ecs.GetComponentType(someComponent{})) {
		component := someComponent{}
		if err := system.World.GetComponent(entity, &component); err != nil {
			continue
		}
		component.PastTime += args.Delta
		component.Frame += 1
		system.World.SaveComponent(entity, component)

		text += "\n"
		text += fmt.Sprintf("first frame %s\n", component.StartedAt.Format(format))
		text += fmt.Sprintf("current frame %d\n", component.Frame)
		// text += fmt.Sprintf("time in game %d\n", int(component.PastTime.Seconds()))
		text += fmt.Sprintf("time in game %f\n", component.PastTime.Seconds())
	}
	text += "\n"

	if system.Backend == nil {
		text += fmt.Sprintf("backend is nil\n")
	} else {
		tacticalMap, err := relay.Handle(system.Backend.Relay(), tacticalmap.NewGetReq())
		if err == nil {
			text += fmt.Sprintf("found shit %v\n", tacticalMap)
		} else {
			text += fmt.Sprintf("found error %s\n", err)
		}
	}

	text += "entities found: {\n"
	for _, entity := range system.World.GetEntities() {
		text += fmt.Sprintf(" - %s\n", entity)
	}
	text += "}\n"
	text += "\n"

	system.Console.Print(text)
	system.Console.Flush()
}

//

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
