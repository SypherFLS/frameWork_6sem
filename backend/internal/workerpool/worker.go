package workerpool

import (
	"log"
	"sync"
)

type Worker struct {
	ID       int
	taskChan chan *TaskChain 
	Mu       *sync.Mutex
}

func NewWorker(id int, ch chan *TaskChain, mu *sync.Mutex) *Worker {
	return &Worker{
		ID:       id,
		taskChan: ch,
		Mu:       mu,
	}
}

func (w *Worker) StartW(wg *sync.WaitGroup) {
	log.Printf("worker %d: запущен\n", w.ID)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			taskChain, ok := <-w.taskChan
			if !ok {
				log.Printf("worker %d: канал закрыт, завершение работы\n", w.ID)
				break
			}
			if err := Process(w.ID, w.Mu, taskChain); err != nil {
				log.Printf("worker %d: ошибка при обработке цепочки: %v\n", w.ID, err)
			}
		}
		log.Printf("worker %d: завершил работу\n", w.ID)
	}()
}
