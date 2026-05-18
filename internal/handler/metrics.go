package handler

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var metrics = struct {
	requestTotal atomic.Int64
	requestOK    atomic.Int64
	requestErr   atomic.Int64
	cacheHits    atomic.Int64
	cacheMisses  atomic.Int64
	startTime    time.Time
}{}

func init() {
	metrics.startTime = time.Now()
}

func Metrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Since(metrics.startTime).Seconds()

		w.Header().Set("Content-Type", "text/plain; version=0.0.4")

		fmt.Fprintf(w, "# HELP requests_total Общее количество запросов\n")
		fmt.Fprintf(w, "# TYPE requests_total counter\n")
		fmt.Fprintf(w, "requests_total %d\n\n", metrics.requestTotal.Load())

		fmt.Fprintf(w, "# HELP requests_ok Успешные запросы\n")
		fmt.Fprintf(w, "# TYPE requests_ok counter\n")
		fmt.Fprintf(w, "requests_ok %d\n\n", metrics.requestOK.Load())

		fmt.Fprintf(w, "# HELP requests_error Запросы с ошибкой\n")
		fmt.Fprintf(w, "# TYPE requests_error counter\n")
		fmt.Fprintf(w, "requests_error %d\n\n", metrics.requestErr.Load())

		fmt.Fprintf(w, "# HELP cache_hits Ответы из кэша\n")
		fmt.Fprintf(w, "# TYPE cache_hits counter\n")
		fmt.Fprintf(w, "cache_hits %d\n\n", metrics.cacheHits.Load())

		fmt.Fprintf(w, "# HELP cache_misses Новые вычисления\n")
		fmt.Fprintf(w, "# TYPE cache_misses counter\n")
		fmt.Fprintf(w, "cache_misses %d\n\n", metrics.cacheMisses.Load())

		fmt.Fprintf(w, "# HELP uptime_seconds Время работы сервиса в секундах\n")
		fmt.Fprintf(w, "# TYPE uptime_seconds gauge\n")
		fmt.Fprintf(w, "uptime_seconds %.2f\n\n", uptime)
	}
}

func TrackRequest(success bool) {
	metrics.requestTotal.Add(1)

	if success {
		metrics.requestOK.Add(1)
	} else {
		metrics.requestErr.Add(1)
	}
}

func TrackCache(hit bool) {
	if hit {
		metrics.cacheHits.Add(1)
	} else {
		metrics.cacheMisses.Add(1)
	}
}
