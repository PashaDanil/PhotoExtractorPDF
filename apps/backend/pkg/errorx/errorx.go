package errorx

import "errors"

// job redis errors
var (
	ErrNotFound         = errors.New("jobId not found")
	ErrAlreadyCompleted = errors.New("job already completed")
)

// job minio errors
var (
	ErrObjectNotFound = errors.New("object not found in storage")
)
