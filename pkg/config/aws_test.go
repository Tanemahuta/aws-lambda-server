package config_test

import (
	"path"

	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/testing/testcontext"
	"github.com/aws/aws-sdk-go-v2/aws"
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
		It("should apply example", func() {
			exampleCfg, err := config.Read(testcontext.New(), path.Join("testdata", "config.yaml"))
			Expect(err).NotTo(HaveOccurred())
			sut = exampleCfg.AWS
			Expect(sut.Apply(tgt)).NotTo(HaveOccurred())
			Expect(tgt.Retryer).NotTo(BeNil())
		})
	})
})
