package main

import (
	"encoding/json"
	"fmt"
	"framew/internal/db"
	"framew/internal/lib"
	_ "framew/internal/models"
	"net/http"
	"strings"
)

var Storage *db.Conteiner

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		items := Storage.GetAllItems()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(items)

	case "POST":
		var item db.Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			errorResponse := lib.MakeError(400, "bad request", 1, "invalid json")
			w.WriteHeader(errorResponse.Id)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		validateErr := item.Validate()
		if validateErr != nil && !validateErr.IsDone {
			w.WriteHeader(validateErr.Id)
			json.NewEncoder(w).Encode(validateErr)
			return
		}

		addErr := Storage.AddNyItem(item.Name, item.Price)
		if addErr != nil && !addErr.IsDone {
			w.WriteHeader(addErr.Id)
			json.NewEncoder(w).Encode(addErr)
			return
		}

		items := Storage.GetAllItems()
		if len(items) == 0 {
			errorResponse := lib.MakeError(500, "internal server error", 1, "failed to create item")
			w.WriteHeader(errorResponse.Id)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		createdItem := items[len(items)-1]
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdItem)

	default:
		errorResponse := lib.MakeError(405, "method not allowed", 1, "unsupported method")
		w.WriteHeader(errorResponse.Id)
		json.NewEncoder(w).Encode(errorResponse)
	}
}

func itemByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		errorResponse := lib.MakeError(405, "method not allowed", 1, "unsupported method")
		w.WriteHeader(errorResponse.Id)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/items/")
	if path == "" || path == r.URL.Path {
		errorResponse := lib.MakeError(400, "bad request", 1, "missing item id")
		w.WriteHeader(errorResponse.Id)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	item, err := Storage.GetItemById(path)
	if err != nil && !err.IsDone {
		w.WriteHeader(err.Id)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func main() {
	if err := lib.InitLogger(); err != nil {
		fmt.Printf("logers init fails: %v\n", err)
	}
	defer lib.CloseLogger()

	Storage = &db.Conteiner{
		Items: make([]db.Item, 0),
	}

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
