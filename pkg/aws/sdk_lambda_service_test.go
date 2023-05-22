package aws_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("SdkLambdaService", func() {
	var (
		sut             aws.SdkLambdaService
		testArn         arn.ARN
		request         *aws.LambdaRequest
		requestPayload  []byte
		response        *aws.LambdaResponse
		responsePayload []byte
		expectedInv     *lambda.InvokeInput
	)
	BeforeEach(func() {
		sut = nil
		request = &aws.LambdaRequest{
			Host:    "www.example.com",
			Headers: aws.Headers{Header: http.Header{"test": {"test"}}},
			Method:  "POST",
			URI:     "/test",
			Vars:    map[string]string{"a": "b"},
			Body:    []byte("test"),
		}
		var err error
		requestPayload, err = json.Marshal(request)
		Expect(err).NotTo(HaveOccurred())
		response = &aws.LambdaResponse{
			StatusCode: http.StatusAccepted,
			Headers:    map[string]string{"a": "b"},
			Body:       aws.Body{Data: []byte("test")},
		}
		responsePayload, err = json.Marshal(response)
		Expect(err).NotTo(HaveOccurred())
		testArn = arn.ARN{
			Partition: "aws",
			Service:   "lambda",
			Region:    "eu-central-1",
			AccountID: "123456789012",
			Resource:  "function:my-function",
		}
		testArnStr := testArn.String()
		expectedInv = &lambda.InvokeInput{
			FunctionName:   &testArnStr,
			InvocationType: types.InvocationTypeRequestResponse,
			LogType:        types.LogTypeTail,
			Payload:        requestPayload,
		}
	})
	It("should convert response", func() {
		sut = func(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (
			*lambda.InvokeOutput, error,
		) {
			Expect(params).To(Equal(expectedInv))
			logBytes := make([]byte, base64.StdEncoding.EncodedLen(4))
			base64.StdEncoding.Encode(logBytes, []byte("test"))
			logOutput := string(logBytes)
			return &lambda.InvokeOutput{Payload: responsePayload, LogResult: &logOutput}, nil
		}
		Expect(sut.Invoke(context.TODO(), testArn, request)).To(Equal(response))
	})
	It("should error if invocation errors", func() {
		sut = func(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (
			*lambda.InvokeOutput, error,
		) {
			return nil, errors.New("meh")
		}
		_, err := sut.Invoke(context.TODO(), testArn, request)
		Expect(err).To(MatchError(And(
			ContainSubstring("could not invoke lambda"),
			ContainSubstring("my-function"),
			ContainSubstring("meh"),
		)))
	})
	It("should error if payload cannot be unmarshalled", func() {
		sut = func(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (
			*lambda.InvokeOutput, error,
		) {
			return &lambda.InvokeOutput{Payload: []byte("6")}, nil
		}
		_, err := sut.Invoke(context.TODO(), testArn, request)
		Expect(err).To(MatchError(And(
			ContainSubstring("could not unmarshal payload to response"),
			ContainSubstring("cannot unmarshal number into Go value of type aws.LambdaResponse"),
		)))
	})
})

var _ = Describe("NewLambdaService()", func() {
	It("should create service or error", func() {
		res, err := aws.NewLambdaService(context.TODO())
		Expect([]interface{}{res, err}).To(Or(
			ConsistOf(Not(BeNil()), Not(HaveOccurred())),
			ConsistOf(BeNil(), HaveOccurred()),
		))
	})
})
