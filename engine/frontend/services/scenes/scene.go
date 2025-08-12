package scenes

import (
	"frontend/services/ecs"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

//

const (
	_ ioc.Order = iota

	LoadFirst

	// load world
	LoadWorld

	// load event listeners
	LoadBeforeDomain
	_
	_
	_

	LoadDomain
	_
	_
	_

	LoadAfterDomain
	_
	_
	_

	// call events affecting world state
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

type sceneCtx struct {
	World         ecs.World
	EventsBuilder events.Builder
	Events        events.Events
}

type SceneCtx struct{ *sceneCtx }

func NewSceneCtx(world ecs.World, eventsBuilder events.Builder) SceneCtx {
	return SceneCtx{&sceneCtx{
		world,
		eventsBuilder,
		eventsBuilder.Events(),
	}}
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
	var ctxPtr *SceneCtx
	onUnload := func() {
		if ctxPtr == nil {
			return
		}
		ctx := *ctxPtr
		for _, listener := range n.onUnload {
			listener(ctx)
		}
		ctxPtr = nil
	}
	onLoad := func() SceneCtx {
		if ctxPtr != nil {
			return *ctxPtr
		}

		ctx := NewSceneCtx(
			ecs.NewWorld(),
			events.NewBuilder(),
		)
		for _, listener := range n.onLoad {
			listener(ctx)
		}
		ctxPtr = &ctx
		return ctx
	}
	return newScene(sceneId, onLoad, onUnload)
}

//

type Scene interface {
	Id() SceneId

	Load() SceneCtx
	Unload()
}

type scene struct {
	id       SceneId
	onLoad   func() SceneCtx
	onUnload func()
}

func newScene(
	id SceneId,
	onLoad func() SceneCtx,
	onUnload func(),
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

func (scene *scene) Load() SceneCtx { return scene.onLoad() }
func (scene *scene) Unload()        { scene.onUnload() }
