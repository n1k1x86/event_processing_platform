package main

import (
	"context"
	"event_processing_platform/internal/api/routes"
	"event_processing_platform/internal/api/server"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	app := server.NewHTTPServer()
	routes.SetRoutes(app)
	errChan := make(chan error, 1)

	go func() {
		errChan <- app.Start("localhost:8000")
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
		ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := app.Stop(ctxTimeout)
		if err != nil {
			log.Fatal(err)
		}
	case err := <-errChan:
		if err != nil {
			log.Fatalf("server failed: %s\n", err.Error())
		}
	}
}
