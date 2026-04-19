package jobs

import (
	"sync"

	"go.uber.org/zap"
)

type JobRuntimeManager struct {
	storage map[JobType]*JobRuntime
	logger  *zap.Logger
	mu      sync.RWMutex
}

func NewJobRuntimeManager(logger *zap.Logger) *JobRuntimeManager {
	return &JobRuntimeManager{
		storage: make(map[JobType]*JobRuntime),
		logger:  logger,
	}
}

func (j *JobRuntimeManager) RegisterRuntime(jobType JobType, jobRuntime *JobRuntime) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if _, ok := j.storage[jobType]; ok {
		return ErrJobRuntimeAlreadyRegistered
	}
	j.storage[jobType] = jobRuntime
	return nil
}

type JobSnapshot struct {
	jobRuntime *JobRuntime
	jobType    JobType
}

func (j *JobRuntimeManager) RunAll() {
	j.mu.RLock()
	runtimes := make([]*JobSnapshot, 0, len(j.storage))
	for t, r := range j.storage {
		runtimes = append(runtimes, &JobSnapshot{
			jobRuntime: r,
			jobType:    t,
		})
	}
	j.mu.RUnlock()

	for _, r := range runtimes {
		err := r.jobRuntime.Start()
		if err != nil {
			j.logger.Error("error while starting job runtime", zap.String("job_type", string(r.jobType)), zap.Error(err))
		}
	}
}

func (j *JobRuntimeManager) StopAll() {
	j.mu.Lock()
	runtimes := make([]*JobRuntime, 0, len(j.storage))
	for _, r := range j.storage {
		runtimes = append(runtimes, r)
	}
	j.mu.Unlock()

	for _, r := range runtimes {
		r.Stop()
	}
	j.logger.Info("all job runtimes were stopped")
}
