package render

import (
	"engine/modules/transform"
	"engine/services/ecs"
)

type ToolFactory ecs.ToolFactory[World, RenderTool]
type RenderTool interface {
	Render() Interface
}
type World interface {
	ecs.World
	transform.TransformTool
}
type Interface interface {
	Color() ecs.ComponentsArray[ColorComponent]
	Mesh() ecs.ComponentsArray[MeshComponent]
	Texture() ecs.ComponentsArray[TextureComponent]
	TextureFrame() ecs.ComponentsArray[TextureFrameComponent]

	Error() error
}
