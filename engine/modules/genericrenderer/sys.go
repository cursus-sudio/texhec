package genericrenderer

import (
	"engine/services/ecs"
)

type System ecs.SystemRegister

type Service interface {
	Direct() ecs.ComponentsArray[DirectComponent]
}
