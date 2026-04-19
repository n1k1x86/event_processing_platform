package jobs

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

type JobRuntime struct {
	workers         int
	jobQueueManager *JobQueueManager
	jobRegistry     *JobRegistry
	jobStorage      *JobStorage
	resultSink      ResultSink
	logger          *zap.Logger
	parent          context.Context
	ctx             context.Context
	cancel          context.CancelFunc
	jobType         JobType
	wg              sync.WaitGroup
	once            sync.Once
	started         bool
	mu              sync.Mutex
}

func (j *JobRuntime) Start() error {
	j.mu.Lock()
	if j.started {
		j.mu.Unlock()
		return ErrJobRuntimeAlreadyStarted
	}

	select {
	case <-j.parent.Done():
		j.mu.Unlock()
		return ErrParentContextDone
	default:
	}

	queue, err := j.jobQueueManager.GetQueue(j.jobType)
	if err != nil {
		j.mu.Unlock()
		j.logger.Error("error while getting job queue", zap.Error(err))
		return err
	}

	ctx, cancel := context.WithCancel(j.parent)

	j.ctx = ctx
	j.cancel = cancel

	j.wg.Add(j.workers)

	j.started = true
	j.mu.Unlock()

	for range j.workers {
		go func() {
			defer j.wg.Done()
			defer func() {
				r := recover()
				if r != nil {
					j.logger.Error("panic was recovered", zap.Any("panic", r))
				}
			}()

			for {
				select {
				case <-j.ctx.Done():
					j.logger.Info("exit worker context done")
					return
				default:
				}

				job, ok := queue.Pop()
				if ok {
					j.jobStorage.Set(job.ID, JobProcessing)
					result, err := j.jobRegistry.Handle(j.ctx, job)
					if err != nil {
						j.logger.Error("error while handling job", zap.Error(err))
						j.jobStorage.Set(job.ID, JobFinishedWithError)
						continue
					} else if result != nil {
						err = j.resultSink.PushResult(result)
						if err != nil {
							j.logger.Error("error while pushing result", zap.Error(err))
							j.jobStorage.Set(job.ID, JobFinishedWithError)
							continue
						}
					}
					j.jobStorage.Set(job.ID, JobFinished)
				} else {
					return
				}
			}
		}()
	}

	return nil
}

func (j *JobRuntime) Stop() {
	j.mu.Lock()
	if !j.started {
		j.mu.Unlock()
		return
	}
	j.mu.Unlock()
	j.once.Do(j.stop)
}

func (j *JobRuntime) stop() {
	if j.ctx == nil || j.cancel == nil {
		return
	}
	j.cancel()
	err := j.jobQueueManager.CloseQueue(j.jobType)
	if err != nil {
		j.logger.Error("error while closing job queue", zap.Error(err))
	}

	j.wg.Wait()
}

func NewJobRuntime(parent context.Context, workers int, jobQueueManager *JobQueueManager,
	jobRegistry *JobRegistry, jobType JobType, logger *zap.Logger, resultSink ResultSink, jobStorage *JobStorage) *JobRuntime {
	return &JobRuntime{
		workers:         workers,
		jobQueueManager: jobQueueManager,
		jobRegistry:     jobRegistry,
		parent:          parent,
		jobType:         jobType,
		logger:          logger,
		wg:              sync.WaitGroup{},
		resultSink:      resultSink,
		jobStorage:      jobStorage,
	}
}
