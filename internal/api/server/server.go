package server

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
)

type HTTPServer struct {
	app *fiber.App
}

func NewHTTPServer(readTimeout, writeTimeout time.Duration) *HTTPServer {
	return &HTTPServer{
		app: fiber.New(
			fiber.Config{
				ReadTimeout:  readTimeout,
				WriteTimeout: writeTimeout,
			},
		),
	}
}

func (h *HTTPServer) GetApp() *fiber.App {
	return h.app
}

func (h *HTTPServer) Start(addr string) error {
	return h.app.Listen(addr)
}

func (h *HTTPServer) Stop(ctx context.Context) error {
	err := h.app.ShutdownWithContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
