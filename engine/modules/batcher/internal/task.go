package internal

import (
	"engine/modules/batcher"
)

type task struct {
	allSteps    int64
	currentStep int

	currentBatch int
	batches      []Batch
}

func NewTask(batches []Batch) batcher.Task {
	var allSteps int64
	for _, batch := range batches {
		allSteps += batch.Steps()
	}
	return &task{
		allSteps,
		0,

		0,
		batches,
	}
}

func (t *task) Step() {
	batch := t.batches[t.currentBatch]
	if finished := batch.Step(); finished {
		t.currentBatch++
	}
	t.currentStep++
}

func (t *task) Progress() float32 {
	return float32(t.currentStep) / float32(t.allSteps)
}
