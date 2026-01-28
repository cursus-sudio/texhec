package batcher

import "engine/services/ecs"

type System ecs.SystemRegister

type Batch struct {
	Steps   int
	Handler func(int)
}

func NewBatch(steps int, handler func(int)) Batch {
	return Batch{steps, handler}
}

type TaskFactory interface {
	AddBatch(Batch) TaskFactory
	Build() Task
}

type Task interface {
	Step()
	Progress() float32
}

type Service interface {
	Queue(Task)
	NewTask() TaskFactory
}
