package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type HandlerClientError struct {
	Details string `json:"details"`
}

type HandlerServerError struct {
	Details string `json:"details"`
}

func handleClientError(c fiber.Ctx, logger *zap.Logger, msg string, err error, status int) {
	logger.Info(msg, zap.Error(err))

	e := HandlerClientError{
		Details: msg,
	}
	data, err := json.Marshal(e)
	if err != nil {
		logger.Error("unexpected error", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(status)
	c.Write(data)
}

func handleServerError(c fiber.Ctx, logger *zap.Logger, msg string, err error, status int) {
	logger.Error(msg, zap.Error(err))

	e := HandlerServerError{
		Details: msg,
	}
	data, err := json.Marshal(e)
	if err != nil {
		logger.Error("unexpected error", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(status)
	c.Write(data)
}
