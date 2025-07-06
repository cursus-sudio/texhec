module frontend

go 1.24.3

require shared v0.0.0

require (
	github.com/ogiusek/events v1.0.2
	github.com/ogiusek/ioc/v2 v2.0.11
	github.com/ogiusek/null v1.1.0
	github.com/veandco/go-sdl2 v0.4.40
)

require (
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20250301202403-da16c1255728 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/ogiusek/lockset v1.0.1 // indirect
	github.com/ogiusek/relay/v2 v2.0.6 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
)

replace shared => ../shared
