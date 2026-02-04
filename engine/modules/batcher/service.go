package batcher

import "engine/services/ecs"

type System ecs.SystemRegister

type Batch struct {
	Steps   int64
	Handler func(int64)
}

func NewBatch(steps int64, handler func(int64)) Batch {
	return Batch{steps, handler}
}

type TaskFactory interface {
	AddOrderedBatch(Batch) TaskFactory
	AddConcurrentBatch(Batch) TaskFactory
	Build() Task
}

type Task interface {
	Step()
	Progress() float32
}

type Service interface {
	NewTask() TaskFactory

	Queue(Task)
	// progress of first task in queue
	// when there is no tasks in queue -1 is returned
	Progress() float32
}
