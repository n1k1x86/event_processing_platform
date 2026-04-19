package jobs

import (
	"sync"

	"github.com/google/uuid"
)

type JobStorage struct {
	mu      sync.Mutex
	storage map[uuid.UUID]JobStatus
}

func NewJobStorager() *JobStorage {
	return &JobStorage{
		storage: make(map[uuid.UUID]JobStatus),
	}
}

func (j *JobStorage) Set(id uuid.UUID, status JobStatus) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.storage[id] = status
}

func (j *JobStorage) GetJobStatus(id uuid.UUID) (JobStatus, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	if status, ok := j.storage[id]; ok {
		return status, nil
	}
	return "", ErrJobNotFound
}
