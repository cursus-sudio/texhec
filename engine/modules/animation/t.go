package animation

type AnimationState float64

// this is easing function id
// easing functions are registered in separate service
type EasingFunctionID uint32

type EasingFunction func(t AnimationState) AnimationState

//

// ctor will be type safe and will ensure type safety
type Transition struct {
	From, To       any // this should be of type T
	Start, End     AnimationState
	EasingFunction EasingFunctionID
}

func NewTransition[Component any](
	from, to Component,
	// start, end AnimationState,
	easingFunction EasingFunctionID,
) Transition {
	return Transition{
		From:           from,
		To:             to,
		Start:          0,
		End:            1,
		EasingFunction: easingFunction,
	}
}

func (t Transition) SetStart(start AnimationState) Transition {
	t.Start = start
	return t
}

func (t Transition) SetEnd(end AnimationState) Transition {
	t.End = end
	return t
}

//

type Event struct {
	Event   any
	Trigger AnimationState
}

func NewEvent[EventT any](event EventT, trigger AnimationState) Event {
	return Event{
		Event:   event,
		Trigger: trigger,
	}
}

//

type AnimationID uint32
type Animation struct {
	Events      []Event
	Transitions []Transition
}

func NewAnimation(
	events []Event,
	transitions []Transition,
) Animation {
	return Animation{
		Events:      events,
		Transitions: transitions,
	}
}
