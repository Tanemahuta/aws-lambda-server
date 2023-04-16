package config

// Route configuration.
type Route struct {
	// Name for the route (optional).
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Host for the route (optional)
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
	// Methods for the route (optional)
	Methods []string `json:"methods,omitempty" yaml:"methods,omitempty"`
	// Path for the route (either that or PathPrefix).
	Path string `json:"path,omitempty" yaml:"path,omitempty" validate:"required_without=PathPrefix"`
	// PathPrefix for the route (either that or Path).
	PathPrefix string `json:"pathPrefix,omitempty" yaml:"pathPrefix,omitempty" validate:"required_without=Path"`
	// Headers for the route (optional).
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty" validate:"dive,keys,required,endkeys"`
	// HeadersRegexp for the route (optional).
	HeadersRegexp map[string]string `json:"headersRegexp" yaml:"headersRegexp" validate:"dive,keys,required,endkeys,required"` //nolint:lll // tags.
}
