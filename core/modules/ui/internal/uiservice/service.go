package uiservice

import (
	"core/modules/registry"
	"core/modules/ui"
	"engine"
	"engine/modules/transform"
	"engine/modules/transition"
	"engine/services/ecs"
	"time"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type menuComponent struct {
	Visible bool
}
type childrenComponent struct{}

type service struct {
	GameAssets   registry.Assets `inject:"1"`
	engine.World `inject:"1"`

	animationDuration time.Duration

	bgTimePerFrame time.Duration

	uiCameraArray           ecs.ComponentsArray[ui.UiCameraComponent]
	cursorCameraArray       ecs.ComponentsArray[ui.CursorCameraComponent]
	animatedBackgroundArray ecs.ComponentsArray[ui.AnimatedBackgroundComponent]
	menuArray               ecs.ComponentsArray[menuComponent]
	childrenWrapperArray    ecs.ComponentsArray[childrenComponent]
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
	t.cursorCameraArray = ecs.GetComponentsArray[ui.CursorCameraComponent](t.World)
	t.animatedBackgroundArray = ecs.GetComponentsArray[ui.AnimatedBackgroundComponent](t.World)
	t.menuArray = ecs.GetComponentsArray[menuComponent](t.World)
	t.childrenWrapperArray = ecs.GetComponentsArray[childrenComponent](t.World)

	events.Listen(t.EventsBuilder, func(e ui.HideUiEvent) {
		t.Hide()
	})

	t.EnsureExists()

	return t
}

func (t *service) ResetChildWrapper() {
	t.EnsureExists()

	for _, childWrapper := range t.childrenWrapperArray.GetEntities() {
		for _, child := range t.Hierarchy.Children(childWrapper).GetIndices() {
			t.RemoveEntity(child)
		}
	}
}

func (t *service) UiCamera() ecs.ComponentsArray[ui.UiCameraComponent] { return t.uiCameraArray }
func (t *service) AnimatedBackground() ecs.ComponentsArray[ui.AnimatedBackgroundComponent] {
	return t.animatedBackgroundArray
}
func (s *service) CursorCamera() ecs.ComponentsArray[ui.CursorCameraComponent] {
	return s.cursorCameraArray
}

func (t *service) Show() []ecs.EntityID {
	t.ResetChildWrapper()

	for _, menu := range t.menuArray.GetEntities() {
		if component, _ := t.menuArray.Get(menu); !component.Visible {
			t.menuArray.Set(menu, menuComponent{true})
			events.Emit(t.Events, transition.NewTransitionEvent(
				menu,
				transform.NewPivotPoint(0, 1, .5),
				transform.NewPivotPoint(1, 1, .5),
				t.animationDuration,
			))
		}
	}
	return t.childrenWrapperArray.GetEntities()
}

func (t *service) Hide() {
	t.EnsureExists()

	for _, menu := range t.menuArray.GetEntities() {
		if component, _ := t.menuArray.Get(menu); component.Visible {
			t.menuArray.Set(menu, menuComponent{false})
			events.Emit(t.Events, transition.NewTransitionEvent(
				menu,
				transform.NewPivotPoint(1, 1, .5),
				transform.NewPivotPoint(0, 1, .5),
				t.animationDuration,
			))
		}
	}
}
