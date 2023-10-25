package config

import "github.com/aws/aws-sdk-go-v2/aws"

// AWS configuration.
type AWS struct {
	// Retry configuration.
	Retry *AWSRetry `json:"retry,omitempty" yaml:"retry,omitempty" validate:"omitempty"`
}

func (a *AWS) Apply(cfg *aws.Config) error {
	if a == nil {
		return nil
	}
	return a.Retry.Apply(cfg)
}
