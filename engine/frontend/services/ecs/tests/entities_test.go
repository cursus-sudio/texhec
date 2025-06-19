package ecs_test

import (
	"frontend/services/ecs"
	"testing"

	"github.com/ogiusek/ioc/v2"
)

func TestEntities(t *testing.T) {
	b := ioc.NewBuilder()
	ecs.Package().Register(b)
	c := b.Build()
	world := ioc.Get[ecs.WorldFactory](c)()

	if len(world.GetEntities()) != 0 {
		t.Errorf("entity has unexpected entities")
	}

	entity := world.NewEntity()
	if !world.EntityExists(entity) {
		t.Errorf("added entity do not exists")
	}

	if len(world.GetEntities()) != 1 {
		t.Errorf("world do not have one more entity after creating new entity")
	}

	if world.GetEntities()[0] != entity {
		t.Errorf("world returns not new entity")
	}

	world.RemoveEntity(entity)

	if len(world.GetEntities()) != 0 {
		t.Errorf("world do not have one less entity upon deletion")
	}

	if world.EntityExists(entity) {
		t.Errorf("removed entity exists")
	}
}
