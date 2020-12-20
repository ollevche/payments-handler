package rest

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
)

type failResponseBody struct {
	ErrorMessage string `json:"error_message"`
}

func respondJSON(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Err(err).Msg("Failed to respond with JSON")
		return
	}
}

func HandlePanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error().Msgf("Panic happened in http handler: %v", r)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// TODO: implement http.Hijacker
type deferredResponseWriter struct {
	bytes.Buffer
	original       http.ResponseWriter
	responseStatus int
}

func (w *deferredResponseWriter) Header() http.Header {
	return w.original.Header()
}

func (w *deferredResponseWriter) WriteHeader(status int) {
	if w.responseStatus == 0 {
		w.responseStatus = status
	}
}

func (w *deferredResponseWriter) writeToOriginal() error {
	if _, err := w.original.Write(w.Buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var startedAt = time.Now()

		var dw = &deferredResponseWriter{
			original: w,
		}

		next.ServeHTTP(dw, r)

		var timeTook = time.Now().Sub(startedAt).Microseconds()

		dw.Header().Set("X-Response-Time", strconv.FormatInt(timeTook, 10))

		if err := dw.writeToOriginal(); err != nil {
			log.Error().Err(err).Msg("Failed to write to original http.ResponseWriter from deferred")
		}
	})
}

func GetServerNameMiddleware(serverName string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Server-Name", serverName)
			next.ServeHTTP(w, r)
		})
	}
}
