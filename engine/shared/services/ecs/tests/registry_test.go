package ecs_test

import (
	"shared/services/ecs"
	"testing"
)

type register struct{ value int }
type Register struct{ *register }

func (r Register) Release() {
	r.value = 0
}

func TestRegistry(t *testing.T) {
	world := ecs.NewWorld()

	r1 := Register{&register{1}}
	r2 := Register{&register{2}}

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

	if r1.value != 0 {
		t.Errorf("register wasn't cleaned up properly on replace")
		return
	}

	if r2.value == 0 {
		t.Errorf("register was cleaned up prematurely")
		return
	}

	world.Release()
	if r2.value != 0 {
		t.Errorf("register wasn't cleaned up properly on clean up")
		return
	}
}
