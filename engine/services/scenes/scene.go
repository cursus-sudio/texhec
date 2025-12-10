package scenes

import (
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

//

const (
	_ ioc.Order = iota
	LoadConfig
	LoadObjects
	LoadSystems
	LoadInitialEvents
)

//

type SceneId struct {
	sceneId string
}

func NewSceneId(sceneId string) SceneId {
	return SceneId{sceneId: sceneId}
}

//

//

type SceneBuilder interface {
	OnLoad(func(world ecs.World)) SceneBuilder
	OnUnload(func(world ecs.World)) SceneBuilder
	Build(id SceneId) Scene
}

type sceneBuilder struct {
	onLoad   []func(ecs.World)
	onUnload []func(ecs.World)
}

func NewSceneBuilder() SceneBuilder {
	return &sceneBuilder{}
}

func (b *sceneBuilder) OnLoad(listener func(ecs.World)) SceneBuilder {
	b.onLoad = append(b.onLoad, listener)
	return b
}

func (b *sceneBuilder) OnUnload(listener func(ecs.World)) SceneBuilder {
	b.onUnload = append(b.onUnload, listener)
	return b
}

func (n *sceneBuilder) Build(sceneId SceneId) Scene {
	onUnload := func(world ecs.World) {
		for _, listener := range n.onUnload {
			listener(world)
		}
		world.Release()
	}
	onLoad := func() ecs.World {
		world := ecs.NewWorld()
		for _, listener := range n.onLoad {
			listener(world)
		}
		return world
	}
	return newScene(sceneId, onLoad, onUnload)
}

//

type Scene interface {
	Id() SceneId

	Load() ecs.World
	Unload(ecs.World)
}

type scene struct {
	id       SceneId
	onLoad   func() ecs.World
	onUnload func(ecs.World)
}

func newScene(
	id SceneId,
	onLoad func() ecs.World,
	onUnload func(ecs.World),
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

func (scene *scene) Load() ecs.World        { return scene.onLoad() }
func (scene *scene) Unload(world ecs.World) { scene.onUnload(world) }
