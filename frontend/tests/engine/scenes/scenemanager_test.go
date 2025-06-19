package scenes_test

import (
	"frontend/src/engine/ecs"
	"frontend/src/engine/scenes"
	"testing"

	"github.com/ogiusek/ioc/v2"
)

// type SceneId struct {
// 	sceneId string
// }
//
// // scene manager
// 	ErrSceneAlreadyExists error = errors.New("scene already exists")
// 	ErrSceneDoNotExists   error = errors.New("scene do not exists")
//
// type SceneManager interface {
// 	// - ErrSceneAlreadyExists
// 	AddScene(Scene) error

// 	CurrentScene() Scene

// 	GetScene(SceneId) null.Nullable[Scene]

// 	GetScenes() []Scene

// 	// - ErrSceneDoNotExists
// 	LoadScene(SceneId) error
// }

func TestScenes(t *testing.T) {
	b := ioc.NewBuilder()
	scenes.Package().Register(b)
	ecs.Package().Register(b)
	c := b.Build()
	worldFactory := ioc.Get[ecs.WorldFactory](c)
	sceneManager := ioc.Get[scenes.SceneManager](c)
	sceneBuilder := ioc.Get[scenes.SceneBuilder](c)

	loaded := false
	unloaded := false
	world := worldFactory()
	scene := sceneBuilder.
		OnLoad(func(s scenes.Scene) { loaded = true }).
		OnUnload(func(s scenes.Scene) { unloaded = true }).
		Build(scenes.NewSceneId("main scene"), world)

	loadedOther := false
	otherScene := sceneBuilder.
		OnLoad(func(s scenes.Scene) { loadedOther = true }).
		OnUnload(func(s scenes.Scene) {}).
		Build(scenes.NewSceneId("other scene"), world)

	if loaded {
		t.Errorf("scene was loaded on creation")
	}

	if err := sceneManager.AddScene(scene); err != nil {
		t.Errorf("unexpected error when adding scene")
	}

	if err := sceneManager.AddScene(scene); err != scenes.ErrSceneAlreadyExists {
		t.Errorf("expected error when adding second time the same scene")
	}

	if sceneManager.CurrentScene().Id() != scene.Id() {
		t.Errorf("currect scene doesn't match first loaded scene")
	}

	if !loaded {
		t.Errorf("scene didn't got loaded by default")
	}

	if err := sceneManager.LoadScene(otherScene.Id()); err != scenes.ErrSceneDoNotExists {
		t.Errorf("succesfully loaded not added scene")
	}

	if err := sceneManager.AddScene(otherScene); err != nil {
		t.Errorf("encnountered error when adding other scene")
	}

	if sceneManager.CurrentScene().Id() != scene.Id() {
		t.Errorf("current scene isn't equal to latest loaded scene")
	}

	if unloaded {
		t.Errorf("unloaded scene without switching scenes")
	}

	if loadedOther {
		t.Errorf("prematurely loaded other scene")
	}

	if err := sceneManager.LoadScene(otherScene.Id()); err != nil {
		t.Errorf("encountered error when loading other scene")
	}

	if !unloaded {
		t.Errorf("didn't unloaded scene when loaded other")
	}

	if !loadedOther {
		t.Errorf("didn't trigger scene on load")
	}

	if sceneManager.CurrentScene().Id() != otherScene.Id() {
		t.Errorf("current scene id isn't equal to curent scene ")
	}
}
