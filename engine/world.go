package engine

import (
	"engine/modules/assets"
	"engine/modules/audio"
	"engine/modules/batcher"
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/connection"
	"engine/modules/groups"
	"engine/modules/hierarchy"
	"engine/modules/inputs"
	"engine/modules/layout"
	"engine/modules/netsync"
	"engine/modules/noise"
	"engine/modules/record"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/modules/transition"
	"engine/modules/uuid"
	"engine/services/clock"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"

	"github.com/ogiusek/events"
)

type World struct {
	ecs.World  `inject:"1"`
	Assets     assets.Service     `inject:"1"`
	Audio      audio.Service      `inject:"1"`
	Batcher    batcher.Service    `inject:"1"`
	Camera     camera.Service     `inject:"1"`
	Collider   collider.Service   `inject:"1"`
	Connection connection.Service `inject:"1"`
	Groups     groups.Service     `inject:"1"`
	Hierarchy  hierarchy.Service  `inject:"1"`
	Inputs     inputs.Service     `inject:"1"`
	Layout     layout.Service     `inject:"1"`
	NetSync    netsync.Service    `inject:"1"`
	Noise      noise.Service      `inject:"1"`
	Record     record.Service     `inject:"1"`
	Render     render.Service     `inject:"1"`
	Text       text.Service       `inject:"1"`
	Transform  transform.Service  `inject:"1"`
	Transition transition.Service `inject:"1"`
	UUID       uuid.Service       `inject:"1"`

	Clock  clock.Clock   `inject:"1"`
	Logger logger.Logger `inject:"1"`

	EventsBuilder events.Builder `inject:"1"`
	Events        events.Events  `inject:"1"`

	Window window.Api `inject:"1"`
}
