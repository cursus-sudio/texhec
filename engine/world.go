package engine

import (
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/connection"
	"engine/modules/genericrenderer"
	"engine/modules/groups"
	"engine/modules/hierarchy"
	"engine/modules/inputs"
	"engine/modules/netsync"
	"engine/modules/record"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/modules/transition"
	"engine/modules/uuid"
	"engine/services/ecs"
)

type World interface {
	ecs.World
	camera.CameraTool
	collider.ColliderTool
	connection.ConnectionTool
	genericrenderer.GenericRendererTool
	groups.GroupsTool
	hierarchy.HierarchyTool
	inputs.InputsTool
	netsync.NetSyncTool
	record.RecordTool
	render.RenderTool
	text.TextTool
	transform.TransformTool
	transition.TransitionTool
	uuid.UUIDTool
}
