package service

import (
	"engine/modules/render"
	"engine/services/ecs"
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type service struct {
	World             ecs.World `inject:"1"`
	colorArray        ecs.ComponentsArray[render.ColorComponent]
	meshArray         ecs.ComponentsArray[render.MeshComponent]
	textureArray      ecs.ComponentsArray[render.TextureComponent]
	textureFrameArray ecs.ComponentsArray[render.TextureFrameComponent]

	// directArray     ecs.ComponentsArray[render.DirectComponent]
	// instancingArray ecs.ComponentsArray[render.InstancingComponent]
}

func NewService(c ioc.Dic) render.Service {
	s := ioc.GetServices[*service](c)
	s.colorArray = ecs.GetComponentsArray[render.ColorComponent](s.World)
	s.meshArray = ecs.GetComponentsArray[render.MeshComponent](s.World)
	s.textureArray = ecs.GetComponentsArray[render.TextureComponent](s.World)
	s.textureFrameArray = ecs.GetComponentsArray[render.TextureFrameComponent](s.World)

	// s.directArray = ecs.GetComponentsArray[render.DirectComponent](s.World)
	// s.instancingArray = ecs.GetComponentsArray[render.InstancingComponent](s.World)

	// defaults
	s.colorArray.SetEmpty(render.NewColor(mgl32.Vec4{1, 1, 1, 1}))
	// no default mesh
	// no default texture
	s.textureFrameArray.SetEmpty(render.NewTextureFrameComponent(0))

	return s
}

//

var GlErrorStrings = map[uint32]string{
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

//	func (t *service) Direct() ecs.ComponentsArray[render.DirectComponent] {
//		return t.directArray
//	}
//
//	func (t *service) Instancing() ecs.ComponentsArray[render.InstancingComponent] {
//		return t.instancingArray
//	}
func (t *service) Render(entity ecs.EntityID) {
	// when instancing will be working change this to instancing by default
	// t.directArray.Set(entity, render.DirectComponent{})
}

func (*service) Error() error {
	if glErr := gl.GetError(); glErr != gl.NO_ERROR {
		return fmt.Errorf("opengl error: %x %s", glErr, GlErrorStrings[glErr])
	}
	return nil
}
