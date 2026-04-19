package routes

import (
	"event_processing_platform/internal/api/handlers"
	"event_processing_platform/internal/api/server"
)

func SetHealthRoutes(s *server.HTTPServer) {
	app := s.GetApp()
	app.Get("/healthz", handlers.Healthz)
	app.Get("/readz", handlers.Readz)
}
