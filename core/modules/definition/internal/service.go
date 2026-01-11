package internal

import (
	"core/modules/definition"
	"engine/services/ecs"
)

type service struct {
	definitionArray     ecs.ComponentsArray[definition.DefinitionComponent]
	definitionLinkArray ecs.ComponentsArray[definition.DefinitionLinkComponent]
}

func NewService(
	w ecs.World,
) definition.Service {
	t := &service{
		definitionArray:     ecs.GetComponentsArray[definition.DefinitionComponent](w),
		definitionLinkArray: ecs.GetComponentsArray[definition.DefinitionLinkComponent](w),
	}

	return t
}

func (t *service) Component() ecs.ComponentsArray[definition.DefinitionComponent] {
	return t.definitionArray
}
func (t *service) Link() ecs.ComponentsArray[definition.DefinitionLinkComponent] {
	return t.definitionLinkArray
}
