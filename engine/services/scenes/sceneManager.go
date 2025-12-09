package scenes

import "engine/services/ecs"

type SceneManagerBuilder interface {
	AddScene(Scene) SceneManagerBuilder
	MakeActive(SceneId) SceneManagerBuilder
	OnSceneLoad(func(ecs.World)) SceneManagerBuilder
	OnSceneUnload(func(ecs.World)) SceneManagerBuilder
	Build() SceneManager
}

type sceneManagerBuilder struct {
	scenes          map[SceneId]Scene
	currentScene    *SceneId
	loadListeners   []func(ecs.World)
	unloadListeners []func(ecs.World)
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

// listener is triggered after loading scene
func (b *sceneManagerBuilder) OnSceneLoad(listener func(ecs.World)) SceneManagerBuilder {
	b.loadListeners = append(b.loadListeners, listener)
	return b
}

// listener is triggered after unloading scene
func (b *sceneManagerBuilder) OnSceneUnload(listener func(ecs.World)) SceneManagerBuilder {
	b.unloadListeners = append(b.unloadListeners, listener)
	return b
}

func (b *sceneManagerBuilder) Build() SceneManager {
	if b.currentScene == nil {
		panic(ErrNoActiveScene.Error())
	}
	id := *b.currentScene
	_, ok := b.scenes[id]
	if !ok {
		panic(ErrSceneDoNotExists.Error())
	}

	manager := &sceneManager{
		scenes:             b.scenes,
		activeWorld:        nil,
		loadedCurrentScene: false,
		currentSceneId:     id,

		onLoad: func(sc ecs.World) {
			for _, listener := range b.loadListeners {
				listener(sc)
			}
		},
		onUnload: func(sc ecs.World) {
			for _, listener := range b.unloadListeners {
				listener(sc)
			}
		},
	}
	return manager
}

//

type SceneManager interface {
	// if scene is first then it automatically is loaded
	// this method returns error:
	// - ErrSceneAlreadyExists
	// AddScene(Scene) error

	CurrentScene() SceneId
	CurrentSceneWorld() ecs.World

	// this method returns error:
	// - ErrSceneDoNotExists
	LoadScene(SceneId) error
}

type sceneManager struct {
	scenes             map[SceneId]Scene
	activeWorld        ecs.World
	loadedCurrentScene bool
	currentSceneId     SceneId

	onLoad   func(ecs.World)
	onUnload func(ecs.World)
}

func (sceneManager *sceneManager) CurrentScene() SceneId {
	return sceneManager.currentSceneId
}

func (sceneManager *sceneManager) CurrentSceneWorld() ecs.World {
	if !sceneManager.loadedCurrentScene {
		sceneManager.LoadScene(sceneManager.currentSceneId)
	}
	return sceneManager.activeWorld
}

func (sceneManager *sceneManager) LoadScene(sceneId SceneId) error {
	loadedScene, ok := sceneManager.scenes[sceneId]
	if !ok {
		return ErrSceneDoNotExists
	}

	var unloadedScene Scene
	previousCtx := sceneManager.activeWorld
	if sceneManager.loadedCurrentScene {
		unloadedScene = sceneManager.scenes[sceneManager.currentSceneId]
	}

	sceneManager.loadedCurrentScene = true

	// load
	sceneManager.currentSceneId = sceneId
	world := loadedScene.Load()
	sceneManager.onLoad(world)
	sceneManager.activeWorld = world

	// unload
	if unloadedScene != nil {
		unloadedScene.Unload(previousCtx)
		sceneManager.onUnload(previousCtx)
	}

	return nil
}
