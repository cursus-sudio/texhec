package ecs_test

import (
	"frontend/services/ecs"
	"testing"
)

func TestRegistry(t *testing.T) {
	type Register struct{ value int }
	world := ecs.NewWorld()

	r1 := Register{1}
	r2 := Register{2}

	if _, err := ecs.GetRegister[Register](world); err == nil {
		t.Errorf("got register from empty world")
		return
	}

	world.SaveRegister(r1)
	if register, err := ecs.GetRegister[Register](world); err != nil {
		t.Errorf("expected to get register but got error \"%s\"", err)
		return
	} else if register != r1 {
		t.Errorf("expected to get register but got invalid register %v", register)
		return
	}

	world.SaveRegister(r2)
	if register, err := ecs.GetRegister[Register](world); err != nil {
		t.Errorf("expected to get register but got error \"%s\"", err)
		return
	} else if register != r2 {
		t.Errorf("expected to get register but got invalid register %v", register)
		return
	}
}
