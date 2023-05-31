package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws"
	"github.com/Tanemahuta/aws-lambda-server/pkg/handler"
	"github.com/Tanemahuta/aws-lambda-server/pkg/metrics"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

type UnreadableBody struct{}

func (u UnreadableBody) Read([]byte) (int, error) {
	return 0, errors.New("meh")
}

func (u UnreadableBody) Close() error {
	return errors.New("meh")
}

var _ = Describe("Lambda", func() {
	var (
		sut            *handler.Lambda
		lambdaRequest  *aws.LambdaRequest
		lambdaResponse *aws.LambdaResponse
		errLambda      error
		requestVars    map[string]string
		httpRequest    *http.Request
		httpResponse   *httptest.ResponseRecorder
	)
	BeforeEach(func() {
		lambdaArn := arn.ARN{
			Partition: "aws", Service: "lambda", Region: "eu-central-1", AccountID: "0123456789012",
			Resource: "function:my-function",
		}
		requestVars = map[string]string{"a": "b"}
		lambdaRequest = &aws.LambdaRequest{
			Host: "www.example.com",
			Headers: aws.Headers{
				"Test": []string{"test"},
			},
			Method: http.MethodPost,
			URI:    "http://www.example.com/test",
			Vars:   requestVars,
			Body:   []byte("test"),
		}
		httpRequest = mux.SetURLVars(
			httptest.NewRequest(http.MethodPost, "http://www.example.com/test", bytes.NewBufferString("test")),
			requestVars,
		)
		httpRequest.Header.Set("test", "test")
		httpResponse = httptest.NewRecorder()
		lambdaResponse = nil
		errLambda = nil
		sut = &handler.Lambda{
			Invoker: LambdaServiceFn(func(ctx context.Context, arn arn.ARN, r *aws.LambdaRequest) (
				*aws.LambdaResponse, error,
			) {
				Expect(ctx).NotTo(BeNil())
				Expect(arn).To(Equal(lambdaArn))
				Expect(r).To(Equal(lambdaRequest))
				return lambdaResponse, errLambda
			}),
			ARN: lambdaArn,
		}
	})
	AfterEach(func() {
		metrics.AwsLambdaInvocationTotal.Reset()
		metrics.AwsLambdaInvocationErrors.Reset()
		metrics.AwsLambdaInvocationDuration.Reset()
	})
	When("receiving a valid lambda response", func() {
		BeforeEach(func() {
			lambdaResponse = &aws.LambdaResponse{
				StatusCode: http.StatusAccepted,
				Headers:    aws.Headers{"test": []string{"test"}},
				Body:       aws.Body{Data: []byte("test")},
			}
			Expect(func() { sut.ServeHTTP(httpResponse, httpRequest) }).NotTo(Panic())
		})
		It("should adapt to http correctly", func() {
			Expect(httpResponse.Code).To(Equal(lambdaResponse.StatusCode))
			Expect(httpResponse.Header()).To(Equal(http.Header{"Test": []string{"test"}}))
			Expect(httpResponse.Body.String()).To(Equal(lambdaResponse.Body.String()))
		})
		It("should add metrics", func() {
			Expect(metrics.Collect(metrics.AwsLambdaInvocationTotal)).To(HaveKeyWithValue(
				"functionArn=arn:aws:lambda:eu-central-1:0123456789012:function:my-function",
				BeNumerically("==", 1),
			))
			Expect(metrics.Collect(metrics.AwsLambdaInvocationErrors)).To(BeEmpty())
			Expect(metrics.Collect(metrics.AwsLambdaInvocationDuration)).To(HaveKeyWithValue(
				"functionArn=arn:aws:lambda:eu-central-1:0123456789012:function:my-function",
				BeNumerically("==", 1),
			))
		})
	})
	When("body cannot be read", func() {
		BeforeEach(func() {
			httpRequest.Body = UnreadableBody{}
			Expect(func() { sut.ServeHTTP(httpResponse, httpRequest) }).NotTo(Panic())
		})
		It("should convert error to response code", func() {
			Expect(httpResponse.Code).To(Equal(http.StatusBadRequest))
		})
		It("should not add metrics", func() {
			Expect(metrics.Collect(metrics.AwsLambdaInvocationTotal)).To(BeEmpty())
			Expect(metrics.Collect(metrics.AwsLambdaInvocationErrors)).To(BeEmpty())
			Expect(metrics.Collect(metrics.AwsLambdaInvocationDuration)).To(BeEmpty())
		})
	})
	When("lambda errors", func() {
		BeforeEach(func() {
			errLambda = errors.New("meh")
			Expect(func() { sut.ServeHTTP(httpResponse, httpRequest) }).NotTo(Panic())
		})
		It("should convert error to response code", func() {
			Expect(httpResponse.Code).To(Equal(http.StatusInternalServerError))
		})
		It("should add metrics", func() {
			Expect(metrics.Collect(metrics.AwsLambdaInvocationTotal)).NotTo(BeNil())
			Expect(metrics.Collect(metrics.AwsLambdaInvocationErrors)).NotTo(BeNil())
			Expect(metrics.Collect(metrics.AwsLambdaInvocationDuration)).NotTo(BeNil())
		})
	})
})
