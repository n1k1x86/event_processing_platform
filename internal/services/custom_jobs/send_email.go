package custom_jobs

import (
	"context"
	"encoding/json"
	"event_processing_platform/internal/services/jobs"
	"fmt"

	"go.uber.org/zap"
)

var SendEmailJob jobs.JobType = "send_email"

type SendEmailHandler struct {
	logger *zap.Logger
}

func NewSendEmailHandler(logger *zap.Logger) *SendEmailHandler {
	return &SendEmailHandler{
		logger: logger,
	}
}

type SendEmailPayload struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SendEmailResult struct {
	Info string
}

func (s *SendEmailHandler) Execute(ctx context.Context, job *jobs.Job) (*jobs.Result, error) {
	var p SendEmailPayload
	err := json.Unmarshal(job.Payload, &p)
	if err != nil {
		s.logger.Error("error while unmarshaling payload job", zap.Error(err))
		return nil, err
	}
	r := SendEmailResult{
		Info: fmt.Sprintf("message was sended from: %s, to: %s, title: %s, message: %s\n", p.From, p.To, p.Title, p.Message),
	}
	return &jobs.Result{
		JobType: job.Type,
		Data:    r,
	}, nil
}

type SendEmailResultSink struct {
	logger *zap.Logger
}

func NewSendEmailResultSink(logger *zap.Logger) *SendEmailResultSink {
	return &SendEmailResultSink{
		logger: logger,
	}
}

func (s *SendEmailResultSink) PushResult(result *jobs.Result) error {
	s.logger.Info(result.Data.(SendEmailResult).Info, zap.String("job_type", string(result.JobType)))
	return nil
}
