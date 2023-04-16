package handler

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Traceparent for traceparent in http.Request.
type Traceparent struct {
	Delegate http.Handler
	Prop     propagation.TextMapPropagator
}

func (t *Traceparent) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	t.Prop.Inject(request.Context(), propagation.HeaderCarrier(writer.Header()))
	t.Delegate.ServeHTTP(writer, request)
}

// NewTraceparent creates a new http.Handler for tracing using the operation name.
func NewTraceparent(delegate http.Handler, operation string) http.Handler {
	return otelhttp.NewHandler(
		&Traceparent{Delegate: delegate, Prop: otel.GetTextMapPropagator()},
		operation,
	)
}

//nolint:gochecknoinits // nope.
func init() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
}
