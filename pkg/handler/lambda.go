package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws/lambda"
	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/pkg/metrics"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

// Lambda for invocation.
type Lambda struct {
	// Invoker to be used.
	Invoker lambda.Facade
	// FnRef of the function to be invoked.
	FnRef lambda.FnRef
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
		log.Error(err, "invocation of lambda failed", "ref", r.FnRef)
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

func (r *Lambda) invokeMetered(request *http.Request, event *lambda.Request) (*lambda.Response, error) {
	now := time.Now()
	var (
		result *lambda.Response
		err    error
	)
	defer func() {
		lbls := prometheus.Labels{
			metrics.FunctionNameLabel:      r.FnRef.Name,
			metrics.InvocationRoleArnLabel: config.ArnAsString(r.FnRef.RoleARN),
		}
		metrics.AwsLambdaInvocationTotal.With(lbls).Inc()
		metrics.AwsLambdaInvocationDuration.With(lbls).Observe(float64(time.Since(now)))
		if err != nil {
			lbls[metrics.ErrorLabel] = err.Error()
			metrics.AwsLambdaInvocationErrors.With(lbls).Inc()
		}
	}()
	result, err = r.Invoker.Invoke(request.Context(), r.FnRef, event)
	return result, err
}

func (r *Lambda) adaptRequest(request *http.Request) (*lambda.Request, error) {
	result := lambda.Request{
		Host:    request.Host,
		Headers: lambda.Headers(request.Header),
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
