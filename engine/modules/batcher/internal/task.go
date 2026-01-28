package internal

import "engine/modules/batcher"

type task struct {
	allSteps    int
	currentStep int

	currentBatch     int
	currentBatchStep int
	batches          []batcher.Batch
}

func NewTask(batches []batcher.Batch) batcher.Task {
	allSteps := 0
	for _, batch := range batches {
		allSteps += batch.Steps
	}
	return &task{
		allSteps,
		0,

		0,
		0,
		batches,
	}
}

func (t *task) Step() {
	batch := t.batches[t.currentBatch]
	batch.Handler(t.currentBatchStep)

	t.currentStep++
	t.currentBatchStep++
	if t.currentBatchStep != batch.Steps {
		return
	}
	t.currentBatch++
	t.currentBatchStep = 0
}

func (t *task) Progress() float32 {
	return float32(t.currentStep) / float32(t.allSteps)
}
