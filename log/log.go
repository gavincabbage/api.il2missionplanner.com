// package log is adapted from the chi logrus middleware example found here:
// https://raw.githubusercontent.com/go-chi/chi/master/_examples/logging/main.go
package log

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

func Middleware(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&formatter{logger})
}

type formatter struct {
	Logger *logrus.Logger
}

func (f *formatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &entry{Logger: logrus.NewEntry(f.Logger)}

	entry.Logger = entry.Logger.WithFields(logrus.Fields{
		"ts":          time.Now().UTC().Format(time.RFC1123),
		"request_id":  middleware.GetReqID(r.Context()),
		"request_uri": r.RequestURI,
		"request_tls": r.TLS,
		"http_proto":  r.Proto,
		"http_method": r.Method,
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
	})

	entry.Logger.Infoln("request started")

	return entry
}

type entry struct {
	Logger logrus.FieldLogger
}

func (e *entry) Write(status, bytes int, elapsed time.Duration) {
	e.Logger = e.Logger.WithFields(logrus.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	e.Logger.Infoln("request complete")
}

func (e *entry) Panic(v interface{}, stack []byte) {
	e.Logger = e.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

// helper methods for manipulating the request-scoped logger from within request handlers

func RequestEntry(r *http.Request) logrus.FieldLogger {
	entry := middleware.GetLogEntry(r).(*entry)
	return entry.Logger
}

func RequestEntrySetField(r *http.Request, key string, value interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*entry); ok {
		entry.Logger = entry.Logger.WithField(key, value)
	}
}

func RequestEntrySetFields(r *http.Request, fields map[string]interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*entry); ok {
		entry.Logger = entry.Logger.WithFields(fields)
	}
}
