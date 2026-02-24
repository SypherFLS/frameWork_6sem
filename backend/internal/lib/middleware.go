package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var requestCounter int64

func generateRequestID() string {
	id := atomic.AddInt64(&requestCounter, 1)
	return fmt.Sprintf("req-%d-%d", time.Now().Unix(), id)
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.statusCode = code
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) StatusCode() int {
	if sr.statusCode == 0 {
		return 200
	}
	return sr.statusCode
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := generateRequestID()

		ctx := context.WithValue(r.Context(), "requestID", requestID)
		r = r.WithContext(ctx)

		start := time.Now()

		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode: 200,
		}

		LogOperation("REQUEST_START",
			fmt.Sprintf("Method: %s | Path: %s | RequestID: %s",
				r.Method, r.URL.Path, requestID))

		next(recorder, r)

		duration := time.Since(start)

		LogOperation("REQUEST_END",
			fmt.Sprintf("Method: %s | Path: %s | Status: %d | Duration: %v | RequestID: %s",
				r.Method, r.URL.Path, recorder.StatusCode(), duration, requestID))
	}
}

func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := "unknown"
				if id := r.Context().Value("requestID"); id != nil {
					requestID = id.(string)
				}
				_ = requestID
				LogError(&Err{
					Id: 500,
					Comment: fmt.Sprintf("panic recovered: %v", err),
					Identification: Iden{
						Num: 0,
						Process: "panic_recovery",
					},
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)

				errorResponse := MakeError(500, "internal server error", 0, "panic_recovery")
				json.NewEncoder(w).Encode(errorResponse)
			}
		}()

		next(w, r)
	}
}

func TimingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next(w, r)

		duration := time.Since(start)

		ctx := context.WithValue(r.Context(), "duration", duration)
		r = r.WithContext(ctx)

		requestID := "unknown"
		if id := r.Context().Value("requestID"); id != nil {
			requestID = id.(string)
		}

		LogOperation("TIMING",
			fmt.Sprintf("RequestID: %s | Duration: %v", requestID, duration))
	}
}

func ChainMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func GetRequestID(r *http.Request) string {
	if id := r.Context().Value("requestID"); id != nil {
		return id.(string)
	}
	return "unknown"
}
