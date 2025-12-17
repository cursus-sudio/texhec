package render

import "engine/services/ecs"

type RenderTool interface {
	Render() Interface
}

type Interface interface {
	Color() ecs.ComponentsArray[ColorComponent]
	Mesh() ecs.ComponentsArray[MeshComponent]
	Texture() ecs.ComponentsArray[TextureComponent]

	Error() error
}
