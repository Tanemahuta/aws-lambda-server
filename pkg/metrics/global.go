//nolint:gochecknoglobals // global registrations.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	FunctionNameLabel      = "functionName"
	InvocationRoleArnLabel = "invocationRole"
	ErrorLabel             = "error"
)

var (
	AwsLambdaInvocationLabels = []string{FunctionNameLabel, InvocationRoleArnLabel}
	HTTPRequestLabels         = []string{"method", "code", FunctionNameLabel, InvocationRoleArnLabel}

	HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total", Help: "total count of http requests",
	}, HTTPRequestLabels)
	HTTPRequestsDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds", Help: "duration of http requests in seconds",
	}, HTTPRequestLabels)
	HTTPRequestsSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_size_bytes", Help: "total size of http requests",
	}, HTTPRequestLabels)
	HTTPResponsesSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_response_size_bytes", Help: "total size of http responses",
	}, HTTPRequestLabels)
	AwsLambdaInvocationTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "aws_lambda_invocation_total", Help: "total count of AWS lambda invocations by ARN",
	}, AwsLambdaInvocationLabels)
	AwsLambdaInvocationErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "aws_lambda_invocation_errors_total", Help: "AWS lambda invocation errors by ARN",
	}, append(AwsLambdaInvocationLabels, ErrorLabel))
	AwsLambdaInvocationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "aws_lambda_invocation_duration_seconds", Help: "duration of AWS lambda invocations by ARN",
	}, AwsLambdaInvocationLabels)
)
