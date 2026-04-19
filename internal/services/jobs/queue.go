package jobs

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
)

type JobType string
type JobStatus string

var (
	JobPending           JobStatus = "Pending"
	JobProcessing        JobStatus = "Processing"
	JobFinishedWithError JobStatus = "Finished with error"
	JobFinished          JobStatus = "Finished"
)

type Job struct {
	Payload json.RawMessage `json:"payload"`
	ID      uuid.UUID       `json:"id"`
	Type    JobType         `json:"type"`
}

func NewJob(payload json.RawMessage, jobType JobType) (*Job, uuid.UUID) {
	jobID := uuid.New()
	return &Job{
		Payload: payload,
		ID:      jobID,
		Type:    jobType,
	}, jobID
}

type JobQueue struct {
	queue      []*Job
	mu         sync.Mutex
	size       int
	cond       *sync.Cond
	jobType    JobType
	closed     bool
	jobStorage *JobStorage
}

func InitJobQueue(size int, jobType JobType, jobStorage *JobStorage) *JobQueue {
	q := &JobQueue{
		queue:      make([]*Job, 0, size),
		size:       size,
		jobType:    jobType,
		jobStorage: jobStorage,
	}
	q.cond = sync.NewCond(&q.mu)

	return q
}

func (j *JobQueue) Push(job *Job) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if j.closed {
		return ErrQueueClosed
	}

	if len(j.queue) == j.size {
		return ErrQueueFull
	}

	j.queue = append(j.queue, job)
	j.jobStorage.Set(job.ID, JobPending)
	j.cond.Signal()

	return nil
}

func (j *JobQueue) Pop() (*Job, bool) {
	j.mu.Lock()
	defer j.mu.Unlock()

	for len(j.queue) == 0 {
		if j.closed {
			return nil, false
		}
		j.cond.Wait()
	}

	job := j.queue[0]
	j.queue[0] = nil

	j.queue = j.queue[1:]
	return job, true
}

func (j *JobQueue) Close() {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.closed {
		return
	}

	j.closed = true

	j.cond.Broadcast()
}
