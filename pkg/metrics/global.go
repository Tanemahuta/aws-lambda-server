package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	FunctionArnLabel = "functionArn"
	ErrorLabel       = "error"
)

var (
	AwsLambdaInvocationLabels = []string{FunctionArnLabel}
	HttpRequestsLabels        = []string{"method", "code", FunctionArnLabel}

	HttpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total", Help: "total count of http requests",
	}, HttpRequestsLabels)
	HttpRequestsDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds", Help: "duration of http requests in seconds",
	}, HttpRequestsLabels)
	HttpRequestsSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_size_bytes", Help: "total size of http requests",
	}, HttpRequestsLabels)
	HttpResponsesSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_response_size_bytes", Help: "total size of http responses",
	}, HttpRequestsLabels)
	AwsLambdaInvocationTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "aws_lambda_invocation_total", Help: "total count of AWS lambda invocations by ARN",
	}, AwsLambdaInvocationLabels)
	AwsLambdaInvocationErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "aws_lambda_invocation_errors", Help: "AWS lambda invocation errors by ARN",
	}, append(AwsLambdaInvocationLabels, ErrorLabel))
	AwsLambdaInvocationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "aws_lambda_invocation_duration_seconds", Help: "duration of AWS lambda invocations by ARN",
	}, AwsLambdaInvocationLabels)
)
