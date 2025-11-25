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

func NewTool() ecs.ToolFactory[render.Tool] {
	return ecs.NewToolFactory(func(w ecs.World) render.Tool { return &tool{world: w} })
}

type transaction struct {
	*tool

	colorTransaction   ecs.ComponentsArrayTransaction[render.ColorComponent]
	meshTransaction    ecs.ComponentsArrayTransaction[render.MeshComponent]
	textureTransaction ecs.ComponentsArrayTransaction[render.TextureComponent]
}

func (t *tool) Transaction() render.Transaction {
	return &transaction{
		t,
		t.colorArray.Transaction(),
		t.meshArray.Transaction(),
		t.textureArray.Transaction(),
	}
}

//

type object struct {
	color   ecs.EntityComponent[render.ColorComponent]
	mesh    ecs.EntityComponent[render.MeshComponent]
	texture ecs.EntityComponent[render.TextureComponent]
}

func (t *transaction) GetObject(entity ecs.EntityID) render.Object {
	return &object{
		t.colorTransaction.GetEntityComponent(entity),
		t.meshTransaction.GetEntityComponent(entity),
		t.textureTransaction.GetEntityComponent(entity),
	}
}
func (t *transaction) Transactions() []ecs.AnyComponentsArrayTransaction {
	return []ecs.AnyComponentsArrayTransaction{
		t.colorTransaction,
		t.meshTransaction,
		t.textureTransaction,
	}
}
func (t *transaction) Flush() error {
	return ecs.FlushMany(t.Transactions()...)
}

//

func (o *object) Color() ecs.EntityComponent[render.ColorComponent]     { return o.color }
func (o *object) Mesh() ecs.EntityComponent[render.MeshComponent]       { return o.mesh }
func (o *object) Texture() ecs.EntityComponent[render.TextureComponent] { return o.texture }

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

func (*tool) Error() error {
	if glErr := gl.GetError(); glErr != gl.NO_ERROR {
		return fmt.Errorf("opengl error: %x %s\n", glErr, glErrorStrings[glErr])
	}
	return nil
}
