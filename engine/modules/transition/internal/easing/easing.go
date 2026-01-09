package easing

import (
	"engine/modules/transition"
	"engine/services/datastructures"
)

type easingService struct {
	easingFunctions datastructures.SparseArray[transition.EasingID, transition.EasingFunction]
}

func NewEasingService() transition.EasingService {
	return &easingService{
		easingFunctions: datastructures.NewSparseArray[transition.EasingID, transition.EasingFunction](),
	}
}

func (s *easingService) Set(id transition.EasingID, fn transition.EasingFunction) {
	s.easingFunctions.Set(id, fn)
}
func (s *easingService) Get(id transition.EasingID) (transition.EasingFunction, bool) {
	return s.easingFunctions.Get(id)
}
