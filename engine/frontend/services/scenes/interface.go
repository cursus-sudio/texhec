package scenes

import (
	"errors"
	"frontend/services/ecs"

	"github.com/ogiusek/events"
	"github.com/ogiusek/null"
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

	Load()
	Unload()

	Events() events.Events
	World() ecs.World
}

// scene builder

type SceneBuilder interface {
	OnLoad(func(Scene)) SceneBuilder
	OnUnload(func(Scene)) SceneBuilder
	Events(func(events.Builder)) SceneBuilder
	Build(SceneId, ecs.World) Scene
}

// scene manager

var (
	ErrSceneAlreadyExists error = errors.New("scene already exists")
	ErrSceneDoNotExists   error = errors.New("scene do not exists")
)

type SceneManager interface {
	// if scene is first then it automatically is loaded
	// this method returns error:
	// - ErrSceneAlreadyExists
	AddScene(Scene) error

	// can panic when no scene is loaded
	CurrentScene() Scene

	GetScene(SceneId) null.Nullable[Scene]
	GetScenes() []Scene

	// this method returns error:
	// - ErrSceneDoNotExists
	LoadScene(SceneId) error
}
