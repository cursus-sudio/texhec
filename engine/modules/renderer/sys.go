package renderer

import (
	"engine/services/ecs"
)

type System ecs.SystemRegister

type Service interface {
	// adds default render component
	Render(ecs.EntityID)

	Direct() ecs.ComponentsArray[DirectComponent]
}
