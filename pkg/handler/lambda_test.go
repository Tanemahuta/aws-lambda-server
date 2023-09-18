package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Tanemahuta/aws-lambda-server/mocks/mocklambda"
	"github.com/Tanemahuta/aws-lambda-server/pkg/aws/lambda"
	"github.com/Tanemahuta/aws-lambda-server/pkg/handler"
	"github.com/Tanemahuta/aws-lambda-server/pkg/metrics"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/golang/mock/gomock"
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
		ctrl          *gomock.Controller
		invokerMock   *mocklambda.MockFacade
		sut           *handler.Lambda
		lambdaRequest *lambda.Request
		requestVars   map[string]string
		httpRequest   *http.Request
		httpResponse  *httptest.ResponseRecorder
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		invokerMock = mocklambda.NewMockFacade(ctrl)
		requestVars = map[string]string{"a": "b"}
		lambdaRequest = &lambda.Request{
			Host: "www.example.com",
			Headers: lambda.Headers{
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
		sut = &handler.Lambda{
			Invoker: invokerMock,
			FnRef: lambda.FnRef{
				Name: "test-function",
				RoleARN: &arn.ARN{
					Partition: "aws", Service: "iam", AccountID: "123456789012", Resource: "role/test-role",
				},
			},
		}
	})
	AfterEach(func() {
		metrics.AwsLambdaInvocationTotal.Reset()
		metrics.AwsLambdaInvocationErrors.Reset()
		metrics.AwsLambdaInvocationDuration.Reset()
	})
	When("receiving a valid lambda response", func() {
		var lambdaResponse *lambda.Response
		BeforeEach(func() {
			lambdaResponse = &lambda.Response{
				StatusCode: http.StatusAccepted,
				Headers:    lambda.Headers{"test": []string{"test"}},
				Body:       lambda.Body{Data: []byte("test")},
			}
			invokerMock.EXPECT().Invoke(gomock.Any(), gomock.Eq(sut.FnRef), lambdaRequest).Return(lambdaResponse, nil)
			Expect(func() { sut.ServeHTTP(httpResponse, httpRequest) }).NotTo(Panic())
		})
		It("should adapt to http correctly", func() {
			Expect(httpResponse.Code).To(Equal(lambdaResponse.StatusCode))
			Expect(httpResponse.Header()).To(Equal(http.Header{"Test": []string{"test"}}))
			Expect(httpResponse.Body.String()).To(Equal(lambdaResponse.Body.String()))
		})
		It("should add metrics", func() {
			labelMatcher := HaveKeyWithValue(
				"functionName=test-function,invocationRole=arn:aws:iam::123456789012:role/test-role",
				BeNumerically("==", 1),
			)
			Expect(metrics.Collect(metrics.AwsLambdaInvocationTotal)).To(labelMatcher)
			Expect(metrics.Collect(metrics.AwsLambdaInvocationErrors)).To(BeEmpty())
			Expect(metrics.Collect(metrics.AwsLambdaInvocationDuration)).To(labelMatcher)
		})
	})
	When("invoking with timeout", func() {
		var lambdaResponse *lambda.Response
		BeforeEach(func() {
			lambdaResponse = &lambda.Response{
				StatusCode: http.StatusAccepted,
				Headers:    lambda.Headers{"test": []string{"test"}},
				Body:       lambda.Body{Data: []byte("test")},
			}
			sut.Timeout = time.Minute
			invokerMock.EXPECT().Invoke(gomock.Any(), gomock.Eq(sut.FnRef), lambdaRequest).DoAndReturn(
				func(ctx context.Context, ref lambda.FnRef, request *lambda.Request) (*lambda.Response, error) {
					defer GinkgoRecover()
					_, ok := ctx.Deadline()
					Expect(ok).To(BeTrue())
					return lambdaResponse, nil
				})
			Expect(func() { sut.ServeHTTP(httpResponse, httpRequest) }).NotTo(Panic())
		})
		It("should adapt to http correctly", func() {
			Expect(httpResponse.Code).To(Equal(lambdaResponse.StatusCode))
			Expect(httpResponse.Header()).To(Equal(http.Header{"Test": []string{"test"}}))
			Expect(httpResponse.Body.String()).To(Equal(lambdaResponse.Body.String()))
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
			invokerMock.EXPECT().Invoke(gomock.Any(), gomock.Eq(sut.FnRef), lambdaRequest).
				Return(nil, errors.New("meh"))
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
