package config_test

import (
	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LambdaARN validation", func() {
	var test config.LambdaARN
	BeforeEach(func() {
		test = config.LambdaARN{ARN: config.ARN{ARN: arn.ARN{
			Partition: "aws",
			Service:   "lambda",
			Region:    "eu-central-1",
			AccountID: "123456789012",
			Resource:  "function:my-function",
		}}}
	})
	It("should error on empty partition", func() {
		test.Partition = ""
		Expect(config.Validate(test)).To(MatchError(And(
			ContainSubstring("Partition"),
			ContainSubstring("required"),
		)))
	})
	It("should error on empty service", func() {
		test.Service = ""
		Expect(config.Validate(test)).To(MatchError(And(
			ContainSubstring("Service"),
			ContainSubstring("required"),
		)))
	})
	It("should error on empty account", func() {
		test.AccountID = ""
		Expect(config.Validate(test)).To(MatchError(And(
			ContainSubstring("AccountID"),
			ContainSubstring("required"),
		)))
	})
	It("should error on empty region", func() {
		test.Region = ""
		Expect(config.Validate(test)).To(MatchError(And(
			ContainSubstring("Region"),
			ContainSubstring("required"),
		)))
	})
	It("should error on invalid service", func() {
		test.Service = "iam"
		Expect(config.Validate(test)).To(MatchError(And(
			ContainSubstring("Service"),
			ContainSubstring("match=lambda"),
		)))
	})
	It("should error on invalid resource", func() {
		test.Resource = "layer:my-layer"
		Expect(config.Validate(test)).To(MatchError(And(
			ContainSubstring("Resource"),
			ContainSubstring("prefix=function:"),
		)))
	})
})
