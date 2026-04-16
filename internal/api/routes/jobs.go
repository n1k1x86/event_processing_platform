package routes

import (
	"event_processing_platform/internal/api/handlers"
	"event_processing_platform/internal/api/server"
	"event_processing_platform/internal/services/jobs"

	"go.uber.org/zap"
)

func SetJobsRoutes(s *server.HTTPServer, jobQueueManager *jobs.JobQueueManager, logger *zap.Logger) {
	app := s.GetApp()
	app.Post("/jobs", handlers.JobsHandler(jobQueueManager, logger))
}
