package errors

import "errors"

// job redis errors
var (
	ErrNotFound      = errors.New("jobId not found")
	ErrAlreadyQueued = errors.New("job already queued")
)
