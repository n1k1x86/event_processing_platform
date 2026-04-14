package main

import (
	"context"
	"event_processing_platform/internal/api/routes"
	"event_processing_platform/internal/api/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	app := server.NewHTTPServer()

	routes.SetRouters(app)

	app.Start("localhost:8000")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := app.Stop(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
