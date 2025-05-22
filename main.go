package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	storage := NewMemoryStorage()

	job := NewJob(storage)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Запускаем джоб для отправки нотификации
	wg.Add(1)
	go func() {
		defer wg.Done()
		job.Start(ctx)
	}()

	handler := NewEventHandler(storage)

	mux := http.NewServeMux()
	mux.HandleFunc("/events", handler.AddEvent)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Запускаем http сервер
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil {
			log.Println("Can't start http server", err)
		}
	}()

	// канал для прослушивания сигналов от ОС
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Ждем до получения сингала
	<-stop
	log.Println("Shutting down...")

	// Отключаем джоб
	cancel()

	// Отключаем http сервер
	ctxShutdown, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Printf("Server Shutdown failed: %v", err)
	}

	wg.Wait()
	log.Println("Gracefully shutdown successfully")
}
