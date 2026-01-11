package engine

import (
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/connection"
	"engine/modules/genericrenderer"
	"engine/modules/groups"
	"engine/modules/hierarchy"
	"engine/modules/inputs"
	"engine/modules/layout"
	"engine/modules/netsync"
	"engine/modules/record"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/modules/transition"
	"engine/modules/uuid"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/events"
)

type World struct {
	Logger        logger.Logger  `inject:"1"`
	EventsBuilder events.Builder `inject:"1"`
	Events        events.Events  `inject:"1"`

	ecs.World       `inject:"1"`
	Camera          camera.Service          `inject:"1"`
	Collider        collider.Service        `inject:"1"`
	Connection      connection.Service      `inject:"1"`
	GenericRenderer genericrenderer.Service `inject:"1"`
	Groups          groups.Service          `inject:"1"`
	Hierarchy       hierarchy.Service       `inject:"1"`
	Inputs          inputs.Service          `inject:"1"`
	Layout          layout.Service          `inject:"1"`
	NetSync         netsync.Service         `inject:"1"`
	Record          record.Service          `inject:"1"`
	Render          render.Service          `inject:"1"`
	Text            text.Service            `inject:"1"`
	Transform       transform.Service       `inject:"1"`
	Transition      transition.Service      `inject:"1"`
	UUID            uuid.Service            `inject:"1"`
}
