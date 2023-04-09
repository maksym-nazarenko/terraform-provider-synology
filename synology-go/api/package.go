package api

// Request defines a contract for all Request implementations.
type Request interface {
}

// Response defines an interface for all responses from Synology API.
type Response interface {
	ErrorDescriber

	GetError() SynologyError
	SetError(SynologyError)
	Success() bool
}

// GenericResponse is a concrete Response implementation.
// It is a generic struct with common to all responses fields.
type GenericResponse struct {
	Success bool
	Data    interface{}
	Error   SynologyError
}
