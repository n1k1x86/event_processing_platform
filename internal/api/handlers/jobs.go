package handlers

import (
	"encoding/json"
	"event_processing_platform/internal/services/jobs"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type JobHandlerRequestBody struct {
	JobType string          `json:"job_type"`
	Payload json.RawMessage `json:"payload"`
}

type JobHandlerResponseBody struct {
	JobID string `json:"job_id"`
}

func JobsHandler(jobQueueManager *jobs.JobQueueManager, logger *zap.Logger) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		c.Set("Content-Type", ContentTypeJSON)

		data := c.Body()
		var reqBody JobHandlerRequestBody
		err := json.Unmarshal(data, &reqBody)
		if err != nil {
			handleClientError(c, logger, "error while unmarshaling body into json", err, http.StatusUnprocessableEntity)
			return nil
		}
		queue, err := jobQueueManager.GetQueue(jobs.JobType(reqBody.JobType))
		if err != nil {
			handleClientError(c, logger, "error while getting queue", err, http.StatusNotFound)
			return nil
		}

		job, jobID := jobs.NewJob(reqBody.Payload, jobs.JobType(reqBody.JobType))
		err = queue.Push(job)
		if err != nil {
			handleServerError(c, logger, "error while pushing job", err, http.StatusTooManyRequests)
			return nil
		}

		r := JobHandlerResponseBody{
			JobID: jobID.String(),
		}
		resp, err := json.Marshal(r)
		if err != nil {
			handleServerError(c, logger, "unexpected error while marshaling resulting body", err, http.StatusInternalServerError)
			return nil
		}
		c.Status(http.StatusAccepted).Write(resp)

		return nil
	}
}

type JobStatusHandlerResponseBody struct {
	Status string `json:"status"`
}

func JobsStatusHandler(jobsRuntimeManager *jobs.JobRuntimeManager, logger *zap.Logger) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		c.Set("Content-Type", ContentTypeJSON)
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			handleClientError(c, logger, "error while parsing id into uuid", err, http.StatusUnprocessableEntity)
			return nil
		}

		status, err := jobsRuntimeManager.GetJobStatus(id)
		if err != nil {
			handleClientError(c, logger, "error while getting job status", err, http.StatusNotFound)
			return nil
		}

		r := JobStatusHandlerResponseBody{
			Status: status,
		}

		data, err := json.Marshal(r)
		if err != nil {
			handleServerError(c, logger, "unexpected error while marshaling resp body", err, http.StatusInternalServerError)
			return nil
		}

		c.Status(http.StatusOK).Write(data)
		return nil
	}
}
