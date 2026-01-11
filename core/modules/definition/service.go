package definition

import (
	"engine/services/ecs"
)

type Service interface {
	Component() ecs.ComponentsArray[DefinitionComponent]
	Link() ecs.ComponentsArray[DefinitionLinkComponent]
}
