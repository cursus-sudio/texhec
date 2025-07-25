package scenes

import "github.com/ogiusek/events"

//

type SceneId struct {
	sceneId string
}

func NewSceneId(sceneId string) SceneId {
	return SceneId{sceneId: sceneId}
}

//

type SceneBuilder interface {
	// load func(SceneManager, Scene) events.Events
	// Load(func(SceneManager, ))
	OnLoad(func(SceneManager, Scene, events.Builder))
	OnUnload(func(Scene))
	Build(id SceneId) Scene
}

type sceneBuilder struct {
	onLoad   []func(SceneManager, Scene, events.Builder)
	onUnload []func(Scene)
}

func NewSceneBuilder() SceneBuilder {
	return &sceneBuilder{}
}

func (b *sceneBuilder) OnLoad(listener func(SceneManager, Scene, events.Builder)) {
	b.onLoad = append(b.onLoad, listener)
}

func (b *sceneBuilder) OnUnload(listener func(Scene)) {
	b.onUnload = append(b.onUnload, listener)
}

func (n *sceneBuilder) Build(sceneId SceneId) Scene {
	onLoad := func(manager SceneManager, s Scene) events.Events {
		b := events.NewBuilder()
		for _, listener := range n.onLoad {
			listener(manager, s, b)
		}
		return b.Build()
	}
	onUnload := func(s Scene) {
		for _, listener := range n.onUnload {
			listener(s)
		}
	}
	return newScene(sceneId, onLoad, onUnload)
}

//

type Scene interface {
	Id() SceneId

	Load(SceneManager) events.Events
	Unload()
}

type scene struct {
	id       SceneId
	onLoad   func(SceneManager, Scene) events.Events
	onUnload func(Scene)
}

func newScene(
	id SceneId,
	onLoad func(SceneManager, Scene) events.Events,
	onUnload func(Scene),
) Scene {
	return &scene{
		id:       id,
		onLoad:   onLoad,
		onUnload: onUnload,
	}
}

func (scene *scene) Id() SceneId {
	return scene.id
}

func (scene *scene) Load(manager SceneManager) events.Events { return scene.onLoad(manager, scene) }
func (scene *scene) Unload()                                 { scene.onUnload(scene) }
