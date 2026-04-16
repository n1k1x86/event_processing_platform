package handlers

import (
	"encoding/json"
	"event_processing_platform/internal/services/jobs"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type JobHandlerRequestBody struct {
	JobType string          `json:"job_type"`
	Payload json.RawMessage `json:"payload"`
}

type JobHandlerResponseBody struct {
	Details string `json:"details"`
}

func JobsHandler(jobQueueManager *jobs.JobQueueManager, logger *zap.Logger) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		c.Set("Content-Type", ContentTypeJSON)

		data := c.Body()
		var reqBody JobHandlerRequestBody
		err := json.Unmarshal(data, &reqBody)
		if err != nil {
			handleServerError(c, logger, "error while unmarshaling body into json", err)
			return err
		}
		queue, err := jobQueueManager.GetQueue(jobs.JobType(reqBody.JobType))
		if err != nil {
			handleClientError(c, logger, "error while getting queue", err)
			return err
		}

		job := jobs.NewJob(reqBody.Payload, uuid.New(), jobs.JobType(reqBody.JobType))
		err = queue.Push(job)
		if err != nil {
			handleServerError(c, logger, "error while pushing job", err)
			return err
		}

		return nil
	}
}
