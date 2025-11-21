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
				ecs.GetComponentsArray[animation.AnimationComponent](w),
				ecs.GetComponentsArray[groups.GroupsComponent](w),
				ecs.GetComponentsArray[collider.ColliderComponent](w),
				ecs.GetComponentsArray[inputs.MouseLeftClickComponent](w),
				ecs.GetComponentsArray[render.MeshComponent](w),
				ecs.GetComponentsArray[render.TextureComponent](w),
				ecs.GetComponentsArray[genericrenderer.PipelineComponent](w),
				ecs.GetComponentsArray[text.TextComponent](w),
				ecs.GetComponentsArray[text.BreakComponent](w),
				ecs.GetComponentsArray[text.TextAlignComponent](w),
				ecs.GetComponentsArray[text.FontSizeComponent](w),
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
		{"sese", inputs.QuitEvent{}},
		{"quit 2", inputs.QuitEvent{}},
		{"quit 2", inputs.QuitEvent{}},
	}
	optionsLen := float32(len(options))
	menuWrapper := s.w.NewEntity()
	ecs.SaveComponent(s.w, menuWrapper, transform.NewSize(100, (100/3)*optionsLen, 1))
	ecs.SaveComponent(s.w, menuWrapper, transform.NewPivotPoint(.5, 1, .5))
	ecs.SaveComponent(s.w, menuWrapper, transform.NewParent(e.Tile, transform.RelativePos))
	ecs.SaveComponent(s.w, menuWrapper, transform.NewParentPivotPoint(.5, 1, .5))

	menu := s.w.NewEntity()
	ecs.SaveComponent(s.w, menu, transform.NewSize(1, 1, 1))
	ecs.SaveComponent(s.w, menu, transform.NewParent(menuWrapper, transform.RelativePos|transform.RelativeSize))
	ecs.SaveComponent(s.w, menu, animation.NewAnimationComponent(gameassets.ShowMenuAnimation, time.Second))

	size := 1 / optionsLen

	for i, option := range options {
		entity := s.w.NewEntity()
		normalizedI := float32(i) / optionsLen
		if normalizedI == 0 {
		}

		// transform
		ecs.SaveComponent(s.w, entity, transform.NewSize(1, size, 1))
		ecs.SaveComponent(s.w, entity, transform.NewPivotPoint(.5, 1, .5))
		ecs.SaveComponent(s.w, entity, transform.NewParent(menu, transform.RelativePos|transform.RelativeSize))
		ecs.SaveComponent(s.w, entity, transform.NewParentPivotPoint(.5, 1-normalizedI, .5))

		//
		ecs.SaveComponent(s.w, entity, groups.EmptyGroups().Ptr().Enable(gamescene.GameGroup).Val())

		// mouse
		ecs.SaveComponent(s.w, entity, collider.NewCollider(gameassets.SquareColliderID))
		ecs.SaveComponent(s.w, entity, inputs.NewMouseLeftClick(option.OnClick))

		// texture
		ecs.SaveComponent(s.w, entity, render.NewMesh(gameassets.SquareMesh))
		ecs.SaveComponent(s.w, entity, render.NewTexture(gameassets.MountainTileTextureID))
		ecs.SaveComponent(s.w, entity, render.NewColor(mgl32.Vec4{1, normalizedI, 1, 1}))
		ecs.SaveComponent(s.w, entity, genericrenderer.PipelineComponent{})

		// text
		ecs.SaveComponent(s.w, entity, text.TextComponent{Text: option.Text})
		ecs.SaveComponent(s.w, entity, text.BreakComponent{Break: text.BreakNone})
		ecs.SaveComponent(s.w, entity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})
		ecs.SaveComponent(s.w, entity, text.FontSizeComponent{FontSize: 24})
	}
}
