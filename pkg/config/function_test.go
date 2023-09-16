package config_test

import (
	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Function", func() {
	var sut *config.Function
	BeforeEach(func() {
		sut = &config.Function{}
	})
	Context("GetName()", func() {
		It("should favour Name before ARN", func() {
			sut.Name = "test-function"
			sut.ARN = &config.LambdaARN{}
			Expect(sut.GetName()).To(Equal("test-function"))
		})
		It("should provide from Name", func() {
			sut.Name = "test-function"
			Expect(sut.GetName()).To(Equal("test-function"))
		})
		It("should fall back to ARN", func() {
			sut.ARN = &config.LambdaARN{}
			Expect(sut.GetName()).To(Equal("arn:::::"))
		})
	})
	Context("GetInvocationRoleARN()", func() {
		It("should return non-nil", func() {
			sut.InvocationRole = &config.RoleARN{}
			Expect(sut.GetInvocationRoleARN()).To(Equal(&sut.InvocationRole.ARN.ARN))
		})
		It("should return nil", func() {
			Expect(sut.GetInvocationRoleARN()).To(BeNil())
		})
	})
})
