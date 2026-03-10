package workerpool

import (
	"log"
	"sync"
)


type TaskChain struct {
	Tasks      []*Task
	ResultChan chan any 
	Err        error    
}


type Task struct {
	Err        error
	F          func(any, *sync.Mutex) (any, error) 
	Data       any
	ResultChan chan any 
}

func NewTask(f func(any, *sync.Mutex) (any, error), data any) *Task {
	return &Task{
		F:    f,
		Data: data,
	}
}


func NewTaskChain(tasks []*Task) *TaskChain {
	return &TaskChain{
		Tasks:      tasks,
		ResultChan: make(chan any, 1),
	}
}


func Process(workerId int, mutex *sync.Mutex, tc *TaskChain) error {
	if len(tc.Tasks) == 0 {
		log.Printf("worker %d: цепочка заданий пуста\n", workerId)
		return nil
	}

	log.Printf("worker %d: начало обработки цепочки из %d заданий\n", workerId, len(tc.Tasks))

	currentData := tc.Tasks[0].Data

	for i, task := range tc.Tasks {
		log.Printf("worker %d: обработка задачи %d/%d\n", workerId, i+1, len(tc.Tasks))

		result, err := task.F(currentData, mutex)
		if err != nil {
			task.Err = err
			tc.Err = err
			log.Printf("worker %d: ошибка в задаче %d: %v\n", workerId, i+1, err)
			return err
		}

		if task.ResultChan != nil {
			select {
			case task.ResultChan <- result:
			default:
			}
		}

		currentData = result
	}

	select {
	case tc.ResultChan <- currentData:
	default:
	}

	log.Printf("worker %d: цепочка заданий успешно обработана\n", workerId)
	return nil
}
