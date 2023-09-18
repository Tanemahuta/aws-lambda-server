package config

// HTTP config.
type HTTP struct {
	// RequestTimeout for http.Server
	RequestTimeout Duration `json:"requestTimeout,omitempty" yaml:"requestTimeout,omitempty"`
	// EnableTraceparent injects traceparent.
	EnableTraceparent bool `json:"enableTraceparent,omitempty" yaml:"enableTraceparent,omitempty"`
}
