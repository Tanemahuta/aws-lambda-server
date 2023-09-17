package lambda

// Request to be sent to the lambda.
type Request struct {
	// Host in the request.
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
	// URI (complete) from the request.
	URI string `json:"uri,omitempty" yaml:"uri,omitempty"`
	// Headers of the request.
	Headers Headers `json:"headers,omitempty" yaml:"header,omitempty"`
	// Method from the request.
	Method string `json:"method,omitempty" yaml:"method,omitempty"`
	// Vars parsed from the path.
	Vars map[string]string `json:"vars,omitempty" yaml:"vars,omitempty"`
	// Body from the request
	Body []byte `json:"body,omitempty" yaml:"body,omitempty"`
}
