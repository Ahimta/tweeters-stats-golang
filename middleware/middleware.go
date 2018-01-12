package middleware

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// https://httpd.apache.org/docs/2.2/logs.html#combined + execution time.
const apacheFormatPattern = "%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\" %.3f\n"

// Apply blablabla
func Apply(handler http.Handler, writer io.Writer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			var err error

			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}

				log.Println("panic", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
			}
		}()

		clientIP := r.RemoteAddr
		if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
			clientIP = clientIP[:colon]
		}

		referer := r.Referer()
		if referer == "" {
			referer = "-"
		}

		userAgent := r.UserAgent()
		if userAgent == "" {
			userAgent = "-"
		}

		startTime := time.Now()
		w1 := &responseWriter{w, 200, 0}
		handler.ServeHTTP(w1, r)

		finishTime := time.Now()
		time := finishTime.UTC()
		elapsedTime := finishTime.Sub(startTime)
		timeFormatted := time.Format("02/Jan/2006 03:04:05")

		status := w1.status
		responseBytes := w1.responseBytes

		fmt.Fprintf(writer, apacheFormatPattern, clientIP, timeFormatted, r.Method,
			r.URL, r.Proto, status, responseBytes, referer, userAgent,
			elapsedTime.Seconds())
	})
}

type responseWriter struct {
	http.ResponseWriter
	status        int
	responseBytes int64
}

func (w *responseWriter) Write(p []byte) (int, error) {
	written, err := w.ResponseWriter.Write(p)
	w.responseBytes += int64(written)
	return written, err
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}