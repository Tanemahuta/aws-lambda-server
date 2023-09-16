package sdk_test

import (
	"reflect"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws/sdk"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AssumeClients", func() {
	var sut sdk.AssumeClients[sdk.Lambda]
	BeforeEach(func() {
		sut = sdk.NewAssumeClients[sdk.Lambda](sdk.LambdaClientProps(aws.Config{}))
	})
	Context("Get()", func() {
		It("should return a lambda client for nil and cache it", func() {
			actual := sut.Get(nil)
			Expect(actual).NotTo(BeNil())
			Expect(reflect.ValueOf(sut.Get(nil)).UnsafePointer()).To(Equal(reflect.ValueOf(actual).UnsafePointer()))
		})
		It("should handle a role and cache it", func() {
			role := &arn.ARN{}
			actual := sut.Get(role)
			Expect(actual).NotTo(BeNil())
			Expect(reflect.ValueOf(sut.Get(role)).UnsafePointer()).To(Equal(reflect.ValueOf(actual).UnsafePointer()))
		})
	})
})
