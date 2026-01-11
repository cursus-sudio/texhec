package render

import (
	"engine/services/ecs"
)

type Service interface {
	Color() ecs.ComponentsArray[ColorComponent]
	Mesh() ecs.ComponentsArray[MeshComponent]
	Texture() ecs.ComponentsArray[TextureComponent]
	TextureFrame() ecs.ComponentsArray[TextureFrameComponent]

	Error() error
}
