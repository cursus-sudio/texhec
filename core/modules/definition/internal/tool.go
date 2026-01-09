package internal

import (
	"core/modules/definition"
	"engine/services/ecs"
	"sync"
)

type tool struct {
	definitionArray     ecs.ComponentsArray[definition.DefinitionComponent]
	definitionLinkArray ecs.ComponentsArray[definition.DefinitionLinkComponent]
}

func NewToolFactory() definition.ToolFactory {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w definition.World) definition.DefinitionTool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[tool](w); ok {
			return t
		}
		t := &tool{
			definitionArray:     ecs.GetComponentsArray[definition.DefinitionComponent](w),
			definitionLinkArray: ecs.GetComponentsArray[definition.DefinitionLinkComponent](w),
		}

		return t
	})
}

func (t *tool) Definition() definition.Interface { return t }
func (t *tool) Component() ecs.ComponentsArray[definition.DefinitionComponent] {
	return t.definitionArray
}
func (t *tool) Link() ecs.ComponentsArray[definition.DefinitionLinkComponent] {
	return t.definitionLinkArray
}
