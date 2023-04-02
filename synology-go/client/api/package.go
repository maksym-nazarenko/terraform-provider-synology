package api

// Request defines a contract for all Request implementations.
type Request interface {
	ErrorDescriber

	// APIPath defines an API base path for this request.
	// Example: /webapi/entry.cgi
	APIPath() string

	// APIName defines an API to be called.
	// Example: SYNO.FileStation.Info
	APIName() string

	// APIVersion defines version number to use during call.
	APIVersion() int

	// APIMethod defines a method to use for this request.
	// Don't be confused with REST HTTP methods like GET/POST/PUT/etc
	// this method is RPC-style string which could be anything.
	// Example: list, get, rename
	APIMethod() string

	// RequestParams is a request arguments mapping.
	// It is request-specific.
	RequestParams() RequestParams

	// todo(maksym): change the description
	// ParseResponse parses Data field from Synology API response
	// into concrete expected response type.
	// It is responsibility of the caller to make a type assertion.
	NewResponseInstance() Response
}

// Response defines an interface for all responses from Synology API.
type Response interface{}

// RequestParams represents a mapping of name=>value request parameters.
type RequestParams map[string]string

// GenericResponse is a concrete Response implementation.
// It is a generic struct with common to all responses fields.
type GenericResponse struct {
	Success bool
	Data    interface{}
	Error   SynologyError
}
