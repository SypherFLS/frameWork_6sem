package workerpool

import (
	"encoding/json"
	"framew/internal/db"
	"framew/internal/lib"
	"net/http"
	"strings"
	"sync"
)

type RequestContext struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Storage    *db.Conteiner
	Method     string
	Path       string
	ItemID     string
	Item       *db.Item
	Items      []db.Item
	Error      *lib.Err
	StatusCode int
	Response   interface{}
}

func NewRequestContext(w http.ResponseWriter, r *http.Request, storage *db.Conteiner) *RequestContext {
	return &RequestContext{
		Writer:  w,
		Request: r,
		Storage: storage,
		Method:  r.Method,
		Path:    r.URL.Path,
	}
}

func (ctx *RequestContext) WriteResponse() {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	if ctx.StatusCode == 0 {
		ctx.StatusCode = http.StatusOK
	}
	ctx.Writer.WriteHeader(ctx.StatusCode)
	if ctx.Response != nil {
		json.NewEncoder(ctx.Writer).Encode(ctx.Response)
	} else if ctx.Error != nil {
		json.NewEncoder(ctx.Writer).Encode(ctx.Error)
	} else if ctx.Items != nil {
		json.NewEncoder(ctx.Writer).Encode(ctx.Items)
	} else if ctx.Item != nil {
		json.NewEncoder(ctx.Writer).Encode(ctx.Item)
	}
}

func ParseRequestTask(data any, mu *sync.Mutex) (any, error) {
	ctx := data.(*RequestContext)

	if ctx.Method == "POST" {
		var item db.Item
		if err := json.NewDecoder(ctx.Request.Body).Decode(&item); err != nil {
			ctx.Error = lib.MakeError(400, "bad request", 1, "invalid json")
			ctx.StatusCode = ctx.Error.Id
			ctx.Response = ctx.Error
			return ctx, nil
		}
		ctx.Item = &item
	}

	if strings.HasPrefix(ctx.Path, "/api/items/") {
		ctx.ItemID = strings.TrimPrefix(ctx.Path, "/api/items/")
		if ctx.ItemID == "" || ctx.ItemID == ctx.Path {
			ctx.Error = lib.MakeError(400, "bad request", 1, "missing item id")
			ctx.StatusCode = ctx.Error.Id
			ctx.Response = ctx.Error
			return ctx, nil
		}
	}

	return ctx, nil
}

func ValidateItemTask(data any, mu *sync.Mutex) (any, error) {
	ctx := data.(*RequestContext)

	if ctx.Item == nil {
		return ctx, nil
	}

	validateErr := ctx.Item.Validate()
	if validateErr != nil && !validateErr.IsDone {
		ctx.Error = validateErr
		ctx.StatusCode = validateErr.Id
		ctx.Response = validateErr
		return ctx, nil
	}

	return ctx, nil
}

func GetAllItemsTask(data any, mu *sync.Mutex) (any, error) {
	ctx := data.(*RequestContext)

	mu.Lock()
	ctx.Items = ctx.Storage.GetAllItems()
	mu.Unlock()

	ctx.StatusCode = http.StatusOK
	ctx.Response = ctx.Items
	return ctx, nil
}

func GetItemByIdTask(data any, mu *sync.Mutex) (any, error) {
	ctx := data.(*RequestContext)

	mu.Lock()
	item, err := ctx.Storage.GetItemById(ctx.ItemID)
	mu.Unlock()

	if err != nil && !err.IsDone {
		ctx.Error = err
		ctx.StatusCode = err.Id
		ctx.Response = err
		return ctx, nil
	}

	ctx.Item = &item
	ctx.StatusCode = http.StatusOK
	ctx.Response = item
	return ctx, nil
}

func AddItemTask(data any, mu *sync.Mutex) (any, error) {
	ctx := data.(*RequestContext)

	if ctx.Item == nil {
		ctx.Error = lib.MakeError(400, "bad request", 1, "item is required")
		ctx.StatusCode = ctx.Error.Id
		ctx.Response = ctx.Error
		return ctx, nil
	}

	mu.Lock()
	addErr := ctx.Storage.AddNyItem(ctx.Item.Name, ctx.Item.Price)
	mu.Unlock()

	if addErr != nil && !addErr.IsDone {
		ctx.Error = addErr
		ctx.StatusCode = addErr.Id
		ctx.Response = addErr
		return ctx, nil
	}

	mu.Lock()
	items := ctx.Storage.GetAllItems()
	mu.Unlock()

	if len(items) == 0 {
		ctx.Error = lib.MakeError(500, "internal server error", 1, "failed to create item")
		ctx.StatusCode = ctx.Error.Id
		ctx.Response = ctx.Error
		return ctx, nil
	}

	createdItem := items[len(items)-1]
	ctx.Item = &createdItem
	ctx.StatusCode = http.StatusCreated
	ctx.Response = createdItem
	return ctx, nil
}

func WriteResponseTask(data any, mu *sync.Mutex) (any, error) {
	ctx := data.(*RequestContext)
	ctx.WriteResponse()
	return ctx, nil
}
