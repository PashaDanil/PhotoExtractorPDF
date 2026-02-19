package job

import "errors"

// job redis errors
var (
	ErrAlreadyExists    = errors.New("job already exists")
	ErrNotFound         = errors.New("jobId not found")
	ErrAlreadyCompleted = errors.New("job already completed")
	ErrAlreadyQueued    = errors.New("job already queued")
	ErrInvalidState     = errors.New("job in invalid state for this operation")
)

// job minio errors
var (
	ErrObjectNotFound     = errors.New("object not found in storage")
	ErrStorageUnavailable = errors.New("storage unavailable")
	ErrStorageForbidden   = errors.New("storage forbidden")
)
