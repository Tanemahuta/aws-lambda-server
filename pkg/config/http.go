package config

// HTTP config.
type HTTP struct {
	// ReadTimeout for http.Server
	ReadTimeout Duration `json:"readTimeout,omitempty" yaml:"readTimeout,omitempty"`
	// ReadHeaderTimeout for http.Server
	ReadHeaderTimeout Duration `json:"readHeaderTimeout,omitempty" yaml:"readHeaderTimeout,omitempty"`
	// WriteTimeout for http.Server
	WriteTimeout Duration `json:"writeTimeout,omitempty" yaml:"writeTimeout,omitempty"`
	// EnableTraceparent injects traceparent.
	EnableTraceparent bool `json:"enableTraceparent,omitempty" yaml:"enableTraceparent,omitempty"`
}
