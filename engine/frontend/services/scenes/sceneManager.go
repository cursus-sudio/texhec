package scenes

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
	_, ok := b.scenes[id]
	if !ok {
		panic(ErrSceneDoNotExists.Error())
	}

	manager := &sceneManager{
		scenes:             b.scenes,
		currentSceneId:     id,
		loadedCurrentScene: false,
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

	// load
	sceneManager.currentSceneId = sceneId
	ctx := loadedScene.Load()
	sceneManager.loadedCurrentScene = true
	sceneManager.activeSceneCtx = ctx

	// unload
	if sceneManager.loadedCurrentScene {
		sceneManager.scenes[sceneManager.currentSceneId].Unload()
	}

	return nil
}
