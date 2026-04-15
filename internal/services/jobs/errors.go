package jobs

import "errors"

var (
	ErrQueueFull   = errors.New("queue is overflaw")
	ErrQueueClosed = errors.New("queue is closed")
)
