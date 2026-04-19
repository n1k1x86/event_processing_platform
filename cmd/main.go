package main

import (
	"context"
	"event_processing_platform/internal/api/routes"
	"event_processing_platform/internal/api/server"
	"event_processing_platform/internal/config"
	zaplogger "event_processing_platform/internal/logger"
	"event_processing_platform/internal/services/custom_jobs"
	"event_processing_platform/internal/services/jobs"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	parent, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zaplogger.NewLogger(cfg.App.Debug)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	jobRuntimeManager := jobs.NewJobRuntimeManager(logger)
	jobQueueManager := jobs.NewQueueManager()
	jobRegistry := jobs.NewRegistry()

	jobQueueSendEmail := jobs.InitJobQueue(cfg.Jobs.SendEmailJob.QueueSize, custom_jobs.SendEmailJob)
	err = jobQueueManager.RegisterJobQueue(custom_jobs.SendEmailJob, jobQueueSendEmail)
	if err != nil {
		zaplogger.ExitWithError(logger, "error while registering queue", zap.Error(err))
	}

	jobResultSinkSendEmail := custom_jobs.NewSendEmailResultSink(logger)
	jobHandlerSendEmail := custom_jobs.NewSendEmailHandler(logger)
	err = jobRegistry.Register(jobs.JobType(cfg.Jobs.SendEmailJob.JobType), jobHandlerSendEmail)
	if err != nil {
		zaplogger.ExitWithError(logger, "error while registering job handler", zap.Error(err))
	}

	jobSendEmailRuntime := jobs.NewJobRuntime(parent, cfg.Jobs.SendEmailJob.Workers, jobQueueManager,
		jobRegistry, custom_jobs.SendEmailJob, logger, jobResultSinkSendEmail)

	err = jobRuntimeManager.RegisterRuntime(custom_jobs.SendEmailJob, jobSendEmailRuntime)
	if err != nil {
		zaplogger.ExitWithError(logger, "error while registering runtime", zap.Error(err))
	}

	jobRuntimeManager.RunAll()

	app := server.NewHTTPServer(cfg.HTTPServer.ReadTimeout.Duration, cfg.HTTPServer.WriteTimeout.Duration)
	pprofServer := server.NewPprofServer(cfg.PprofServer.Addr, logger)

	routes.SetJobsRoutes(app, jobQueueManager, logger)
	routes.SetHealthRoutes(app)

	errChan := make(chan error, 1)

	go func() {
		errChan <- app.Start(cfg.HTTPServer.Addr)
	}()

	go func() {
		errChan <- pprofServer.Start()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
		ctxTimeout, cancel := context.WithTimeout(context.Background(), cfg.App.GracefulTimeout.Duration)
		defer cancel()

		err := app.Stop(ctxTimeout)
		if err != nil {
			logger.Error("server stopping failed", zap.Error(err))
		}

		err = pprofServer.Stop(ctxTimeout)
		if err != nil {
			logger.Error("server stopping failed", zap.Error(err))
		}

		jobRuntimeManager.StopAll()

	case err := <-errChan:
		if err != nil {
			zaplogger.ExitWithError(logger, "server failed", zap.Error(err))
		}
	}
}
