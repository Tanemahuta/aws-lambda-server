package aws

// LambdaRequest to be sent to the lambda.
type LambdaRequest struct {
	// Host in the request.
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
	// Headers of the request.
	Headers Headers `json:"headers,omitempty" yaml:"header,omitempty"`
	// Method from the request.
	Method string `json:"method,omitempty" yaml:"method,omitempty"`
	// URI from the request.
	URI string `json:"uri" yaml:"uri"`
	// Vars parsed from the path.
	Vars map[string]string `json:"vars" yaml:"vars"`
	// Body from the request
	Body []byte `json:"body" yaml:"body"`
}
