package engine

import (
	"engine/modules/animation"
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/connection"
	"engine/modules/genericrenderer"
	"engine/modules/groups"
	"engine/modules/hierarchy"
	"engine/modules/inputs"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/modules/uuid"
	"engine/services/ecs"
)

type World interface {
	ecs.World
	animation.AnimationTool
	camera.CameraTool
	collider.ColliderTool
	connection.ConnectionTool
	genericrenderer.GenericRendererTool
	groups.GroupsTool
	hierarchy.HierarchyTool
	inputs.InputsTool
	render.RenderTool
	text.TextTool
	transform.TransformTool
	uuid.UUIDTool
}
