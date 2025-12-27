package internal

import (
	gameassets "core/assets"
	"core/modules/settings"
	gamescenes "core/scenes"
	"engine/modules/audio"
	"engine/modules/collider"
	"engine/modules/genericrenderer"
	"engine/modules/groups"
	"engine/modules/inputs"
	"engine/modules/render"
	"engine/modules/scenes"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/assets"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

// 1. settings text
// 2. quit button

type system struct {
	assets     assets.Assets
	gameAssets gameassets.GameAssets

	logger logger.Logger
	settings.World
}

type temporaryToggleColorComponent struct{}

func NewSystem(
	assets assets.Assets,
	logger logger.Logger,
	gameAssets gameassets.GameAssets,
) settings.System {
	return ecs.NewSystemRegister(func(world settings.World) error {
		s := system{
			assets,
			gameAssets,

			logger,
			world,
		}

		events.ListenE(world.EventsBuilder(), func(event settings.EnterSettingsForParentEvent) error {
			return s.ListenRender(event.Parent)
		})
		events.Listen(world.EventsBuilder(), s.ListenOnTick)
		events.Listen(world.EventsBuilder(), func(settings.EnterSettingsEvent) {
			event := settings.EnterSettingsForParentEvent{
				Parent: s.Ui().Show(),
			}
			events.Emit(world.Events(), event)
		})

		return nil
	})
}

func (s system) ListenOnTick(frames.TickEvent) {
	toggleArray := ecs.GetComponentsArray[temporaryToggleColorComponent](s)
	for _, entity := range toggleArray.GetEntities() {
		color, ok := s.Render().Color().Get(entity)
		if !ok {
			color.Color = mgl32.Vec4{1, 1, 1, 1}
		}

		color.Color[1] = 1 - color.Color[1]
		color.Color[2] = 1 - color.Color[2]

		s.Render().Color().Set(entity, color)
	}

}

func (s system) ListenRender(parent ecs.EntityID) error {
	// render
	// collider
	// click

	// changes
	labelEntity := s.NewEntity()
	s.Hierarchy().SetParent(labelEntity, parent)
	s.Groups().Inherit().Set(labelEntity, groups.InheritGroupsComponent{})

	s.Transform().Parent().Set(labelEntity, transform.NewParent(transform.RelativePos|transform.RelativeSizeX))
	s.Transform().Size().Set(labelEntity, transform.NewSize(1, 50, 1))

	s.Text().Content().Set(labelEntity, text.TextComponent{Text: "SETTINGS"})
	s.Text().FontSize().Set(labelEntity, text.FontSizeComponent{FontSize: 25})
	s.Text().Align().Set(labelEntity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})

	//

	type Button struct {
		text  string
		event any
	}
	btns := []Button{
		{"SHOT", audio.NewPlayEvent(gamescenes.EffectChannel, s.gameAssets.ExampleAudio)},
		{"QUIT", scenes.NewChangeSceneEvent(gamescenes.MenuID)},
	}

	btnAsset, err := assets.GetAsset[render.TextureAsset](s.assets, s.gameAssets.Hud.Btn)
	if err != nil {
		return err
	}
	btnAspectRatio := btnAsset.AspectRatio()

	for _, btn := range btns {
		btnEntity := s.NewEntity()
		s.Hierarchy().SetParent(btnEntity, parent)
		s.Groups().Inherit().Set(btnEntity, groups.InheritGroupsComponent{})

		ecs.GetComponentsArray[temporaryToggleColorComponent](s).Set(btnEntity, temporaryToggleColorComponent{})

		s.Transform().AspectRatio().Set(btnEntity, transform.NewAspectRatio(float32(btnAspectRatio.Dx()), float32(btnAspectRatio.Dy()), 0, transform.PrimaryAxisX))
		s.Transform().Parent().Set(btnEntity, transform.NewParent(transform.RelativePos|transform.RelativeSizeX))
		s.Transform().MaxSize().Set(btnEntity, transform.NewMaxSize(0, 50, 0))
		s.Transform().Size().Set(btnEntity, transform.NewSize(1, 50, 1))

		s.Render().Mesh().Set(btnEntity, render.NewMesh(s.gameAssets.SquareMesh))
		s.Render().Texture().Set(btnEntity, render.NewTexture(s.gameAssets.Hud.Btn))
		s.GenericRenderer().Pipeline().Set(btnEntity, genericrenderer.PipelineComponent{})

		s.Text().Content().Set(btnEntity, text.TextComponent{Text: btn.text})
		s.Text().FontSize().Set(btnEntity, text.FontSizeComponent{FontSize: 25})
		s.Text().Align().Set(btnEntity, text.TextAlignComponent{Vertical: .5, Horizontal: .5})

		s.Inputs().MouseLeft().Set(btnEntity, inputs.NewMouseLeftClick(btn.event))
		s.Inputs().KeepSelected().Set(btnEntity, inputs.KeepSelectedComponent{})
		s.Collider().Component().Set(btnEntity, collider.NewCollider(s.gameAssets.SquareCollider))
	}

	// go func() {
	// 	return
	// 	for {
	// 		time.Sleep(time.Millisecond * 500)
	// 		// children := s.Hierarchy().Children(parent).GetIndices()
	// 		children := []ecs.EntityID{5}
	// 		type data struct {
	// 			parentPos        transform.PosComponent
	// 			parentPivotPoint transform.ParentPivotPointComponent
	//
	// 			pos        transform.PosComponent
	// 			size       transform.SizeComponent
	// 			pivotPoint transform.PivotPointComponent
	//
	// 			absolutePos transform.AbsolutePosComponent
	// 		}
	// 		dataPrinted := []data{}
	// 		for _, child := range children {
	// 			parentPos, _ := s.Transform().Pos().Get(parent)
	// 			parentPivotPoint, _ := s.Transform().ParentPivotPoint().Get(child)
	//
	// 			pos, _ := s.Transform().Pos().Get(child)
	// 			size, _ := s.Transform().Size().Get(child)
	// 			pivotPoint, _ := s.Transform().PivotPoint().Get(child)
	//
	// 			absolutePos, _ := s.Transform().AbsolutePos().Get(child)
	//
	// 			dataPrinted = append(dataPrinted, data{
	// 				parentPos:        parentPos,
	// 				parentPivotPoint: parentPivotPoint,
	//
	// 				pos:        pos,
	// 				size:       size,
	// 				pivotPoint: pivotPoint,
	//
	// 				absolutePos: absolutePos,
	// 			})
	// 		}
	// 		// s.logger.Info("children are %v with poses %v with codes\n%v", children, poses, posesCodes)
	// 		s.logger.Info("children are %v with data %v", children, dataPrinted)
	// 		// incorrect: [{{[0 0 0]} {[0.5 0 0.5]} {[0 85 0]} {[1 50 1]} {[0.5 1 0.5]} {[0 -89.606064 0]}}]
	// 		// correct  : [{{[0 0 0]} {[0.5 0 0.5]} {[0 85 0]} {[1 50 1]} {[0.5 1 0.5]} {[0 -40 0]}}]
	// 	}
	// }()
	// this line changes layout for some reason
	// s.logger.Info("before call")
	// s.Hierarchy().Children(0) // this lines destroys transform for uknown reason
	// s.Transform().AbsoluteRotation().GetEntities() // and this repairs it back (no matter where placed) HELP
	// s.logger.Info("after call")
	// its the mix of (removing any of these solves the issue):
	// - aspect ratio
	// - hierarchy
	// - layout
	// and it causes:
	// - absolute pos to do not update despite having the same parentPos, parentPivot, pos, size, pivot
	// s.logger.Info("children are %v", s.Hierarchy().Children(parent).GetIndices())
	return nil
}
