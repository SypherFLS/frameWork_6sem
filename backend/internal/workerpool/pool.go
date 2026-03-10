package workerpool

import (
	"sync"
)

type WorkerPool struct {
	TaskChains []*TaskChain
	workers    int
	queue      chan *TaskChain
	wg         sync.WaitGroup
	mu         sync.Mutex
	running    bool
}

func NewPool(taskChains []*TaskChain, quantity int) *WorkerPool {
	return &WorkerPool{
		TaskChains: taskChains,
		workers:    quantity,
		queue:      make(chan *TaskChain, 100),
		running:    false,
	}
}

func NewRunningPool(quantity int) *WorkerPool {
	return &WorkerPool{
		workers: quantity,
		queue:   make(chan *TaskChain, 100),
		running: true,
	}
}

func (wp *WorkerPool) StartWP() {
	for i := 0; i < wp.workers; i++ {
		worker := NewWorker(i+1, wp.queue, &wp.mu)
		worker.StartW(&wp.wg)
	}

	if !wp.running {
		for _, taskChain := range wp.TaskChains {
			wp.queue <- taskChain
		}
		close(wp.queue)
		wp.wg.Wait()
	}
}

func (wp *WorkerPool) Start() {
	if wp.running {
		for i := 0; i < wp.workers; i++ {
			worker := NewWorker(i+1, wp.queue, &wp.mu)
			worker.StartW(&wp.wg)
		}
	}
}

func (wp *WorkerPool) Stop() {
	if wp.running {
		close(wp.queue)
		wp.wg.Wait()
		wp.running = false
	}
}

func (wp *WorkerPool) AddTaskChain(taskChain *TaskChain) {
	if wp.running {
		wp.queue <- taskChain
	}
}
