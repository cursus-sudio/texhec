package scenes

import "github.com/ogiusek/events"

type SceneManagerBuilder interface {
	AddScene(Scene) SceneManagerBuilder
	MakeActive(SceneId) SceneManagerBuilder
	Build() SceneManager
}

type sceneManagerBuilder struct {
	scenes       map[SceneId]Scene
	currentScene *SceneId
}

func NewSceneManagerBuilder() SceneManagerBuilder {
	return &sceneManagerBuilder{
		scenes:       map[SceneId]Scene{},
		currentScene: nil,
	}
}

func (b *sceneManagerBuilder) AddScene(scene Scene) SceneManagerBuilder {
	id := scene.Id()
	if _, ok := b.scenes[id]; ok {
		panic(ErrSceneAlreadyExists.Error())
	}
	b.scenes[id] = scene
	return b
}

func (b *sceneManagerBuilder) MakeActive(sceneId SceneId) SceneManagerBuilder {
	sceneIdHeap := sceneId
	b.currentScene = &sceneIdHeap
	return b
}

func (b *sceneManagerBuilder) Build() SceneManager {
	if b.currentScene == nil {
		panic(ErrNoActiveScene.Error())
	}
	id := *b.currentScene
	scene, ok := b.scenes[id]
	if !ok {
		panic(ErrSceneDoNotExists.Error())
	}

	manager := &sceneManager{
		scenes:         b.scenes,
		currentSceneId: id,
	}
	manager.events = scene.Load(manager)
	return manager
}

//

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

type sceneManager struct {
	scenes             map[SceneId]Scene
	events             events.Events
	loadedCurrentScene bool
	currentSceneId     SceneId
}

func (sceneManager *sceneManager) CurrentScene() SceneId {
	return sceneManager.currentSceneId
}

func (sceneManager *sceneManager) CurrentSceneEvents() events.Events {
	if !sceneManager.loadedCurrentScene {
		sceneManager.LoadScene(sceneManager.currentSceneId)
	}
	return sceneManager.events
}

func (sceneManager *sceneManager) LoadScene(sceneId SceneId) error {
	loadedScene, ok := sceneManager.scenes[sceneId]
	if !ok {
		return ErrSceneDoNotExists
	}
	if sceneManager.loadedCurrentScene {
		sceneManager.scenes[sceneManager.currentSceneId].Unload()
	}

	sceneManager.currentSceneId = sceneId
	e := loadedScene.Load(sceneManager)
	sceneManager.loadedCurrentScene = true
	sceneManager.events = e
	return nil
}
