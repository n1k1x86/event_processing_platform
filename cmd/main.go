package main

import (
	"context"
	"event_processing_platform/internal/api/routes"
	"event_processing_platform/internal/api/server"
	"event_processing_platform/internal/config"
	zaplogger "event_processing_platform/internal/logger"
	"log"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zaplogger.NewLogger(cfg.App.Debug)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	app := server.NewHTTPServer()
	routes.SetRoutes(app)
	errChan := make(chan error, 1)

	go func() {
		errChan <- app.Start(cfg.HTTPServer.Addr)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
		ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := app.Stop(ctxTimeout)
		if err != nil {
			zaplogger.ExitWithError(logger, "server stopping failed", zap.Error(err))
		}
	case err := <-errChan:
		if err != nil {
			zaplogger.ExitWithError(logger, "server failed", zap.Error(err))
		}
	}
}
