package jobs

import (
	"context"
	"sync"
)

type JobHandler interface {
	Execute(ctx context.Context, job *Job) (*Result, error)
}

type JobRegistry struct {
	handlers map[JobType]JobHandler
	mu       sync.RWMutex
}

func NewRegistry() *JobRegistry {
	return &JobRegistry{
		handlers: make(map[JobType]JobHandler),
	}
}

func (j *JobRegistry) Handle(ctx context.Context, job *Job) (*Result, error) {
	j.mu.RLock()
	handler, ok := j.handlers[job.Type]
	j.mu.RUnlock()

	if !ok {
		return nil, ErrJobHandlerNotFound
	}
	return handler.Execute(ctx, job)
}

func (j *JobRegistry) Register(jobType JobType, handler JobHandler) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if _, ok := j.handlers[jobType]; ok {
		return ErrJobHandlerAlreadyExists
	}
	j.handlers[jobType] = handler
	return nil
}

func (j *JobRegistry) Unregister(jobType JobType) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if _, ok := j.handlers[jobType]; !ok {
		return ErrJobHandlerNotFound
	}
	delete(j.handlers, jobType)
	return nil
}
