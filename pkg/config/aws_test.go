package config_test

import (
	"path"
	"reflect"
	"time"

	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/testing/testcontext"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ratelimit"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AWS", func() {
	var (
		sut *config.AWS
		tgt *aws.Config
	)
	BeforeEach(func() {
		sut = nil
		tgt = &aws.Config{}
	})
	Context("Apply()", func() {
		It("should skip nil", func() {
			sut.Apply(nil)
		})
		It("should skip nil retry", func() {
			sut = &config.AWS{}
			sut.Apply(tgt)
			Expect(tgt.Retryer).To(BeNil())
		})
		It("should skip zero vals", func() {
			sut = &config.AWS{Retry: &config.AWSRetry{}}
			sut.Apply(tgt)
			Expect(tgt.Retryer).To(BeNil())
		})
		It("should map example to config", func() {
			exampleCfg, err := config.Read(testcontext.New(), path.Join("testdata", "config.yaml"))
			Expect(err).NotTo(HaveOccurred())
			sut = exampleCfg.AWS
			sut.Apply(tgt)
			Expect(tgt.Retryer).NotTo(BeNil())
			Expect(extractOptions(tgt.Retryer())).To(And(
				HaveField("MaxAttempts", BeNumerically("==", 11)),
				HaveField("MaxBackoff", BeNumerically("==", 11*time.Second)),
				HaveField("Backoff", Not(BeNil())),
				HaveField("RateLimiter", BeEquivalentTo(ratelimit.NewTokenRateLimit(111))),
				HaveField("RetryCost", BeNumerically("==", 11)),
				HaveField("RetryTimeoutCost", BeNumerically("==", 11)),
				HaveField("NoRetryIncrement", BeNumerically("==", 11)),
			))
		})
	})
})

// I know that's not nice. Sorry for that.
func extractOptions(std aws.Retryer) retry.StandardOptions {
	fld := reflect.ValueOf(std).Elem().FieldByName("options")
	return reflect.NewAt(fld.Type(), fld.Addr().UnsafePointer()).Elem().Interface().(retry.StandardOptions)
}
