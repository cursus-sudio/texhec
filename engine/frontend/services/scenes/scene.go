package scenes

import "frontend/services/ecs"

type scene struct {
	id       SceneId
	world    ecs.World
	onLoad   func(Scene)
	onUnload func(Scene)
	events   SceneEvents
}

func newScene(
	id SceneId,
	world ecs.World,
	onLoad func(Scene),
	onUnload func(Scene),
	events SceneEvents,
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
func (scene *scene) SceneEvents() SceneEvents { return scene.events }
