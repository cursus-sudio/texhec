package internal

import (
	"core/modules/definition"
	"engine"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	World               engine.World `inject:"1"`
	definitionArray     ecs.ComponentsArray[definition.DefinitionComponent]
	definitionLinkArray ecs.ComponentsArray[definition.DefinitionLinkComponent]
}

func NewService(c ioc.Dic) definition.Service {
	t := ioc.GetServices[*service](c)
	t.definitionArray = ecs.GetComponentsArray[definition.DefinitionComponent](t.World)
	t.definitionLinkArray = ecs.GetComponentsArray[definition.DefinitionLinkComponent](t.World)

	return t
}

func (t *service) Component() ecs.ComponentsArray[definition.DefinitionComponent] {
	return t.definitionArray
}
func (t *service) Link() ecs.ComponentsArray[definition.DefinitionLinkComponent] {
	return t.definitionLinkArray
}
