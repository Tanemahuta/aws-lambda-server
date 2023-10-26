package config_test

import (
	"path"
	"reflect"

	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/testing/testcontext"
	"github.com/aws/aws-sdk-go-v2/aws"
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
			Expect(sut.Apply(nil)).NotTo(HaveOccurred())
		})
		It("should skip nil retry", func() {
			sut = &config.AWS{}
			Expect(sut.Apply(tgt)).NotTo(HaveOccurred())
			Expect(tgt.Retryer).To(BeNil())
		})
		It("should should use defaults", func() {
			sut = &config.AWS{Retry: &config.AWSRetry{}}
			Expect(sut.Apply(tgt)).NotTo(HaveOccurred())
			Expect(tgt.Retryer).NotTo(BeNil())
			Expect(tgt.Retryer).NotTo(BeNil())
			Expect(extractOptions(tgt.Retryer())).To(BeEquivalentTo(extractOptions(retry.NewStandard())))
		})
		It("should apply example with default config", func() {
			exampleCfg, err := config.Read(testcontext.New(), path.Join("testdata", "config.yaml"))
			Expect(err).NotTo(HaveOccurred())
			sut = exampleCfg.AWS
			Expect(sut.Apply(tgt)).NotTo(HaveOccurred())
			Expect(tgt.Retryer).NotTo(BeNil())
			Expect(extractOptions(tgt.Retryer())).To(BeEquivalentTo(extractOptions(retry.NewStandard())))
		})
	})
})

// I know that's not nice. Sorry for that.
func extractOptions(std aws.Retryer) retry.StandardOptions {
	fld := reflect.ValueOf(std).Elem().FieldByName("options")
	return reflect.NewAt(fld.Type(), fld.Addr().UnsafePointer()).Elem().Interface().(retry.StandardOptions)
}
