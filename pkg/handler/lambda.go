package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws"
	"github.com/Tanemahuta/aws-lambda-server/pkg/metrics"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

// Lambda for invocation.
type Lambda struct {
	// Invoker to be used.
	Invoker aws.LambdaService
	// ARN of the function to be invoked.
	ARN arn.ARN
}

func (r *Lambda) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log := logr.FromContextOrDiscard(request.Context())
	event, err := r.adaptRequest(request)
	if err != nil {
		// Request could not be converted
		log.Error(err, "reading request failed")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	response, err := r.invokeMetered(request, event)
	if err != nil {
		// Invocation failed
		log.Error(err, "invocation of lambda failed", "arn", r.ARN)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Otherwise apply the response to the writer.
	for key, values := range response.Headers {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}
	writer.WriteHeader(response.StatusCode)
	if _, err = writer.Write(response.Body.Data); err != nil {
		log.Error(err, "could not write response body")
	}
}

func (r *Lambda) invokeMetered(request *http.Request, event *aws.LambdaRequest) (*aws.LambdaResponse, error) {
	now := time.Now()
	var (
		result *aws.LambdaResponse
		err    error
	)
	defer func() {
		lbls := prometheus.Labels{metrics.FunctionArnLabel: r.ARN.String()}
		metrics.AwsLambdaInvocationTotal.With(lbls).Inc()
		metrics.AwsLambdaInvocationDuration.With(lbls).Observe(float64(time.Since(now)))
		if err != nil {
			lbls[metrics.ErrorLabel] = err.Error()
			metrics.AwsLambdaInvocationErrors.With(lbls).Inc()
		}
	}()
	result, err = r.Invoker.Invoke(request.Context(), r.ARN, event)
	return result, err
}

func (r *Lambda) adaptRequest(request *http.Request) (*aws.LambdaRequest, error) {
	result := aws.LambdaRequest{
		Host:    request.Host,
		Headers: aws.Headers(request.Header),
		Method:  request.Method, URI: request.RequestURI,
		Vars: mux.Vars(request),
	}
	if request.Body != nil {
		defer func() { _ = request.Body.Close() }()
		var err error
		result.Body, err = io.ReadAll(request.Body)
		if err != nil {
			return nil, errors.Wrap(err, "could not read request body")
		}
	}
	return &result, nil
}
