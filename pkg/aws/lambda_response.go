package aws

// LambdaResponse from the invocation.
type LambdaResponse struct {
	// StatusCode for the http response.
	StatusCode int `json:"statusCode" yaml:"statusCode"`
	// Headers to be appended.
	Headers map[string]string `json:"headers" yaml:"headers"`
	// Body to be written.
	Body Body `json:"body" yaml:"body"`
}
