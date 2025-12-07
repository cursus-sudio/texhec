package syncpkg

import (
	"engine/modules/sync/internal/config"
	"engine/services/ecs"
	"reflect"

	"github.com/ogiusek/events"
)

type Config struct {
	config *config.Config
}

func NewConfig(
	isClient bool,
) Config {
	return Config{
		config: &config.Config{
			IsClient:       isClient,
			MaxPredictions: 60,
		},
	}
}

func AddEvent[EventType any](config Config) {
	config.config.Events = append(config.config.Events, reflect.TypeFor[EventType]())
	config.config.ListenToEvents = append(config.config.ListenToEvents, func(b events.Builder, f func(any)) {
		events.Listen(b, func(e EventType) { f(e) })
	})
}

func AddComponent[ComponentType any](config Config) {
	config.config.Components = append(config.config.Components, reflect.TypeFor[ComponentType]())
	config.config.ArraysOfComponents = append(config.config.ArraysOfComponents, func(w ecs.World) ecs.AnyComponentArray {
		return ecs.GetComponentsArray[ComponentType](w)
	})
}
