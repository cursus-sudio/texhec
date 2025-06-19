package ecs

import "log"

type systemsImpl struct {
	updateSystems []System
	drawSystems   []System
}

func newSystems() *systemsImpl {
	return &systemsImpl{
		updateSystems: make([]System, 0),
		drawSystems:   make([]System, 0),
	}
}

func (systems *systemsImpl) LoadSystem(system System, systemType SystemType) {
	if system == nil {
		log.Panic("tried to add not implemented system")
	}

	switch systemType {
	case UpdateSystem:
		systems.updateSystems = append(systems.updateSystems, system)
		break
	case DrawSystem:
		systems.drawSystems = append(systems.drawSystems, system)
		break
	default:
		log.Panicf("not yet implemented loading system type")
		break
	}
}

func (systems *systemsImpl) Update(args Args) {
	for _, system := range systems.updateSystems {
		system.Update(args)
	}
	for _, system := range systems.drawSystems {
		system.Update(args)
	}
}
