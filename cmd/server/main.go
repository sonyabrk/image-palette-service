package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sonyabrk/image-palette-service/internal/cache"
	"github.com/sonyabrk/image-palette-service/internal/config"
	"github.com/sonyabrk/image-palette-service/internal/processor"
	"github.com/sonyabrk/image-palette-service/internal/worker"
)

func main() {
	cfg := config.Load()
	memCache := cache.New(time.Duration(cfg.CacheTTL) * time.Second)
	proc := processor.New(cfg.MaxClusters)
	pool := worker.NewPool(cfg.Workers, proc, memCache)
	mux := http.NewServerMux()

	mux.HandleFunc("POST /analyze", handler.Analyze(pool))
	mux.HandleFunc("GET /metrics", handler.Heath())
	mux.HandleFunc("GET /metrics", handler.Metrics())

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("сервер запущен на порту %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ошибка запуска сервера %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("завершение работы...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("принудительное завершение: %v", err)
	}

	pool.Stop()

	log.Println("сервер остановлен")
}
