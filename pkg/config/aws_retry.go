package config

import (
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ratelimit"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/pkg/errors"
)

var _ Wrapped = &AWSRetryRateLimiter{}

// AWSRetryRateLimiter config.
type AWSRetryRateLimiter struct {
	// Tokens to be obtained.
	Tokens uint `json:"tokens,omitempty" yaml:"tokens,omitempty" validate:"omitempty"`
}

func (a AWSRetryRateLimiter) Unwrap() interface{} {
	if a.Tokens == 0 {
		return ratelimit.NewTokenRateLimit(retry.DefaultRetryRateTokens)
	}
	return ratelimit.NewTokenRateLimit(a.Tokens)
}

// AWSRetry configuration.
type AWSRetry struct {
	// MaxBackoff Duration for a retry.
	MaxBackoff Duration `json:"maxBackoff,omitempty" yaml:"maxBackoff,omitempty" validate:"omitempty"`
	// MaxAttempts for a request.
	MaxAttempts int `json:"maxAttempts,omitempty" yaml:"maxAttempts,omitempty" validate:"omitempty"`
	// The cost to deduct from the RateLimiter's token bucket per retry.
	RetryCost uint `json:"retryCost,omitempty" yaml:"retryCost,omitempty" validate:"omitempty"`
	// The cost to deduct from the RateLimiter's token bucket per retry caused by timeout error.
	RetryTimeoutCost uint `json:"retryTimeoutCost,omitempty" yaml:"retryTimeoutCost,omitempty" validate:"omitempty"`
	// The cost to payback to the RateLimiter's token bucket for successful attempts.
	NoRetryIncrement uint `json:"noRetryIncrement,omitempty" yaml:"noRetryIncrement,omitempty" validate:"omitempty"`
	// RateLimiter configuration.
	RateLimiter AWSRetryRateLimiter `json:"rateLimiter,omitempty" yaml:"rateLimiter,omitempty" validate:"omitempty"`
}

func (r *AWSRetry) Apply(cfg *aws.Config) error {
	if r == nil {
		return nil
	}
	srcVal, srcTpe := reflect.ValueOf(r).Elem(), reflect.TypeOf(r).Elem()
	tgtTpe := reflect.TypeOf(retry.StandardOptions{})
	var optsFuncs []func(*retry.StandardOptions)
	for idx := 0; idx < srcTpe.NumField(); idx++ {
		srcFld := srcTpe.Field(idx)
		srcFldVal := reflect.ValueOf(Unwrap(srcVal.FieldByName(srcFld.Name).Interface()))
		tgtFld, ok := tgtTpe.FieldByName(srcFld.Name)
		if !ok || !tgtFld.IsExported() {
			return errors.Errorf("field '%v' not found or exported", srcFld.Name)
		}
		if srcFldVal.CanConvert(tgtFld.Type) {
			srcFldVal = srcFldVal.Convert(tgtFld.Type)
		}
		if !srcFldVal.Type().AssignableTo(tgtFld.Type) {
			return errors.Errorf(
				"field '%v' type '%v' not assignable to type '%v", srcFld.Name, srcFldVal.Type(), tgtFld.Type,
			)
		}
		// Do not use zero values.
		if !srcFldVal.IsZero() {
			optsFuncs = append(optsFuncs, func(options *retry.StandardOptions) {
				reflect.ValueOf(options).Elem().FieldByName(tgtFld.Name).Set(srcFldVal)
			})
		}
	}
	if len(optsFuncs) > 0 {
		cfg.Retryer = func() aws.Retryer { return retry.NewStandard(optsFuncs...) }
	}
	return nil
}
