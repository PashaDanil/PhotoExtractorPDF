package errorx

import "errors"

// job errors
var (
	ErrNotFound         = errors.New("jobId not found")
	ErrAlreadyCompleted = errors.New("job already completed")
)
