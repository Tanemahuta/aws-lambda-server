package routing

import (
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MeterFactory for generic creation of MetricsDecorators handlers.
type MeterFactory[O any, H http.Handler] func(O, http.Handler, ...promhttp.Option) H

//nolint:gochecknoglobals // global decorators.
var MetricsDecorators = []Decorator{
	CurryMeteringFactory[*prometheus.CounterVec](promhttp.InstrumentHandlerCounter, metrics.HTTPRequestsTotal),
	CurryMeteringFactory[prometheus.ObserverVec](promhttp.InstrumentHandlerDuration, metrics.HTTPRequestsDuration),
	CurryMeteringFactory[prometheus.ObserverVec](promhttp.InstrumentHandlerRequestSize, metrics.HTTPRequestsSize),
	CurryMeteringFactory[prometheus.ObserverVec](promhttp.InstrumentHandlerResponseSize, metrics.HTTPResponsesSize),
}

// MeteringTarget allowing to curry labels.
type MeteringTarget[O any] interface {
	MustCurryWith(labels prometheus.Labels) O
}

// CurryMeteringFactory using the provided MeteringTarget and return a routing.Decorator from it.
func CurryMeteringFactory[O MeteringTarget[O], H http.Handler](fn MeterFactory[O, H], o O) Decorator {
	return func(handler http.Handler, functionArn string) http.Handler {
		return fn(o.MustCurryWith(prometheus.Labels{metrics.FunctionArnLabel: functionArn}), handler)
	}
}
