package shapes

import "frontend/services/console"

type collidersService struct {
	console console.Console
}

func newCollidersService(console console.Console) *collidersService {
	return &collidersService{
		console: console,
	}
}
