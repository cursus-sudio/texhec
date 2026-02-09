package uiservice

import (
	gameassets "core/assets"
	"core/modules/ui"
	"engine"
	"engine/modules/transform"
	"engine/modules/transition"
	"engine/services/ecs"
	"time"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type changing struct {
	isInitialized bool
	active        bool
	menu          ecs.EntityID
	childWrapper  ecs.EntityID
}

type service struct {
	GameAssets   gameassets.GameAssets `inject:"1"`
	engine.World `inject:"1"`

	animationDuration time.Duration

	bgTimePerFrame time.Duration

	uiCameraArray           ecs.ComponentsArray[ui.UiCameraComponent]
	cursorCameraArray       ecs.ComponentsArray[ui.CursorCameraComponent]
	animatedBackgroundArray ecs.ComponentsArray[ui.AnimatedBackgroundComponent]
	*changing
}

func NewService(
	c ioc.Dic,
	animationDuration time.Duration,
	bgTimePerFrame time.Duration,
) *service {
	t := ioc.GetServices[*service](c)
	t.animationDuration = animationDuration
	t.bgTimePerFrame = bgTimePerFrame

	t.uiCameraArray = ecs.GetComponentsArray[ui.UiCameraComponent](t.World)
	t.animatedBackgroundArray = ecs.GetComponentsArray[ui.AnimatedBackgroundComponent](t.World)
	t.cursorCameraArray = ecs.GetComponentsArray[ui.CursorCameraComponent](t.World)
	t.changing = &changing{}

	t.Logger.Warn(t.Init())

	return t
}

func (t *service) ResetChildWrapper() error {
	if !t.isInitialized {
		err := t.Init()
		if err != nil {
			return err
		}
	}

	for _, child := range t.Hierarchy.Children(t.childWrapper).GetIndices() {
		t.RemoveEntity(child)
	}
	return nil
}

func (t *service) UiCamera() ecs.ComponentsArray[ui.UiCameraComponent] { return t.uiCameraArray }
func (t *service) AnimatedBackground() ecs.ComponentsArray[ui.AnimatedBackgroundComponent] {
	return t.animatedBackgroundArray
}
func (s *service) CursorCamera() ecs.ComponentsArray[ui.CursorCameraComponent] {
	return s.cursorCameraArray
}

func (t *service) Show() ecs.EntityID {
	t.Logger.Warn(t.ResetChildWrapper())
	if !t.active {
		t.active = true
		events.Emit(t.Events, transition.NewTransitionEvent(
			t.menu,
			transform.NewPivotPoint(0, 1, .5),
			transform.NewPivotPoint(1, 1, .5),
			t.animationDuration,
		))
	}
	return t.childWrapper
}

func (t *service) Hide() {
	if t.active {
		t.active = false
		events.Emit(t.Events, transition.NewTransitionEvent(
			t.menu,
			transform.NewPivotPoint(1, 1, .5),
			transform.NewPivotPoint(0, 1, .5),
			t.animationDuration,
		))
	}
}
