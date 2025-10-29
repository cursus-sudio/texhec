package scenes

import (
	"shared/services/ecs"

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

type SceneCtx ecs.World

func NewSceneCtx(world ecs.World) SceneCtx {
	return world
}

//

type SceneBuilder interface {
	OnLoad(func(ctx SceneCtx)) SceneBuilder
	OnUnload(func(ctx SceneCtx)) SceneBuilder
	Build(id SceneId) Scene
}

type sceneBuilder struct {
	onLoad   []func(SceneCtx)
	onUnload []func(SceneCtx)
}

func NewSceneBuilder() SceneBuilder {
	return &sceneBuilder{}
}

func (b *sceneBuilder) OnLoad(listener func(SceneCtx)) SceneBuilder {
	b.onLoad = append(b.onLoad, listener)
	return b
}

func (b *sceneBuilder) OnUnload(listener func(SceneCtx)) SceneBuilder {
	b.onUnload = append(b.onUnload, listener)
	return b
}

func (n *sceneBuilder) Build(sceneId SceneId) Scene {
	onUnload := func(ctx SceneCtx) {
		for _, listener := range n.onUnload {
			listener(ctx)
		}
		ctx.Release()
	}
	onLoad := func() SceneCtx {
		ctx := NewSceneCtx(ecs.NewWorld())
		for _, listener := range n.onLoad {
			listener(ctx)
		}
		return ctx
	}
	return newScene(sceneId, onLoad, onUnload)
}

//

type Scene interface {
	Id() SceneId

	Load() SceneCtx
	Unload(ctx SceneCtx)
}

type scene struct {
	id       SceneId
	onLoad   func() SceneCtx
	onUnload func(SceneCtx)
}

func newScene(
	id SceneId,
	onLoad func() SceneCtx,
	onUnload func(SceneCtx),
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

func (scene *scene) Load() SceneCtx      { return scene.onLoad() }
func (scene *scene) Unload(ctx SceneCtx) { scene.onUnload(ctx) }
