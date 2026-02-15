package systems

import (
	"core/modules/registry"
	"core/modules/ui"
	"engine"
	"engine/modules/assets"
	"engine/modules/render"
	"engine/modules/transform"
	"engine/modules/transition"
	"engine/services/ecs"
	"engine/services/frames"
	"time"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type UpdateBgEvent struct{}

type System struct {
	GameAssets   registry.Assets `inject:"1"`
	engine.World `inject:"1"`
	Ui           ui.Service `inject:"1"`

	blueprint     ecs.EntityID
	bgDirtySet    ecs.DirtySet
	transitionArr ecs.ComponentsArray[transition.TransitionComponent[render.TextureFrameComponent]]

	bgTimePerFrame time.Duration
	bgTexture      int

	backgrounds       []assets.ID
	backgroundsFrames []int
}

func NewSystem(c ioc.Dic, bgTimePerFrame time.Duration) ui.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*System](c)

		s.blueprint = s.NewEntity()
		s.Ui.AnimatedBackground().Set(s.blueprint, ui.AnimatedBackgroundComponent{})

		s.bgDirtySet = ecs.NewDirtySet()
		s.Ui.AnimatedBackground().AddDirtySet(s.bgDirtySet)

		s.transitionArr = ecs.GetComponentsArray[transition.TransitionComponent[render.TextureFrameComponent]](s.World)
		s.bgTimePerFrame = bgTimePerFrame
		s.bgTexture = 0

		s.backgrounds = []assets.ID{
			s.GameAssets.Hud.Background2,
			s.GameAssets.Hud.Background1,
			s.GameAssets.Hud.Background1,
			s.GameAssets.Hud.Background1,
		}

		s.backgroundsFrames = make([]int, 0, len(s.backgrounds))
		for _, bg := range s.backgrounds {
			texture, err := assets.GetAsset[render.TextureAsset](s.Assets, bg)
			if err != nil {
				return err
			}
			s.backgroundsFrames = append(s.backgroundsFrames, len(texture.Images()))
		}

		s.Transform.Parent().BeforeGet(s.BeforeGet)
		s.Render.Mesh().BeforeGet(s.BeforeGet)
		s.Render.Texture().BeforeGet(s.BeforeGet)
		s.Render.TextureFrame().BeforeGet(s.BeforeGet)

		events.Listen(s.EventsBuilder, s.ListenUpdateBg)
		events.Emit(s.Events, UpdateBgEvent{})
		return nil
	})
}

func (s *System) BeforeGet() {
	entities := s.bgDirtySet.Get()
	if len(entities) == 0 {
		return
	}

	texture, _ := s.Render.Texture().Get(s.blueprint)
	transitionComp, _ := s.transitionArr.Get(s.blueprint)
	for _, entity := range entities {
		if entity == s.blueprint {
			continue
		}
		if _, ok := s.Ui.AnimatedBackground().Get(entity); !ok {
			continue
		}
		if _, ok := s.transitionArr.Get(entity); ok {
			continue
		}
		s.Transform.Parent().Set(entity, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))
		if entity != s.blueprint {
			s.Render.Mesh().Set(entity, render.NewMesh(s.GameAssets.SquareMesh))
		}
		s.Render.Texture().Set(entity, texture)
		s.transitionArr.Set(entity, transitionComp)
	}
}

func (s *System) ListenFrame(e frames.FrameEvent) {
}

func (s *System) ListenUpdateBg(event UpdateBgEvent) {
	i := s.bgTexture % len(s.backgrounds)
	s.bgTexture = i
	bg, size := s.backgrounds[i], s.backgroundsFrames[i]
	duration := s.bgTimePerFrame * time.Duration(size)

	for _, entity := range s.Ui.AnimatedBackground().GetEntities() {
		s.Transform.Parent().Set(entity, transform.NewParent(transform.RelativePos|transform.RelativeSizeXY))
		if entity != s.blueprint {
			s.Render.Mesh().Set(entity, render.NewMesh(s.GameAssets.SquareMesh))
		}
		s.Render.Texture().Set(entity, render.NewTexture(bg))
		s.transitionArr.Set(entity, transition.NewTransition(
			render.NewTextureFrame(0),
			render.NewTextureFrame(1),
			duration,
		))
	}

	events.Emit(s.Events, transition.NewDelayedEvent(UpdateBgEvent{}, duration))
	s.bgTexture += 1
}
