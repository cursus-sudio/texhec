package scenes

import "github.com/ogiusek/events"

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
