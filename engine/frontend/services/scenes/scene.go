package scenes

import (
	"frontend/services/ecs"

	"github.com/ogiusek/events"
)

type scene struct {
	id       SceneId
	world    ecs.World
	onLoad   func(Scene)
	onUnload func(Scene)
	events   events.Events
}

func newScene(
	id SceneId,
	world ecs.World,
	onLoad func(Scene),
	onUnload func(Scene),
	events events.Events,
) Scene {
	return &scene{
		id:       id,
		world:    world,
		onLoad:   onLoad,
		onUnload: onUnload,
		events:   events,
	}
}

func (scene *scene) Id() SceneId {
	return scene.id
}

func (scene *scene) Load()   { scene.onLoad(scene) }
func (scene *scene) Unload() { scene.onUnload(scene) }

func (scene *scene) World() ecs.World {
	return scene.world
}
func (scene *scene) Events() events.Events { return scene.events }
