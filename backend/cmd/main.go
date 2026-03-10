package main

import (
	"encoding/json"
	"fmt"
	"framew/internal/db"
	"framew/internal/lib"
	_ "framew/internal/models"
	"framew/internal/workerpool"
	"net/http"
)

var Storage *db.Conteiner
var Pool *workerpool.WorkerPool

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := workerpool.NewRequestContext(w, r, Storage)

	var tasks []*workerpool.Task

	switch r.Method {
	case "GET":
		tasks = []*workerpool.Task{
			workerpool.NewTask(workerpool.ParseRequestTask, ctx),
			workerpool.NewTask(workerpool.GetAllItemsTask, nil),
			workerpool.NewTask(workerpool.WriteResponseTask, nil),
		}

	case "POST":
		tasks = []*workerpool.Task{
			workerpool.NewTask(workerpool.ParseRequestTask, ctx),
			workerpool.NewTask(workerpool.ValidateItemTask, nil),
			workerpool.NewTask(workerpool.AddItemTask, nil),
			workerpool.NewTask(workerpool.WriteResponseTask, nil),
		}

	default:
		errorResponse := lib.MakeError(405, "method not allowed", 1, "unsupported method")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorResponse.Id)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	chain := workerpool.NewTaskChain(tasks)
	Pool.AddTaskChain(chain)

	<-chain.ResultChan
}

func itemByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		errorResponse := lib.MakeError(405, "method not allowed", 1, "unsupported method")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorResponse.Id)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	ctx := workerpool.NewRequestContext(w, r, Storage)

	tasks := []*workerpool.Task{
		workerpool.NewTask(workerpool.ParseRequestTask, ctx),
		workerpool.NewTask(workerpool.GetItemByIdTask, nil),
		workerpool.NewTask(workerpool.WriteResponseTask, nil),
	}

	chain := workerpool.NewTaskChain(tasks)
	Pool.AddTaskChain(chain)

	<-chain.ResultChan
}

func main() {
	if err := lib.InitLogger(); err != nil {
		fmt.Printf("logers init fails: %v\n", err)
	}
	defer lib.CloseLogger()

	Storage = &db.Conteiner{
		Items: make([]db.Item, 0),
	}

	Pool = workerpool.NewRunningPool(5)
	Pool.Start()

	defer Pool.Stop()

	http.HandleFunc("/api/items", lib.ChainMiddleware(itemsHandler,
		lib.LoggingMiddleware,
		lib.RecoveryMiddleware,
		lib.TimingMiddleware,
	))

	http.HandleFunc("/api/items/", lib.ChainMiddleware(itemByIdHandler,
		lib.LoggingMiddleware,
		lib.RecoveryMiddleware,
		lib.TimingMiddleware,
	))

	http.ListenAndServe(":8080", nil)
}
