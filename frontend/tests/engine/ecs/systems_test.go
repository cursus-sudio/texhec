package ecs_test

import (
	"fmt"
	"frontend/src/engine/ecs"
	"frontend/src/engine/ecs/ecsargs"
	"testing"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type OrderedSystem struct {
	counter *int
	expect  int
}

func NewOrderedSystem(counter *int, expect int) OrderedSystem {
	return OrderedSystem{counter: counter, expect: expect}
}

func (system *OrderedSystem) Update(args ecs.Args) {
	*system.counter += 1
	if system.expect != *system.counter {
		panic(fmt.Sprintf("expected systems ordering. system should be %d but this is iterated as %d", system.expect, *system.counter))
	}
}

func TestSystemsOrder(t *testing.T) {
	b := ioc.NewBuilder()
	ecs.Package().Register(b)
	c := b.Build()
	world := ioc.Get[ecs.WorldFactory](c)()

	counter := 0
	type First struct{ OrderedSystem }
	first := &First{NewOrderedSystem(&counter, 1)}
	type Second struct{ OrderedSystem }
	second := &Second{NewOrderedSystem(&counter, 2)}
	type Third struct{ OrderedSystem }
	third := &Third{NewOrderedSystem(&counter, 3)}
	type Fourth struct{ OrderedSystem }
	fourth := &Fourth{NewOrderedSystem(&counter, 4)}

	world.LoadSystem(third, ecs.DrawSystem)
	world.LoadSystem(first, ecs.UpdateSystem)
	world.LoadSystem(fourth, ecs.DrawSystem)
	world.LoadSystem(second, ecs.UpdateSystem)

	deltaTime := ecsargs.NewDeltaTime(time.Duration(0))
	ecsArgs := ecs.NewArgs(deltaTime)

	if counter != 0 {
		t.Error("system updated counter when it shouldn't")
	}

	world.Update(ecsArgs)

	if counter != 4 {
		t.Error("system updated unexpected amout of times than expected")
	}
}
