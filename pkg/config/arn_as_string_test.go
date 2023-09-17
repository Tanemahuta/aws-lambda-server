package config_test

import (
	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArnAsString()", func() {
	It("should convert nil", func() {
		Expect(config.ArnAsString(nil)).To(BeEmpty())
	})
	It("should delegate to non-nil", func() {
		Expect(config.ArnAsString(&arn.ARN{})).To(Equal("arn:::::"))
	})
})
