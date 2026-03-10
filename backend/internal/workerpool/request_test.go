package workerpool

import (
	"bytes"
	"encoding/json"
	"framew/internal/db"
	"framew/internal/lib"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestNewRequestContext(t *testing.T) {
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	req := httptest.NewRequest("GET", "/api/items", nil)
	w := httptest.NewRecorder()

	ctx := NewRequestContext(w, req, storage)

	if ctx == nil {
		t.Fatal("NewRequestContext вернул nil")
	}

	if ctx.Writer != w {
		t.Error("Writer не установлен правильно")
	}

	if ctx.Request != req {
		t.Error("Request не установлен правильно")
	}

	if ctx.Storage != storage {
		t.Error("Storage не установлен правильно")
	}

	if ctx.Method != "GET" {
		t.Errorf("Ожидался метод GET, получено: %s", ctx.Method)
	}

	if ctx.Path != "/api/items" {
		t.Errorf("Ожидался путь /api/items, получено: %s", ctx.Path)
	}
}

func TestRequestContext_WriteResponse(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/items", nil)
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	ctx := NewRequestContext(w, req, storage)

	ctx.StatusCode = http.StatusOK
	ctx.Items = []db.Item{{Id: "1", Name: "Test", Price: 10.0}}

	ctx.WriteResponse()

	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус код %d, получено: %d", http.StatusOK, w.Code)
	}

	if !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
		t.Error("Content-Type должен быть application/json")
	}

	var items []db.Item
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Errorf("Ошибка декодирования ответа: %v", err)
	}

	if len(items) != 1 {
		t.Errorf("Ожидался 1 элемент, получено: %d", len(items))
	}
}

func TestRequestContext_WriteResponse_WithError(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/items", nil)
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	ctx := NewRequestContext(w, req, storage)

	ctx.StatusCode = http.StatusBadRequest
	ctx.Error = lib.MakeError(400, "bad request", 1, "test error")

	ctx.WriteResponse()

	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус код %d, получено: %d", http.StatusBadRequest, w.Code)
	}
}

func TestParseRequestTask_GET(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/items", nil)
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	ctx := NewRequestContext(w, req, storage)
	mu := &sync.Mutex{}

	result, err := ParseRequestTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if resultCtx.Item != nil {
		t.Error("Item должен быть nil для GET запроса")
	}
}

func TestParseRequestTask_POST_ValidJSON(t *testing.T) {
	item := db.Item{Name: "Test Item", Price: 10.0}
	jsonData, _ := json.Marshal(item)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/items", bytes.NewBuffer(jsonData))
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	ctx := NewRequestContext(w, req, storage)
	mu := &sync.Mutex{}

	result, err := ParseRequestTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if resultCtx.Item == nil {
		t.Fatal("Item не должен быть nil")
	}

	if resultCtx.Item.Name != item.Name {
		t.Errorf("Ожидалось имя %s, получено: %s", item.Name, resultCtx.Item.Name)
	}
}

func TestParseRequestTask_POST_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/items", bytes.NewBufferString("invalid json"))
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	ctx := NewRequestContext(w, req, storage)
	mu := &sync.Mutex{}

	result, err := ParseRequestTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка (ошибка должна быть в контексте), получено: %v", err)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if resultCtx.Error == nil {
		t.Fatal("Ожидалась ошибка в контексте")
	}

	if resultCtx.StatusCode != 400 {
		t.Errorf("Ожидался статус код 400, получено: %d", resultCtx.StatusCode)
	}
}

func TestParseRequestTask_WithItemID(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/items/123", nil)
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	ctx := NewRequestContext(w, req, storage)
	mu := &sync.Mutex{}

	result, err := ParseRequestTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if resultCtx.ItemID != "123" {
		t.Errorf("Ожидался ItemID '123', получено: %s", resultCtx.ItemID)
	}
}

func TestParseRequestTask_InvalidItemID(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/items/", nil)
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	ctx := NewRequestContext(w, req, storage)
	mu := &sync.Mutex{}

	result, err := ParseRequestTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if resultCtx.Error == nil {
		t.Fatal("Ожидалась ошибка для пустого ItemID")
	}
}

func TestValidateItemTask_ValidItem(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/items", nil)
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	ctx := NewRequestContext(w, req, storage)
	ctx.Item = &db.Item{Name: "Test Item", Price: 10.0}
	mu := &sync.Mutex{}

	result, err := ValidateItemTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if resultCtx.Error != nil {
		t.Errorf("Ожидалась nil ошибка для валидного item, получено: %v", resultCtx.Error)
	}
}

func TestValidateItemTask_InvalidItem(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/items", nil)
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	ctx := NewRequestContext(w, req, storage)
	ctx.Item = &db.Item{Name: "", Price: 0.0}
	mu := &sync.Mutex{}

	result, err := ValidateItemTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if resultCtx.Error == nil {
		t.Fatal("Ожидалась ошибка для невалидного item")
	}
}

func TestValidateItemTask_NilItem(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/items", nil)
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	ctx := NewRequestContext(w, req, storage)
	mu := &sync.Mutex{}

	result, err := ValidateItemTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if resultCtx.Error != nil {
		t.Error("Ошибка не должна быть установлена для nil item")
	}
}

func TestGetAllItemsTask(t *testing.T) {
	storage := &db.Conteiner{
		Items: []db.Item{
			{Id: "1", Name: "Item 1", Price: 10.0},
			{Id: "2", Name: "Item 2", Price: 20.0},
		},
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/items", nil)
	ctx := NewRequestContext(w, req, storage)
	mu := &sync.Mutex{}

	result, err := GetAllItemsTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if len(resultCtx.Items) != 2 {
		t.Errorf("Ожидалось 2 элемента, получено: %d", len(resultCtx.Items))
	}

	if resultCtx.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус код %d, получено: %d", http.StatusOK, resultCtx.StatusCode)
	}
}

func TestGetItemByIdTask_Found(t *testing.T) {
	storage := &db.Conteiner{
		Items: []db.Item{
			{Id: "1", Name: "Item 1", Price: 10.0},
		},
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/items/1", nil)
	ctx := NewRequestContext(w, req, storage)
	ctx.ItemID = "1"
	mu := &sync.Mutex{}

	result, err := GetItemByIdTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if resultCtx.Item == nil {
		t.Fatal("Item не должен быть nil")
	}

	if resultCtx.Item.Id != "1" {
		t.Errorf("Ожидался ID '1', получено: %s", resultCtx.Item.Id)
	}

	if resultCtx.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус код %d, получено: %d", http.StatusOK, resultCtx.StatusCode)
	}
}

func TestGetItemByIdTask_NotFound(t *testing.T) {
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/items/999", nil)
	ctx := NewRequestContext(w, req, storage)
	ctx.ItemID = "999"
	mu := &sync.Mutex{}

	result, err := GetItemByIdTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if resultCtx.Error == nil {
		t.Fatal("Ожидалась ошибка для несуществующего item")
	}

	if resultCtx.StatusCode != 404 {
		t.Errorf("Ожидался статус код 404, получено: %d", resultCtx.StatusCode)
	}
}

func TestAddItemTask_Success(t *testing.T) {
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/items", nil)
	ctx := NewRequestContext(w, req, storage)
	ctx.Item = &db.Item{Name: "New Item", Price: 15.0}
	mu := &sync.Mutex{}

	result, err := AddItemTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if resultCtx.Error != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", resultCtx.Error)
	}

	if resultCtx.Item == nil {
		t.Fatal("Item не должен быть nil")
	}

	if resultCtx.StatusCode != http.StatusCreated {
		t.Errorf("Ожидался статус код %d, получено: %d", http.StatusCreated, resultCtx.StatusCode)
	}

	if len(storage.Items) != 1 {
		t.Errorf("Ожидался 1 элемент в хранилище, получено: %d", len(storage.Items))
	}
}

func TestAddItemTask_NilItem(t *testing.T) {
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/items", nil)
	ctx := NewRequestContext(w, req, storage)
	mu := &sync.Mutex{}

	result, err := AddItemTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if resultCtx.Error == nil {
		t.Fatal("Ожидалась ошибка для nil item")
	}

	if resultCtx.StatusCode != 400 {
		t.Errorf("Ожидался статус код 400, получено: %d", resultCtx.StatusCode)
	}
}

func TestAddItemTask_InvalidData(t *testing.T) {
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/items", nil)
	ctx := NewRequestContext(w, req, storage)
	ctx.Item = &db.Item{Name: "", Price: -1.0}
	mu := &sync.Mutex{}

	result, err := AddItemTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if resultCtx.Error == nil {
		t.Fatal("Ожидалась ошибка для невалидных данных")
	}
}

func TestWriteResponseTask(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/items", nil)
	storage := &db.Conteiner{Items: make([]db.Item, 0)}
	ctx := NewRequestContext(w, req, storage)
	ctx.StatusCode = http.StatusOK
	ctx.Items = []db.Item{{Id: "1", Name: "Test", Price: 10.0}}
	mu := &sync.Mutex{}

	result, err := WriteResponseTask(ctx, mu)

	if err != nil {
		t.Errorf("Ожидалась nil ошибка, получено: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус код %d, получено: %d", http.StatusOK, w.Code)
	}

	resultCtx, ok := result.(*RequestContext)
	if !ok {
		t.Fatal("Результат должен быть *RequestContext")
	}

	if resultCtx != ctx {
		t.Error("Результат должен быть тем же контекстом")
	}
}
