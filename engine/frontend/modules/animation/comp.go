package animation

import (
	"time"
)

type AnimationComponent struct {
	AnimationID AnimationID
	PreviousState,
	State AnimationState
	Duration time.Duration

	// if component isn't removed then state keeps being update.
	// this allows looping
}

func NewAnimationComponent(
	animationID AnimationID,
	duration time.Duration,
) AnimationComponent {
	return AnimationComponent{
		AnimationID:   animationID,
		PreviousState: 0,
		State:         0,
		Duration:      duration,
	}
}

func (c *AnimationComponent) AddElapsedTime(time time.Duration) {
	ratio := time.Seconds() / c.Duration.Seconds()
	c.PreviousState = c.State
	c.State += AnimationState(ratio)
	c.State = min(c.State, 1)
}

// //
//
// type LoopComponent struct {
// }
