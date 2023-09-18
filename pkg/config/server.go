package config

// Server for the handler.
type Server struct {
	// HTTP configuration
	HTTP HTTP `json:"http,omitempty" yaml:"http,omitempty"`
	// DisableValidation disables the validation of the function using a dry-run invocation.
	DisableValidation bool `json:"disableValidation,omitempty" yaml:"disableValidation,omitempty"`
	// Functions providing the Function.
	Functions []Function `json:"functions" yaml:"functions" validate:"required,dive"`
}
