package workerpool

import (
	"sync"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	tasks := []*Task{
		NewTask(func(data any, mu *sync.Mutex) (any, error) { return "result", nil }, "data1"),
	}
	chains := []*TaskChain{
		NewTaskChain(tasks),
	}

	pool := NewPool(chains, 3)

	if pool == nil {
		t.Fatal("NewPool вернул nil")
	}

	if pool.workers != 3 {
		t.Errorf("Ожидалось 3 воркера, получено: %d", pool.workers)
	}

	if pool.running {
		t.Error("Пул не должен быть в режиме running при создании через NewPool")
	}

	if len(pool.TaskChains) != len(chains) {
		t.Errorf("Ожидалось %d цепочек, получено: %d", len(chains), len(pool.TaskChains))
	}
}

func TestNewRunningPool(t *testing.T) {
	pool := NewRunningPool(5)

	if pool == nil {
		t.Fatal("NewRunningPool вернул nil")
	}

	if pool.workers != 5 {
		t.Errorf("Ожидалось 5 воркеров, получено: %d", pool.workers)
	}

	if !pool.running {
		t.Error("Пул должен быть в режиме running при создании через NewRunningPool")
	}
}

func TestPool_StartWP(t *testing.T) {
	tasks := []*Task{
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			return "result1", nil
		}, "data1"),
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			return "result2", nil
		}, nil),
	}

	chains := []*TaskChain{
		NewTaskChain(tasks),
		NewTaskChain(tasks),
	}

	pool := NewPool(chains, 2)
	pool.StartWP()

	for _, chain := range chains {
		select {
		case result := <-chain.ResultChan:
			if result == nil {
				t.Error("Результат не должен быть nil")
			}
		case <-time.After(2 * time.Second):
			t.Error("Таймаут при ожидании результата")
		}
	}
}

func TestPool_Start(t *testing.T) {
	pool := NewRunningPool(2)
	pool.Start()

	time.Sleep(100 * time.Millisecond)

	tasks := []*Task{
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			return "result", nil
		}, "data"),
	}

	chain := NewTaskChain(tasks)
	pool.AddTaskChain(chain)

	select {
	case result := <-chain.ResultChan:
		if result == nil {
			t.Error("Результат не должен быть nil - воркеры должны быть запущены")
		}
	case <-time.After(2 * time.Second):
		t.Error("Таймаут - воркеры не обработали задачу")
	}

	pool.Stop()
}

func TestPool_Stop(t *testing.T) {
	pool := NewRunningPool(2)
	pool.Start()

	time.Sleep(100 * time.Millisecond)

	pool.Stop()

	if pool.running {
		t.Error("Пул должен быть остановлен")
	}
}

func TestPool_AddTaskChain(t *testing.T) {
	pool := NewRunningPool(2)
	pool.Start()

	tasks := []*Task{
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			return "result", nil
		}, "data"),
	}

	chain := NewTaskChain(tasks)
	pool.AddTaskChain(chain)

	select {
	case result := <-chain.ResultChan:
		if result == nil {
			t.Error("Результат не должен быть nil")
		}
	case <-time.After(2 * time.Second):
		t.Error("Таймаут при ожидании результата")
	}

	pool.Stop()
}

func TestPool_AddTaskChain_NotRunning(t *testing.T) {
	pool := NewPool([]*TaskChain{}, 2)

	tasks := []*Task{
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			return "result", nil
		}, "data"),
	}

	chain := NewTaskChain(tasks)
	pool.AddTaskChain(chain)

	select {
	case <-chain.ResultChan:
		t.Error("Результат не должен быть получен, так как пул не запущен")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestPool_ConcurrentTasks(t *testing.T) {
	pool := NewRunningPool(3)
	pool.Start()

	completed := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		tasks := []*Task{
			NewTask(func(data any, mu *sync.Mutex) (any, error) {
				time.Sleep(50 * time.Millisecond)
				return "result", nil
			}, "data"),
		}

		chain := NewTaskChain(tasks)
		pool.AddTaskChain(chain)

		go func() {
			<-chain.ResultChan
			completed <- true
		}()
	}

	for i := 0; i < 5; i++ {
		select {
		case <-completed:
		case <-time.After(3 * time.Second):
			t.Errorf("Таймаут при ожидании завершения задачи %d", i+1)
		}
	}

	pool.Stop()
}
