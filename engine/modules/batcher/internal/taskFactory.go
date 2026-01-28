package internal

import "engine/modules/batcher"

type taskFactory struct {
	batches []batcher.Batch
}

func NewTaskFactory() batcher.TaskFactory {
	return &taskFactory{
		make([]batcher.Batch, 0),
	}
}

func (f *taskFactory) AddBatch(b batcher.Batch) batcher.TaskFactory {
	f.batches = append(f.batches, b)
	return f
}

func (f *taskFactory) Build() batcher.Task {
	return NewTask(f.batches)
}
