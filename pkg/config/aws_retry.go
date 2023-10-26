package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ratelimit"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
)

// AWSRetryRateLimiter config.
type AWSRetryRateLimiter struct {
	// RetryCost to deduct from the token bucket per retry.
	RetryCost uint `json:"retryCost,omitempty" yaml:"retryCost,omitempty" validate:"omitempty"`
	// RetryTimeoutCost to deduct from the token bucket per retry caused by timeout error.
	RetryTimeoutCost uint `json:"retryTimeoutCost,omitempty" yaml:"retryTimeoutCost,omitempty" validate:"omitempty"`
	// NoRetryIncrement to pay back to the token bucket for successful attempts.
	NoRetryIncrement uint `json:"noRetryIncrement,omitempty" yaml:"noRetryIncrement,omitempty" validate:"omitempty"`
	// Tokens to be obtained.
	Tokens uint `json:"tokens,omitempty" yaml:"tokens,omitempty" validate:"omitempty" map:"-"`
}

// AWSRetry configuration.
type AWSRetry struct {
	// MaxBackoff Duration for a retry.
	MaxBackoff Duration `json:"maxBackoff,omitempty" yaml:"maxBackoff,omitempty" validate:"omitempty"`
	// MaxAttempts for a request.
	MaxAttempts int `json:"maxAttempts,omitempty" yaml:"maxAttempts,omitempty" validate:"omitempty"`
	// RateLimiter configuration.
	RateLimiter AWSRetryRateLimiter `json:"rateLimiter,omitempty" yaml:"rateLimiter,omitempty" validate:"omitempty"`
}

func (r *AWSRetry) Apply(cfg *aws.Config) {
	if r == nil {
		return
	}
	var opts []func(*retry.StandardOptions)
	withNonZero(r.MaxBackoff, func(t Duration) {
		opts = append(opts, func(o *retry.StandardOptions) {
			o.MaxBackoff = t.Duration
			o.Backoff = retry.NewExponentialJitterBackoff(t.Duration)
		})
	})
	withNonZero(r.MaxAttempts, func(t int) {
		opts = append(opts, func(o *retry.StandardOptions) {
			o.MaxAttempts = t
		})
	})
	withNonZero(r.RateLimiter.RetryCost, func(t uint) {
		opts = append(opts, func(o *retry.StandardOptions) {
			o.RetryCost = t
		})
	})
	withNonZero(r.RateLimiter.RetryTimeoutCost, func(t uint) {
		opts = append(opts, func(o *retry.StandardOptions) {
			o.RetryTimeoutCost = t
		})
	})
	withNonZero(r.RateLimiter.NoRetryIncrement, func(t uint) {
		opts = append(opts, func(o *retry.StandardOptions) {
			o.NoRetryIncrement = t
		})
	})
	withNonZero(r.RateLimiter.Tokens, func(t uint) {
		opts = append(opts, func(o *retry.StandardOptions) {
			o.RateLimiter = ratelimit.NewTokenRateLimit(t)
		})
	})
	if len(opts) > 0 {
		cfg.Retryer = func() aws.Retryer { return retry.NewStandard(opts...) }
	}
}

func withNonZero[T comparable](t T, fn func(t T)) {
	var z T
	if t != z {
		fn(t)
	}
}
