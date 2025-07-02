package scenes

import (
	"frontend/services/ecs"

	"github.com/ogiusek/events"
)

type sceneBuilder struct {
	onLoad        []func(Scene)
	onUnload      []func(Scene)
	eventsBuilder events.Builder
}

func newSceneBuilder() SceneBuilder {
	return sceneBuilder{}
}

func (b sceneBuilder) OnLoad(listener func(Scene)) SceneBuilder {
	b.onLoad = append(b.onLoad, listener)
	return b
}

func (b sceneBuilder) OnUnload(listener func(Scene)) SceneBuilder {
	b.onUnload = append(b.onUnload, listener)
	return b
}

func (b sceneBuilder) Events(event func(events.Builder)) SceneBuilder {
	event(b.eventsBuilder)
	return b
}

func (n sceneBuilder) Build(sceneId SceneId, world ecs.World) Scene {
	onLoad := func(s Scene) {
		for _, listener := range n.onLoad {
			listener(s)
		}
	}
	onUnload := func(s Scene) {
		for _, listener := range n.onUnload {
			listener(s)
		}
	}
	return newScene(sceneId, world, onLoad, onUnload, n.eventsBuilder.Build())
}
