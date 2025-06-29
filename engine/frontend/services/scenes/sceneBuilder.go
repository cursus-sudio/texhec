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

func (sceneBuilder sceneBuilder) OnLoad(listener func(Scene)) SceneBuilder {
	sceneBuilder.onLoad = append(sceneBuilder.onLoad, listener)
	return sceneBuilder
}

func (sceneBuilder sceneBuilder) OnUnload(listener func(Scene)) SceneBuilder {
	sceneBuilder.onUnload = append(sceneBuilder.onUnload, listener)
	return sceneBuilder
}

func (sceneBuilder sceneBuilder) Build(sceneId SceneId, world ecs.World) Scene {
	onLoad := func(s Scene) {
		for _, listener := range sceneBuilder.onLoad {
			listener(s)
		}
	}
	onUnload := func(s Scene) {
		for _, listener := range sceneBuilder.onUnload {
			listener(s)
		}
	}
	return newScene(sceneId, world, onLoad, onUnload, NewSceneEvents(sceneBuilder.eventsBuilder.Build()))
}
