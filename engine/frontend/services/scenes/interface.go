package scenes

import (
	"errors"

	"github.com/ogiusek/events"
)

// scene

type SceneId struct {
	sceneId string
}

func NewSceneId(sceneId string) SceneId {
	return SceneId{sceneId: sceneId}
}

type Scene interface {
	Id() SceneId

	Load(SceneManager) events.Events
	Unload()
}

// scene builder

type SceneBuilder interface {
	OnLoad(func(SceneManager, Scene)) SceneBuilder
	OnUnload(func(Scene)) SceneBuilder
	Build(SceneId, func(SceneManager, Scene) events.Events) Scene
}

// scene manager

var (
	ErrSceneAlreadyExists error = errors.New("scene already exists")
	ErrNoActiveScene      error = errors.New("no active scene")
	ErrSceneDoNotExists   error = errors.New("scene do not exists")
)

type SceneManagerBuilder interface {
	AddScene(Scene) SceneManagerBuilder
	MakeActive(SceneId) SceneManagerBuilder
	Build() SceneManager
}

// TODO
// type SceneManagerBuilder interface {
// 	AddScene(Scene) SceneManagerBuilder
// }

type SceneManager interface {
	// if scene is first then it automatically is loaded
	// this method returns error:
	// - ErrSceneAlreadyExists
	// AddScene(Scene) error

	// can panic when no scene is loaded
	CurrentScene() SceneId
	CurrentSceneEvents() events.Events

	// this method returns error:
	// - ErrSceneDoNotExists
	LoadScene(SceneId) error
}
