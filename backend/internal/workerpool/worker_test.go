package workerpool

import (
	"sync"
	"testing"
	"time"
)

func TestNewWorker(t *testing.T) {
	taskChan := make(chan *TaskChain, 10)
	mu := &sync.Mutex{}

	worker := NewWorker(1, taskChan, mu)

	if worker == nil {
		t.Fatal("NewWorker вернул nil")
	}

	if worker.ID != 1 {
		t.Errorf("Ожидался ID 1, получено: %d", worker.ID)
	}

	if worker.taskChan != taskChan {
		t.Error("taskChan не установлен правильно")
	}

	if worker.Mu != mu {
		t.Error("Mutex не установлен правильно")
	}
}

func TestWorker_StartW(t *testing.T) {
	taskChan := make(chan *TaskChain, 10)
	mu := &sync.Mutex{}
	worker := NewWorker(1, taskChan, mu)
	wg := &sync.WaitGroup{}

	worker.StartW(wg)

	tasks := []*Task{
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			return "result", nil
		}, "data"),
	}

	chain := NewTaskChain(tasks)
	taskChan <- chain

	select {
	case result := <-chain.ResultChan:
		if result == nil {
			t.Error("Результат не должен быть nil")
		}
	case <-time.After(2 * time.Second):
		t.Error("Таймаут при ожидании результата")
	}

	close(taskChan)
	wg.Wait()
}

func TestWorker_StartW_MultipleTasks(t *testing.T) {
	taskChan := make(chan *TaskChain, 10)
	mu := &sync.Mutex{}
	worker := NewWorker(1, taskChan, mu)
	wg := &sync.WaitGroup{}

	worker.StartW(wg)

	completed := 0
	completedMutex := &sync.Mutex{}

	for i := 0; i < 3; i++ {
		tasks := []*Task{
			NewTask(func(data any, mu *sync.Mutex) (any, error) {
				time.Sleep(10 * time.Millisecond)
				return "result", nil
			}, "data"),
		}

		chain := NewTaskChain(tasks)
		taskChan <- chain

		go func() {
			<-chain.ResultChan
			completedMutex.Lock()
			completed++
			completedMutex.Unlock()
		}()
	}

	time.Sleep(200 * time.Millisecond)

	completedMutex.Lock()
	if completed != 3 {
		t.Errorf("Ожидалось завершение 3 задач, получено: %d", completed)
	}
	completedMutex.Unlock()

	close(taskChan)
	wg.Wait()
}

func TestWorker_StartW_ErrorHandling(t *testing.T) {
	taskChan := make(chan *TaskChain, 10)
	mu := &sync.Mutex{}
	worker := NewWorker(1, taskChan, mu)
	wg := &sync.WaitGroup{}

	worker.StartW(wg)

	tasks := []*Task{
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			return nil, &TestError{message: "test error"}
		}, "data"),
	}

	chain := NewTaskChain(tasks)
	taskChan <- chain

	time.Sleep(100 * time.Millisecond)

	if chain.Err == nil {
		t.Error("Ожидалась ошибка в цепочке")
	}

	close(taskChan)
	wg.Wait()
}

func TestWorker_StartW_ChannelClosed(t *testing.T) {
	taskChan := make(chan *TaskChain, 10)
	mu := &sync.Mutex{}
	worker := NewWorker(1, taskChan, mu)
	wg := &sync.WaitGroup{}

	worker.StartW(wg)

	close(taskChan)

	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Error("Таймаут при ожидании завершения воркера")
	}
}

func TestWorker_StartW_EmptyChain(t *testing.T) {
	taskChan := make(chan *TaskChain, 10)
	mu := &sync.Mutex{}
	worker := NewWorker(1, taskChan, mu)
	wg := &sync.WaitGroup{}

	worker.StartW(wg)

	chain := NewTaskChain([]*Task{})
	taskChan <- chain

	time.Sleep(100 * time.Millisecond)

	close(taskChan)
	wg.Wait()
}

func TestWorker_StartW_ChainWithMultipleTasks(t *testing.T) {
	taskChan := make(chan *TaskChain, 10)
	mu := &sync.Mutex{}
	worker := NewWorker(1, taskChan, mu)
	wg := &sync.WaitGroup{}

	worker.StartW(wg)

	callOrder := make([]int, 0)
	orderMutex := &sync.Mutex{}

	tasks := []*Task{
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			orderMutex.Lock()
			callOrder = append(callOrder, 1)
			orderMutex.Unlock()
			return "result1", nil
		}, "initial"),
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			orderMutex.Lock()
			callOrder = append(callOrder, 2)
			orderMutex.Unlock()
			return "result2", nil
		}, nil),
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			orderMutex.Lock()
			callOrder = append(callOrder, 3)
			orderMutex.Unlock()
			return "result3", nil
		}, nil),
	}

	chain := NewTaskChain(tasks)
	taskChan <- chain

	select {
	case result := <-chain.ResultChan:
		if result != "result3" {
			t.Errorf("Ожидался результат 'result3', получено: %v", result)
		}
	case <-time.After(2 * time.Second):
		t.Error("Таймаут при ожидании результата")
	}

	orderMutex.Lock()
	if len(callOrder) != 3 {
		t.Errorf("Ожидалось 3 вызова, получено: %d", len(callOrder))
	}
	if callOrder[0] != 1 || callOrder[1] != 2 || callOrder[2] != 3 {
		t.Errorf("Неправильный порядок вызовов: %v", callOrder)
	}
	orderMutex.Unlock()

	close(taskChan)
	wg.Wait()
}

type TestError struct {
	message string
}

func (e *TestError) Error() string {
	return e.message
}
