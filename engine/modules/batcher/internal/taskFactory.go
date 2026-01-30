package internal

import "engine/modules/batcher"

type taskFactory struct {
	batches                []Batch
	concurrentRoutinesUsed int
}

func NewTaskFactory(concurrentRoutineUsed int) batcher.TaskFactory {
	return &taskFactory{
		make([]Batch, 0),
		concurrentRoutineUsed,
	}
}

func (f *taskFactory) AddOrderedBatch(b batcher.Batch) batcher.TaskFactory {
	f.batches = append(f.batches, NewOrderedBatch(b))
	return f
}

func (f *taskFactory) AddConcurrentBatch(b batcher.Batch) batcher.TaskFactory {
	f.batches = append(f.batches, NewConcurrentBatch(b, f.concurrentRoutinesUsed))
	return f
}

func (f *taskFactory) Build() batcher.Task {
	return NewTask(f.batches)
}
