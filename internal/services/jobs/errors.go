package jobs

import "errors"

var (
	ErrQueueFull   = errors.New("queue is overflaw")
	ErrQueueClosed = errors.New("queue is closed")
)

var (
	ErrJobQueueAlreadyExists = errors.New("job queue already exists")
	ErrJobQueueNotFound      = errors.New("job queue not found")
)

var (
	ErrJobHandlerNotFound      = errors.New("job handler not found")
	ErrJobHandlerAlreadyExists = errors.New("job handler already registered")
)
