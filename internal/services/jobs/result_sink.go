package jobs

type ResultSink interface {
	PushResult(result *Result) error
}

type Result struct {
	JobType JobType
	Data    any
}

func NewResult(jobType JobType, data any) *Result {
	return &Result{
		JobType: jobType,
		Data:    data,
	}
}
