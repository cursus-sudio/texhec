package genericrenderer

import (
	"engine/modules/camera"
	"engine/modules/groups"
	"engine/modules/render"
	"engine/modules/transform"
	"engine/services/ecs"
)

type System ecs.SystemRegister[World]

type GenericRendererTool interface {
	GenericRenderer() Interface
}

type World interface {
	ecs.World
	render.RenderTool
	camera.CameraTool
	groups.GroupsTool
	transform.TransformTool
}

type Interface interface {
	Pipeline() ecs.ComponentsArray[PipelineComponent]
}
