package handler

import (
	"io"
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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
	response, err := r.Invoker.Invoke(request.Context(), r.ARN, event)
	if err != nil {
		// Invocation failed
		log.Error(err, "invocation of lambda failed", "arn", r.ARN)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Otherwise apply the response to the writer.
	for key, value := range response.Headers {
		writer.Header().Set(key, value)
	}
	writer.WriteHeader(response.StatusCode)
	if _, err = writer.Write(response.Body.Data); err != nil {
		log.Error(err, "could not write response body")
	}
}

func (r *Lambda) adaptRequest(request *http.Request) (*aws.LambdaRequest, error) {
	result := aws.LambdaRequest{
		Host:    request.Host,
		Headers: aws.Headers{Header: request.Header},
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
