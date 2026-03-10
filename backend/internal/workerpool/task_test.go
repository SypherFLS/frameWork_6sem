package workerpool

import (
	"errors"
	"sync"
	"testing"
)

func TestNewTask(t *testing.T) {
	testFunc := func(data any, mu *sync.Mutex) (any, error) {
		return "result", nil
	}
	testData := "test data"

	task := NewTask(testFunc, testData)

	if task == nil {
		t.Fatal("NewTask вернул nil")
	}

	if task.F == nil {
		t.Error("Функция задачи не установлена")
	}

	if task.Data != testData {
		t.Errorf("Ожидались данные %v, получено %v", testData, task.Data)
	}
}

func TestNewTaskChain(t *testing.T) {
	tasks := []*Task{
		NewTask(func(data any, mu *sync.Mutex) (any, error) { return "task1", nil }, "data1"),
		NewTask(func(data any, mu *sync.Mutex) (any, error) { return "task2", nil }, "data2"),
	}

	chain := NewTaskChain(tasks)

	if chain == nil {
		t.Fatal("NewTaskChain вернул nil")
	}

	if len(chain.Tasks) != len(tasks) {
		t.Errorf("Ожидалось %d задач, получено %d", len(tasks), len(chain.Tasks))
	}

	if chain.ResultChan == nil {
		t.Error("ResultChan не инициализирован")
	}
}

func TestProcess_EmptyChain(t *testing.T) {
	chain := NewTaskChain([]*Task{})
	mu := &sync.Mutex{}

	err := Process(1, mu, chain)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка для пустой цепочки, получено: %v", err)
	}
}

func TestProcess_SuccessfulChain(t *testing.T) {
	callCount := 0
	mu := &sync.Mutex{}

	tasks := []*Task{
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			callCount++
			return "result1", nil
		}, "initial"),
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			callCount++
			if data != "result1" {
				t.Errorf("Ожидался результат 'result1', получено: %v", data)
			}
			return "result2", nil
		}, nil),
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			callCount++
			if data != "result2" {
				t.Errorf("Ожидался результат 'result2', получено: %v", data)
			}
			return "final", nil
		}, nil),
	}

	chain := NewTaskChain(tasks)
	err := Process(1, mu, chain)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	if callCount != 3 {
		t.Errorf("Ожидалось 3 вызова функций, получено: %d", callCount)
	}

	result := <-chain.ResultChan
	if result != "final" {
		t.Errorf("Ожидался результат 'final', получено: %v", result)
	}
}

func TestProcess_ErrorHandling(t *testing.T) {
	testError := errors.New("test error")
	mu := &sync.Mutex{}

	tasks := []*Task{
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			return "result1", nil
		}, "initial"),
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			return nil, testError
		}, nil),
		NewTask(func(data any, mu *sync.Mutex) (any, error) {
			t.Error("Третья задача не должна выполняться при ошибке")
			return "should not execute", nil
		}, nil),
	}

	chain := NewTaskChain(tasks)
	err := Process(1, mu, chain)

	if err != testError {
		t.Errorf("Ожидалась ошибка %v, получено: %v", testError, err)
	}

	if chain.Err != testError {
		t.Errorf("Ожидалась ошибка в цепочке %v, получено: %v", testError, chain.Err)
	}
}

func TestProcess_WithResultChan(t *testing.T) {
	mu := &sync.Mutex{}
	resultChan := make(chan any, 1)

	task := NewTask(func(data any, mu *sync.Mutex) (any, error) {
		return "task result", nil
	}, "initial")
	task.ResultChan = resultChan

	chain := NewTaskChain([]*Task{task})
	err := Process(1, mu, chain)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	select {
	case result := <-resultChan:
		if result != "task result" {
			t.Errorf("Ожидался результат 'task result', получено: %v", result)
		}
	default:
		t.Error("Результат не был отправлен в ResultChan задачи")
	}
}
