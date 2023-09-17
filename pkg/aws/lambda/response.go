package lambda

// Response from the invocation.
type Response struct {
	// StatusCode for the http response.
	StatusCode int `json:"statusCode" yaml:"statusCode"`
	// Headers to be appended.
	Headers Headers `json:"headers" yaml:"headers"`
	// Body to be written.
	Body Body `json:"body" yaml:"body"`
}
