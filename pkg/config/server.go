package config

// Server for the handler.
type Server struct {
	// EnableTraceparent injects traceparent.
	EnableTraceparent bool `json:"enableTraceparent,omitempty" yaml:"enableTraceparent,omitempty"`
	// Functions providing the Function.
	Functions []Function `json:"functions" yaml:"functions" validate:"required,dive"`
}
