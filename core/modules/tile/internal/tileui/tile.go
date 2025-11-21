package tileui

import (
	gameassets "core/assets"
	"core/modules/tile"
	gamescene "core/scenes/game"
	"frontend/modules/animation"
	"frontend/modules/collider"
	"frontend/modules/genericrenderer"
	"frontend/modules/groups"
	"frontend/modules/inputs"
	"frontend/modules/render"
	"frontend/modules/text"
	"frontend/modules/transform"
	"shared/services/ecs"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s tile.System) tile.System {
		toolFactory := ioc.Get[ecs.ToolFactory[transform.TransformTool]](c)
		return ecs.NewSystemRegister(func(w ecs.World) error {
			if err := s.Register(w); err != nil {
				return err
			}
			tool := toolFactory.Build(w)
			s := system{
				w,
				tool.Transaction(),
				ecs.GetComponentsArray[animation.AnimationComponent](w.Components()),
				ecs.GetComponentsArray[groups.GroupsComponent](w.Components()),
				ecs.GetComponentsArray[collider.ColliderComponent](w.Components()),
				ecs.GetComponentsArray[inputs.MouseLeftClickComponent](w.Components()),
				ecs.GetComponentsArray[render.MeshComponent](w.Components()),
				ecs.GetComponentsArray[render.TextureComponent](w.Components()),
				ecs.GetComponentsArray[genericrenderer.PipelineComponent](w.Components()),
				ecs.GetComponentsArray[text.TextComponent](w.Components()),
				ecs.GetComponentsArray[text.BreakComponent](w.Components()),
				ecs.GetComponentsArray[text.TextAlignComponent](w.Components()),
				ecs.GetComponentsArray[text.FontSizeComponent](w.Components()),
			}
			events.Listen(w.EventsBuilder(), s.Listen)
			return nil
		})
	})
}

type system struct {
	w ecs.World

	transform           transform.TransformTransaction
	animationArray      ecs.ComponentsArray[animation.AnimationComponent]
	groupsArray         ecs.ComponentsArray[groups.GroupsComponent]
	colliderArray       ecs.ComponentsArray[collider.ColliderComponent]
	mouseLeftClickArray ecs.ComponentsArray[inputs.MouseLeftClickComponent]
	meshArray           ecs.ComponentsArray[render.MeshComponent]
	textureArray        ecs.ComponentsArray[render.TextureComponent]
	pipelineArray       ecs.ComponentsArray[genericrenderer.PipelineComponent]
	textArray           ecs.ComponentsArray[text.TextComponent]
	breakArray          ecs.ComponentsArray[text.BreakComponent]
	textAlignArray      ecs.ComponentsArray[text.TextAlignComponent]
	fontSizeArray       ecs.ComponentsArray[text.FontSizeComponent]
}

type Option struct {
	Text    string
	OnClick any // any event
}

func (s *system) Listen(e tile.TileClickEvent) {
	options := []Option{
		{"quit", inputs.QuitEvent{}},
		{"quit 2", inputs.QuitEvent{}},
	}
	optionsLen := float32(len(options))
	menuWrapper := s.w.NewEntity()
	ecs.SaveComponent(s.w.Components(), menuWrapper, transform.NewSize(100, 50*optionsLen, 1))
	ecs.SaveComponent(s.w.Components(), menuWrapper, transform.NewPivotPoint(0, 1, .5))
	ecs.SaveComponent(s.w.Components(), menuWrapper, transform.NewParent(e.Tile, transform.RelativePos))
	ecs.SaveComponent(s.w.Components(), menuWrapper, transform.NewParentPivotPoint(1, 1, .5))

	menu := s.w.NewEntity()
	ecs.SaveComponent(s.w.Components(), menu, transform.NewSize(1, 1, 1))
	ecs.SaveComponent(s.w.Components(), menu, transform.NewParent(menuWrapper, transform.RelativePos|transform.RelativeSize))
	ecs.SaveComponent(s.w.Components(), menu, animation.NewAnimationComponent(gameassets.ShowMenuAnimation, time.Second))

	size := 1 / optionsLen

	for i, option := range options {
		entity := s.w.NewEntity()
		normalizedI := float32(i) / optionsLen
		if normalizedI == 0 {
		}

		// transform
		ecs.SaveComponent(s.w.Components(), entity, transform.NewSize(1, size, 1))
		ecs.SaveComponent(s.w.Components(), entity, transform.NewPivotPoint(.5, 1, .5))
		ecs.SaveComponent(s.w.Components(), entity, transform.NewParent(menu, transform.RelativePos|transform.RelativeSize))
		ecs.SaveComponent(s.w.Components(), entity, transform.NewParentPivotPoint(.5, 1-normalizedI, .5))

		//
		ecs.SaveComponent(s.w.Components(), entity, groups.EmptyGroups().Ptr().Enable(gamescene.GameGroup).Val())

		// mouse
		ecs.SaveComponent(s.w.Components(), entity, collider.NewCollider(gameassets.SquareColliderID))
		ecs.SaveComponent(s.w.Components(), entity, inputs.NewMouseLeftClick(option.OnClick))

		// texture
		ecs.SaveComponent(s.w.Components(), entity, render.NewMesh(gameassets.SquareMesh))
		ecs.SaveComponent(s.w.Components(), entity, render.NewTexture(gameassets.MountainTileTextureID))
		ecs.SaveComponent(s.w.Components(), entity, render.NewColor(mgl32.Vec4{1, normalizedI, 1, 1}))
		ecs.SaveComponent(s.w.Components(), entity, genericrenderer.PipelineComponent{})

		// text
		ecs.SaveComponent(s.w.Components(), entity, text.TextComponent{Text: option.Text})
		ecs.SaveComponent(s.w.Components(), entity, text.BreakComponent{Break: text.BreakNone})
		ecs.SaveComponent(s.w.Components(), entity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
		ecs.SaveComponent(s.w.Components(), entity, text.FontSizeComponent{FontSize: 24})
	}
}
