package internal

import (
	"engine/modules/render"
	"engine/services/ecs"
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type tool struct {
	world        ecs.World
	colorArray   ecs.ComponentsArray[render.ColorComponent]
	meshArray    ecs.ComponentsArray[render.MeshComponent]
	textureArray ecs.ComponentsArray[render.TextureComponent]
}

func NewTool() ecs.ToolFactory[render.RenderTool] {
	return ecs.NewToolFactory(func(w ecs.World) render.RenderTool {
		return &tool{
			w,
			ecs.GetComponentsArray[render.ColorComponent](w),
			ecs.GetComponentsArray[render.MeshComponent](w),
			ecs.GetComponentsArray[render.TextureComponent](w),
		}
	})
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

func (t *tool) Render() render.Interface {
	return t
}

func (t *tool) Color() ecs.ComponentsArray[render.ColorComponent] {
	return t.colorArray
}
func (t *tool) Mesh() ecs.ComponentsArray[render.MeshComponent] {
	return t.meshArray
}
func (t *tool) Texture() ecs.ComponentsArray[render.TextureComponent] {
	return t.textureArray
}

func (*tool) Error() error {
	if glErr := gl.GetError(); glErr != gl.NO_ERROR {
		return fmt.Errorf("opengl error: %x %s\n", glErr, glErrorStrings[glErr])
	}
	return nil
}
