package scenes

import "github.com/ogiusek/events"

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

type sceneManager struct {
	scenes         map[SceneId]Scene
	events         events.Events
	currentSceneId SceneId
}

func (sceneManager *sceneManager) CurrentScene() SceneId {
	return sceneManager.currentSceneId
}

func (sceneManager *sceneManager) CurrentSceneEvents() events.Events {
	return sceneManager.events
}

func (sceneManager *sceneManager) LoadScene(sceneId SceneId) error {
	loadedScene, ok := sceneManager.scenes[sceneId]
	if !ok {
		return ErrSceneDoNotExists
	}
	sceneManager.scenes[sceneManager.currentSceneId].Unload()

	sceneManager.currentSceneId = sceneId
	e := loadedScene.Load(sceneManager)
	sceneManager.events = e
	return nil
}
