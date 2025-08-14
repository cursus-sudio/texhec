package grid

import "frontend/services/frames"

type GridSystem struct{}

func NewGridSystem() GridSystem {
	return GridSystem{}
}

func (s *GridSystem) Listen(e frames.FrameEvent) {
	// creates vao
	// sends data to shader
	// renders
}
