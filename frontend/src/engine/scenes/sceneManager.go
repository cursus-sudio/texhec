package scenes

import (
	"github.com/ogiusek/null"
)

type sceneManager struct {
	scenes         map[SceneId]Scene
	currentSceneId SceneId
}

func newSceneManager() SceneManager {
	return &sceneManager{
		scenes:         map[SceneId]Scene{},
		currentSceneId: NewSceneId(""),
	}
}

func (sceneManager *sceneManager) AddScene(scene Scene) error {
	sceneId := scene.Id()
	if _, ok := sceneManager.scenes[sceneId]; ok {
		return ErrSceneAlreadyExists
	}
	sceneManager.scenes[sceneId] = scene
	if sceneManager.currentSceneId == NewSceneId("") {
		sceneManager.LoadScene(scene.Id())
	}
	return nil
}

func (sceneManager *sceneManager) CurrentScene() Scene {
	scene, ok := sceneManager.scenes[sceneManager.currentSceneId]
	if !ok {
		panic("no scene was loaded")
	}
	return scene
}

func (sceneManager *sceneManager) GetScene(sceneId SceneId) null.Nullable[Scene] {
	scene, ok := sceneManager.scenes[sceneId]
	if !ok {
		return null.Null[Scene]()
	}
	return null.New(scene)
}

func (sceneManager *sceneManager) GetScenes() []Scene {
	scenes := make([]Scene, len(sceneManager.scenes))
	for _, scene := range sceneManager.scenes {
		scenes = append(scenes, scene)
	}
	return scenes
}

func (sceneManager *sceneManager) LoadScene(sceneId SceneId) error {
	newScene, ok := sceneManager.scenes[sceneId]
	if !ok {
		return ErrSceneDoNotExists
	}
	currentScene, ok := sceneManager.scenes[sceneManager.currentSceneId]
	if ok {
		currentScene.Unload()
	}

	sceneManager.currentSceneId = sceneId
	newScene.Load()
	return nil
}
