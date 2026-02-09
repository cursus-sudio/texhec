package systems

import (
	gameassets "core/assets"
	"core/modules/ui"
	"engine"
	"engine/modules/render"
	"engine/modules/transform"
	"engine/modules/transition"
	"engine/services/assets"
	"engine/services/ecs"
	"time"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type UpdateBgEvent struct{}

type System struct {
	GameAssets   gameassets.GameAssets `inject:"1"`
	engine.World `inject:"1"`
	Ui           ui.Service `inject:"1"`

	bgTimePerFrame time.Duration
	bgTexture      int

	backgrounds       []assets.AssetID
	backgroundsFrames []int
}

func NewSystem(c ioc.Dic, bgTimePerFrame time.Duration) ui.System {
	s := ioc.GetServices[*System](c)
	s.bgTimePerFrame = bgTimePerFrame
	s.bgTexture = 0

	s.backgrounds = []assets.AssetID{
		s.GameAssets.Hud.Background1,
		s.GameAssets.Hud.Background1,
		s.GameAssets.Hud.Background1,
		s.GameAssets.Hud.Background2,
		s.GameAssets.Hud.Background2,
		s.GameAssets.Hud.Background1,
	}

	return ecs.NewSystemRegister(func() error {
		s.backgroundsFrames = make([]int, 0, len(s.backgrounds))
		for _, bg := range s.backgrounds {
			texture, err := assets.GetAsset[render.TextureAsset](s.Assets, bg)
			if err != nil {
				return err
			}
			s.backgroundsFrames = append(s.backgroundsFrames, len(texture.Images()))
		}
		return s.Init()
	})
}

func (s *System) Init() error {
	dirtySet := ecs.NewDirtySet()
	s.Ui.AnimatedBackground().AddDirtySet(dirtySet)

	blueprint := s.NewEntity()
	s.Ui.AnimatedBackground().Set(blueprint, ui.AnimatedBackgroundComponent{})

	transitionArr := ecs.GetComponentsArray[transition.TransitionComponent[render.TextureFrameComponent]](s.World)

	beforeGet := func() {
		entities := dirtySet.Get()
		if len(entities) == 0 {
			return
		}

		texture, _ := s.Render.Texture().Get(blueprint)
		transitionComp, _ := transitionArr.Get(blueprint)
		for _, entity := range entities {
			if entity == blueprint {
				continue
			}
			if _, ok := s.Ui.AnimatedBackground().Get(entity); !ok {
				continue
			}
			if _, ok := transitionArr.Get(entity); ok {
				continue
			}
			s.Transform.Parent().Set(entity, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))
			s.Render.Mesh().Set(entity, render.NewMesh(s.GameAssets.SquareMesh))
			s.Render.Texture().Set(entity, texture)
			s.Transition.Easing().Set(entity, transition.NewEasing(gameassets.MyEasingFunction))
			transitionArr.Set(entity, transitionComp)
		}
	}

	s.Transform.Parent().BeforeGet(beforeGet)
	s.Render.Mesh().BeforeGet(beforeGet)
	s.Render.Texture().BeforeGet(beforeGet)
	s.Render.TextureFrame().BeforeGet(beforeGet)

	//

	events.Listen(s.EventsBuilder, func(u UpdateBgEvent) {
		i := s.bgTexture % len(s.backgrounds)
		s.bgTexture = i
		bg, size := s.backgrounds[i], s.backgroundsFrames[i]
		duration := s.bgTimePerFrame * time.Duration(size)

		for _, entity := range s.Ui.AnimatedBackground().GetEntities() {
			s.Transform.Parent().Set(entity, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))
			s.Render.Mesh().Set(entity, render.NewMesh(s.GameAssets.SquareMesh))
			s.Render.Texture().Set(entity, render.NewTexture(bg))
			s.Transition.Easing().Set(entity, transition.NewEasing(gameassets.MyEasingFunction))
			transitionArr.Set(entity, transition.NewTransition(
				render.NewTextureFrame(0),
				render.NewTextureFrame(1),
				duration,
			))
		}

		events.Emit(s.Events, transition.NewDelayedEvent(UpdateBgEvent{}, duration))
		s.bgTexture += 1
	})

	events.Emit(s.Events, UpdateBgEvent{})
	return nil
}
