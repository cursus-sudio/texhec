package scenes

type SceneManagerBuilder interface {
	AddScene(Scene) SceneManagerBuilder
	MakeActive(SceneId) SceneManagerBuilder
	OnSceneLoad(func(SceneCtx)) SceneManagerBuilder
	OnSceneUnload(func(SceneCtx)) SceneManagerBuilder
	Build() SceneManager
}

type sceneManagerBuilder struct {
	scenes          map[SceneId]Scene
	currentScene    *SceneId
	loadListeners   []func(SceneCtx)
	unloadListeners []func(SceneCtx)
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
func (b *sceneManagerBuilder) OnSceneLoad(listener func(SceneCtx)) SceneManagerBuilder {
	b.loadListeners = append(b.loadListeners, listener)
	return b
}

// listener is triggered after unloading scene
func (b *sceneManagerBuilder) OnSceneUnload(listener func(SceneCtx)) SceneManagerBuilder {
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
		activeSceneCtx:     nil,
		loadedCurrentScene: false,
		currentSceneId:     id,

		onLoad: func(sc SceneCtx) {
			for _, listener := range b.loadListeners {
				listener(sc)
			}
		},
		onUnload: func(sc SceneCtx) {
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
	CurrentSceneCtx() SceneCtx

	// this method returns error:
	// - ErrSceneDoNotExists
	LoadScene(SceneId) error
}

type sceneManager struct {
	scenes             map[SceneId]Scene
	activeSceneCtx     SceneCtx
	loadedCurrentScene bool
	currentSceneId     SceneId

	onLoad   func(SceneCtx)
	onUnload func(SceneCtx)
}

func (sceneManager *sceneManager) CurrentScene() SceneId {
	return sceneManager.currentSceneId
}

func (sceneManager *sceneManager) CurrentSceneCtx() SceneCtx {
	if !sceneManager.loadedCurrentScene {
		sceneManager.LoadScene(sceneManager.currentSceneId)
	}
	return sceneManager.activeSceneCtx
}

func (sceneManager *sceneManager) LoadScene(sceneId SceneId) error {
	loadedScene, ok := sceneManager.scenes[sceneId]
	if !ok {
		return ErrSceneDoNotExists
	}

	var unloadedScene Scene
	previousCtx := sceneManager.activeSceneCtx
	if sceneManager.loadedCurrentScene {
		unloadedScene = sceneManager.scenes[sceneManager.currentSceneId]
	}

	sceneManager.loadedCurrentScene = true

	// load
	sceneManager.currentSceneId = sceneId
	ctx := loadedScene.Load()
	sceneManager.onLoad(ctx)
	sceneManager.activeSceneCtx = ctx

	// unload
	if unloadedScene != nil {
		unloadedScene.Unload(previousCtx)
		sceneManager.onUnload(previousCtx)
	}

	return nil
}
