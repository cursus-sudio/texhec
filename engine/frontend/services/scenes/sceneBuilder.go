package scenes

import "github.com/ogiusek/events"

type sceneBuilder struct {
	load     func(SceneManager, Scene) events.Events
	onLoad   []func(SceneManager, Scene)
	onUnload []func(Scene)
}

func newSceneBuilder() SceneBuilder {
	return sceneBuilder{}
}

func (b sceneBuilder) OnLoad(listener func(SceneManager, Scene)) SceneBuilder {
	b.onLoad = append(b.onLoad, listener)
	return b
}

func (b sceneBuilder) OnUnload(listener func(Scene)) SceneBuilder {
	b.onUnload = append(b.onUnload, listener)
	return b
}

func (n sceneBuilder) Build(sceneId SceneId, load func(SceneManager, Scene) events.Events) Scene {
	onLoad := func(manager SceneManager, s Scene) events.Events {
		e := load(manager, s)
		for _, listener := range n.onLoad {
			listener(manager, s)
		}
		return e
	}
	onUnload := func(s Scene) {
		for _, listener := range n.onUnload {
			listener(s)
		}
	}
	return newScene(sceneId, onLoad, onUnload)
}
