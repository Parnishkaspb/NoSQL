package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"NoSQL/internal/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	server := http.NewServer()
	go func() {
		if err := server.Start(ctx); err != nil {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	<-sigChan
	log.Println("Завершаем по сигналу...")
	cancel()
}
