package main

import (
	"assignment-backend/internal/application"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	app := application.New()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	errCh := make(chan error, 1)

	go func() {
		errCh <- app.Start()
	}()

	log.Println("product-service started")

	select {
	case <-ctx.Done():
		log.Println("product-service: received shutdown signal")

	case err := <-errCh:
		if err != nil {
			log.Printf("product-service: app exited with error: %v", err)
		} else {
			log.Println("product-service: app stopped unexpectedly")
		}
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		log.Printf("product-service: graceful shutdown failed: %v", err)
	}
}
