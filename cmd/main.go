// package main
//
// import (
//
//	"NoSQL/internal/http"
//
// )
//
//	func main() {
//		http.StartServer()
//	}
package main

import (
	"NoSQL/internal/http"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Перехват SIGINT и SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера
	server := http.NewServer()
	go func() {
		if err := server.Start(ctx); err != nil {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	// Блокируем до получения сигнала
	<-sigChan
	log.Println("🧹 Завершаем по сигналу...")
	cancel() // инициирует ctx.Done()
}
