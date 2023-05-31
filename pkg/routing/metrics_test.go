package routing_test

import (
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/metrics"
	"github.com/Tanemahuta/aws-lambda-server/pkg/routing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var _ = Describe("CurryMeteringFactory()", func() {
	It("should add functionArn", func() {
		metric := prometheus.NewCounterVec(
			prometheus.CounterOpts(prometheus.Opts{Name: "test"}), []string{metrics.FunctionArnLabel},
		)
		factory := routing.CurryMeteringFactory[*prometheus.CounterVec](
			func(o *prometheus.CounterVec, _ http.Handler, option ...promhttp.Option) http.HandlerFunc {
				return func(writer http.ResponseWriter, request *http.Request) {
					o.With(prometheus.Labels{}).Inc()
				}
			},
			metric,
		)
		Expect(func() { factory(nil, "test").ServeHTTP(nil, nil) }).NotTo(Panic())
		Expect(metrics.Collect(metric)).To(HaveKeyWithValue(
			"functionArn=test", BeNumerically("==", 1),
		))
	})
})
