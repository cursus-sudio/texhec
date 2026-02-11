package internal

import (
	"engine/modules/scene"
	"engine/services/ecs"
	"engine/services/logger"
	"fmt"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

// type Service interface {
// 	SetScene(id ID, loader func(sceneParent ecs.EntityID))
// }

type SceneComp struct{}

type Service struct {
	scenes        map[scene.ID]scene.Scene
	SceneArr      ecs.ComponentsArray[SceneComp]
	Logger        logger.Logger  `inject:"1"`
	World         ecs.World      `inject:"1"`
	EventsBuilder events.Builder `inject:"1"`
}

func NewService(c ioc.Dic) scene.Service {
	service := ioc.GetServices[*Service](c)
	service.scenes = make(map[scene.ID]scene.Scene)
	service.SceneArr = ecs.GetComponentsArray[SceneComp](service.World)
	entity := service.World.NewEntity()
	service.SceneArr.Set(entity, SceneComp{})

	events.Listen(service.EventsBuilder, service.ChangeScene)
	return service
}

func (service *Service) ChangeScene(event scene.ChangeSceneEvent) {
	for _, entity := range service.SceneArr.GetEntities() {
		service.World.RemoveEntity(entity)
	}
	sceneEntity := service.World.NewEntity()
	service.SceneArr.Set(sceneEntity, SceneComp{})

	scene, ok := service.scenes[event.ID]
	if !ok {
		service.Logger.Warn(fmt.Errorf("scene with id %v doesn't exist", event.ID))
		return
	}
	scene(sceneEntity)
}

func (service *Service) SetScene(id scene.ID, scene scene.Scene) {
	service.scenes[id] = scene
}
