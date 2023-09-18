package sdk_test

import (
	"context"
	"reflect"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws/sdk"
	"github.com/Tanemahuta/aws-lambda-server/testing/testcontext"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

const credErrMsg = "test-credentials"

var _ aws.CredentialsProvider = &testCredentials{}

type testCredentials struct{}

func (t testCredentials) Retrieve(context.Context) (aws.Credentials, error) {
	return aws.Credentials{}, errors.New(credErrMsg)
}

var _ = Describe("AssumeClients", func() {
	var (
		functionName string
		sut          sdk.AssumeClients[sdk.Lambda]
	)
	BeforeEach(func() {
		functionName = "test-function"
		sut = sdk.NewAssumeClients[sdk.Lambda](sdk.LambdaClientProps(aws.Config{
			Region:      "eu-central-1",
			Credentials: &testCredentials{},
		}))
	})
	Context("Get()", func() {
		It("should return a lambda client for nil and cache it", func() {
			actual := sut.Get(nil)
			Expect(actual).NotTo(BeNil())
			_, err := actual.Invoke(testcontext.New(), &lambda.InvokeInput{FunctionName: &functionName})
			Expect(err).To(MatchError(ContainSubstring(credErrMsg)))
			Expect(reflect.ValueOf(sut.Get(nil)).UnsafePointer()).To(Equal(reflect.ValueOf(actual).UnsafePointer()))
		})
		It("should handle a role and cache it", func() {
			role := &arn.ARN{
				Partition: "aws",
				Service:   "iam",
				Region:    "",
				AccountID: "123456789012",
				Resource:  "role/invocation-role",
			}
			actual := sut.Get(role)
			Expect(actual).NotTo(BeNil())
			_, err := actual.Invoke(testcontext.New(), &lambda.InvokeInput{FunctionName: &functionName})
			Expect(err).To(MatchError(And(
				ContainSubstring(credErrMsg), ContainSubstring("operation error STS: AssumeRole"),
			)))
			Expect(reflect.ValueOf(sut.Get(role)).UnsafePointer()).To(Equal(reflect.ValueOf(actual).UnsafePointer()))
		})
	})
})
