package forwardtest

import "errors"

var (
	// ErrInvalidStatus is returned when the status is invalid.
	ErrInvalidStatus = errors.New("invalid status")
)

// Status represents the status of a forwardtest.
type Status string

const (
	// StatusReady indicates that the forwardtest is ready to start.
	StatusReady Status = "ready"
	// StatusRunning indicates that the forwardtest is currently running.
	StatusRunning Status = "running"
	// StatusFinished indicates that the forwardtest has finished.
	StatusFinished Status = "finished"
)

// String returns the string representation of the status.
func (s Status) String() string {
	return string(s)
}

// Validate checks if the status is valid.
func (s Status) Validate() error {
	switch s {
	case StatusReady, StatusRunning, StatusFinished:
		return nil
	default:
		return ErrInvalidStatus
	}
}
