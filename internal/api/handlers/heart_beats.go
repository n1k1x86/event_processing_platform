package handlers

import "github.com/gofiber/fiber/v3"

func Healthz(c fiber.Ctx) error {
	return c.SendString("OK")
}

func Readz(c fiber.Ctx) error {
	// тут будет проверка всех сервисов (Postgres, Redis, Kafka, Mongo)
	return c.SendString("READY")
}
