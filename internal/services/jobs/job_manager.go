package jobs

import "sync"

type JobManager struct {
	storage map[JobType]*JobQueue
	mu      sync.RWMutex
}

func (j *JobManager) RegisterJobQueue(jobType JobType, queue *JobQueue) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if _, ok := j.storage[jobType]; ok {
		return ErrJobQueueAlreadyExists
	}

	j.storage[jobType] = queue

	return nil
}

func (j *JobManager) GetQueue(jobType JobType) (*JobQueue, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()

	if q, ok := j.storage[jobType]; ok {
		return q, nil
	}

	return nil, ErrJobQueueNotFound
}

func (j *JobManager) CloseQueue(jobType JobType) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if q, ok := j.storage[jobType]; ok {
		q.Close()
		return nil
	}
	return ErrJobQueueNotFound
}

func (j *JobManager) CloseAll() {
	j.mu.Lock()
	defer j.mu.Unlock()

	for _, q := range j.storage {
		q.Close()
	}
}

func NewManager() *JobManager {
	return &JobManager{
		storage: make(map[JobType]*JobQueue),
	}
}
