package routing_test

import (
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws/lambda"
	"github.com/Tanemahuta/aws-lambda-server/pkg/metrics"
	"github.com/Tanemahuta/aws-lambda-server/pkg/routing"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var _ = Describe("CurryMeteringFactory()", func() {
	It("should add functionArn", func() {
		metric := prometheus.NewCounterVec(
			prometheus.CounterOpts(prometheus.Opts{Name: "test"}),
			[]string{metrics.FunctionNameLabel, metrics.InvocationRoleArnLabel},
		)
		factory := routing.CurryMeteringFactory[*prometheus.CounterVec](
			func(o *prometheus.CounterVec, _ http.Handler, option ...promhttp.Option) http.HandlerFunc {
				return func(writer http.ResponseWriter, request *http.Request) {
					o.With(prometheus.Labels{}).Inc()
				}
			},
			metric,
		)
		fnRef := lambda.FnRef{
			Name: "test",
			RoleARN: &arn.ARN{
				Partition: "aws", Service: "iam", Region: "", AccountID: "123456789012", Resource: "role/test-role",
			},
		}
		Expect(func() { factory(nil, fnRef).ServeHTTP(nil, nil) }).NotTo(Panic())
		Expect(metrics.Collect(metric)).To(HaveKeyWithValue(
			"functionName=test,invocationRole=arn:aws:iam::123456789012:role/test-role",
			BeNumerically("==", 1),
		))
	})
})
