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

// func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
// 	// Signal WaitGroup when this worker's loop finally ends
// 	defer wg.Done()
//
// 	for j := range jobs {
// 		fmt.Printf("Worker %d started job %d\n", id, j)
// 		time.Sleep(time.Second) // Simulate work
// 		fmt.Printf("Worker %d finished job %d\n", id, j)
// 	}
// }
//
// func test() {
// 	const numWorkers = 10
// 	const numJobs = 50
//
// 	jobs := make(chan int, numJobs)
// 	var wg sync.WaitGroup
//
// 	// 1. Start 10 workers
// 	for w := 1; w <= numWorkers; w++ {
// 		wg.Add(1)
// 		go worker(w, jobs, &wg)
// 	}
//
// 	// 2. Send jobs to the pool
// 	// This won't block because the channel is buffered
// 	for j := 1; j <= numJobs; j++ {
// 		jobs <- j
// 	}
//
// 	// 3. Close the channel to tell workers no more jobs are coming
// 	close(jobs)
//
// 	// 4. Wait for all workers to finish their current job and exit
// 	wg.Wait()
// 	fmt.Println("All work complete.")
// }

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
