package internal

import (
	"engine/modules/render"
	"engine/services/ecs"
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type service struct {
	world             ecs.World
	colorArray        ecs.ComponentsArray[render.ColorComponent]
	meshArray         ecs.ComponentsArray[render.MeshComponent]
	textureArray      ecs.ComponentsArray[render.TextureComponent]
	textureFrameArray ecs.ComponentsArray[render.TextureFrameComponent]
}

func NewService(
	world ecs.World,
) render.Service {
	return &service{
		world,
		ecs.GetComponentsArray[render.ColorComponent](world),
		ecs.GetComponentsArray[render.MeshComponent](world),
		ecs.GetComponentsArray[render.TextureComponent](world),
		ecs.GetComponentsArray[render.TextureFrameComponent](world),
	}
}

//

var glErrorStrings = map[uint32]string{
	gl.NO_ERROR:                      "GL_NO_ERROR",
	gl.INVALID_ENUM:                  "GL_INVALID_ENUM",
	gl.INVALID_VALUE:                 "GL_INVALID_VALUE",
	gl.INVALID_OPERATION:             "GL_INVALID_OPERATION",
	gl.STACK_OVERFLOW:                "GL_STACK_OVERFLOW",
	gl.STACK_UNDERFLOW:               "GL_STACK_UNDERFLOW",
	gl.OUT_OF_MEMORY:                 "GL_OUT_OF_MEMORY",
	gl.INVALID_FRAMEBUFFER_OPERATION: "GL_INVALID_FRAMEBUFFER_OPERATION",
	gl.CONTEXT_LOST:                  "GL_CONTEXT_LOST",
	// gl.TABLE_TOO_LARGE:               "GL_TABLE_TOO_LARGE", // Less common in modern GL
}

func (t *service) Color() ecs.ComponentsArray[render.ColorComponent] {
	return t.colorArray
}
func (t *service) Mesh() ecs.ComponentsArray[render.MeshComponent] {
	return t.meshArray
}
func (t *service) Texture() ecs.ComponentsArray[render.TextureComponent] {
	return t.textureArray
}
func (t *service) TextureFrame() ecs.ComponentsArray[render.TextureFrameComponent] {
	return t.textureFrameArray
}

func (*service) Error() error {
	if glErr := gl.GetError(); glErr != gl.NO_ERROR {
		return fmt.Errorf("opengl error: %x %s", glErr, glErrorStrings[glErr])
	}
	return nil
}
