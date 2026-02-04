package internal

import (
	"engine/modules/batcher"
	"sync"
)

type Batch interface {
	Step() (finished bool)
	Steps() int
}

type orderedBatch struct {
	blueprint batcher.Batch
	index     int
}

func (b *orderedBatch) Step() (finished bool) {
	if b.blueprint.Steps == b.index {
		return true
	}
	b.blueprint.Handler(b.index)
	b.index++
	return false
}

func (b *orderedBatch) Steps() int {
	return b.blueprint.Steps
}

func NewOrderedBatch(b batcher.Batch) Batch {
	return &orderedBatch{
		blueprint: b,
		index:     0,
	}
}

//

type concurrentBatch struct {
	blueprint batcher.Batch
	index     int

	concurrentRoutinesUsed int
	wg                     *sync.WaitGroup
	todo                   chan int
}

func (b *concurrentBatch) Step() (finished bool) {
	b.todo <- b.index
	b.index++

	if b.blueprint.Steps == b.index {
		close(b.todo)
		b.wg.Wait()
		return true
	}

	if b.index != 1 {
		return false
	}

	// initialize workers
	for range b.concurrentRoutinesUsed {
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			for i := range b.todo {
				b.blueprint.Handler(i)
			}
		}()
	}
	return false
}

func (b *concurrentBatch) Steps() int {
	return b.blueprint.Steps
}

func NewConcurrentBatch(b batcher.Batch, concurrentRoutinesUsed int) Batch {
	batch := &concurrentBatch{
		blueprint: b,
		index:     0,

		concurrentRoutinesUsed: concurrentRoutinesUsed,
		wg:                     &sync.WaitGroup{},
		todo:                   make(chan int, concurrentRoutinesUsed),
	}
	return batch
}
