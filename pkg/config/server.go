package config

// Server for the handler.
type Server struct {
	// DisableValidation disables the validation of the function using a dry-run invocation.
	DisableValidation bool `json:"disableValidation,omitempty" yaml:"disableValidation,omitempty"`
	// HTTP configuration
	HTTP HTTP `json:"http,omitempty" yaml:"http,omitempty"`
	// AWS configuration
	AWS *AWS `json:"aws,omitempty" yaml:"aws,omitempty"`
	// Functions providing the Function.
	Functions []Function `json:"functions" yaml:"functions" validate:"required,dive"`
}
