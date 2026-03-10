package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

	middleware := LoggingMiddleware(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус код %d, получено: %d", http.StatusOK, w.Code)
	}

	requestID := req.Context().Value("requestID")
	if requestID == nil {
		t.Error("RequestID должен быть установлен в контексте")
	}

	if requestIDStr, ok := requestID.(string); ok {
		if !strings.HasPrefix(requestIDStr, "req-") {
			t.Errorf("RequestID должен начинаться с 'req-', получено: %s", requestIDStr)
		}
	}
}

func TestRecoveryMiddleware_NoPanic(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

	middleware := RecoveryMiddleware(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус код %d, получено: %d", http.StatusOK, w.Code)
	}
}

func TestRecoveryMiddleware_WithPanic(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}

	middleware := RecoveryMiddleware(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус код %d, получено: %d", http.StatusInternalServerError, w.Code)
	}

	if !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
		t.Error("Content-Type должен быть application/json")
	}

	var errResponse Err
	if err := json.NewDecoder(w.Body).Decode(&errResponse); err != nil {
		t.Errorf("Ошибка декодирования ответа: %v", err)
	}

	if errResponse.Id != 500 {
		t.Errorf("Ожидался код ошибки 500, получено: %d", errResponse.Id)
	}
}

func TestRecoveryMiddleware_WithRequestID(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}

	loggingMiddleware := LoggingMiddleware(handler)
	recoveryMiddleware := RecoveryMiddleware(loggingMiddleware)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	recoveryMiddleware(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус код %d, получено: %d", http.StatusInternalServerError, w.Code)
	}
}

func TestTimingMiddleware(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	middleware := TimingMiddleware(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус код %d, получено: %d", http.StatusOK, w.Code)
	}

	duration := req.Context().Value("duration")
	if duration == nil {
		t.Error("Duration должен быть установлен в контексте")
	}
}

func TestChainMiddleware(t *testing.T) {
	callOrder := []string{}

	handler := func(w http.ResponseWriter, r *http.Request) {
		callOrder = append(callOrder, "handler")
		w.WriteHeader(http.StatusOK)
	}

	middleware1 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "middleware1")
			next(w, r)
		}
	}

	middleware2 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "middleware2")
			next(w, r)
		}
	}

	chained := ChainMiddleware(handler, middleware1, middleware2)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	chained(w, req)

	if len(callOrder) != 3 {
		t.Errorf("Ожидалось 3 вызова, получено: %d", len(callOrder))
	}

	if callOrder[0] != "middleware2" {
		t.Errorf("Ожидался первым middleware2, получено: %s", callOrder[0])
	}

	if callOrder[1] != "middleware1" {
		t.Errorf("Ожидался вторым middleware1, получено: %s", callOrder[1])
	}

	if callOrder[2] != "handler" {
		t.Errorf("Ожидался третьим handler, получено: %s", callOrder[2])
	}
}

func TestGetRequestID_WithID(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	middleware := LoggingMiddleware(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware(w, req)

	requestID := GetRequestID(req)

	if requestID == "" || requestID == "unknown" {
		t.Error("RequestID должен быть установлен")
	}
}

func TestGetRequestID_WithoutID(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	requestID := GetRequestID(req)

	if requestID != "unknown" {
		t.Errorf("Ожидался 'unknown', получено: %s", requestID)
	}
}

func TestStatusRecorder(t *testing.T) {
	w := httptest.NewRecorder()
	recorder := &statusRecorder{
		ResponseWriter: w,
		statusCode:     0,
	}

	recorder.WriteHeader(http.StatusNotFound)

	if recorder.StatusCode() != http.StatusNotFound {
		t.Errorf("Ожидался статус код %d, получено: %d", http.StatusNotFound, recorder.StatusCode())
	}

	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус код в ResponseWriter %d, получено: %d", http.StatusNotFound, w.Code)
	}
}

func TestStatusRecorder_DefaultStatusCode(t *testing.T) {
	w := httptest.NewRecorder()
	recorder := &statusRecorder{
		ResponseWriter: w,
		statusCode:     0,
	}

	if recorder.StatusCode() != http.StatusOK {
		t.Errorf("Ожидался статус код по умолчанию %d, получено: %d", http.StatusOK, recorder.StatusCode())
	}
}

func TestGenerateRequestID(t *testing.T) {
	id1 := generateRequestID()
	id2 := generateRequestID()

	if id1 == id2 {
		t.Error("RequestID должны быть уникальными")
	}

	if !strings.HasPrefix(id1, "req-") {
		t.Errorf("RequestID должен начинаться с 'req-', получено: %s", id1)
	}
}
