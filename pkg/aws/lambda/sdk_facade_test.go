package lambda_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/mocks/mocksdk"
	"github.com/Tanemahuta/aws-lambda-server/pkg/aws/lambda"
	"github.com/Tanemahuta/aws-lambda-server/pkg/aws/sdk"
	"github.com/Tanemahuta/aws-lambda-server/testing/testcontext"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	awslambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("SdkFacade", func() {
	var (
		ctrl              *gomock.Controller
		mockAssumeClients *mocksdk.MockAssumeClients[sdk.Lambda]
		mockLambda        *mocksdk.MockLambda
		sut               *lambda.SdkFacade
		functionName      string
		roleArn           arn.ARN
		request           *lambda.Request
		response          *lambda.Response
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockAssumeClients = mocksdk.NewMockAssumeClients[sdk.Lambda](ctrl)
		mockLambda = mocksdk.NewMockLambda(ctrl)
		sut = &lambda.SdkFacade{Clients: mockAssumeClients}
		request = &lambda.Request{
			Method:  "POST",
			Host:    "www.example.com",
			URI:     "http://www.example.com/test/a/b",
			Headers: lambda.Headers{"test": {"test"}},
			Vars:    map[string]string{"a": "b"},
			Body:    []byte("test"),
		}
		response = &lambda.Response{
			StatusCode: http.StatusAccepted,
			Headers:    lambda.Headers{"A": []string{"b"}},
			Body:       lambda.Body{Data: []byte("test")},
		}
		functionName = arn.ARN{
			Partition: "aws", Service: "lambda", Region: "eu-central-1", AccountID: "123456789012",
			Resource: "function:my-function",
		}.String()
		roleArn = arn.ARN{
			Partition: "aws", Service: "iam", AccountID: "234567890123", Resource: "role/my-function-invocation-role",
		}
	})
	AfterEach(func() {
		ctrl.Finish()
	})
	Context("CanInvoke()", func() {
		var err error
		BeforeEach(func() {
			err = errors.New("meh")
		})
		It("should pass-through error with role", func() {
			mockAssumeClients.EXPECT().Get(gomock.Eq(&roleArn)).Return(mockLambda)
			mockLambda.EXPECT().Invoke(gomock.Any(), gomock.Eq(createCanInput(functionName))).
				Return(nil, err)
			Expect(sut.CanInvoke(testcontext.New(), lambda.FnRef{Name: functionName, RoleARN: &roleArn})).
				To(MatchError(err))
		})
		It("should pass-through error without role", func() {
			mockAssumeClients.EXPECT().Get(gomock.Nil()).Return(mockLambda)
			mockLambda.EXPECT().Invoke(gomock.Any(), gomock.Eq(createCanInput(functionName))).
				Return(nil, err)
			Expect(sut.CanInvoke(testcontext.New(), lambda.FnRef{Name: functionName})).
				To(MatchError(err))
		})
	})
	Context("Invoke()", func() {
		It("should convert response with role", func() {
			mockAssumeClients.EXPECT().Get(gomock.Eq(&roleArn)).Return(mockLambda)
			mockLambda.EXPECT().Invoke(gomock.Any(), gomock.Eq(createInput(functionName, request))).
				Return(createOutput(response, true), nil)
			Expect(sut.Invoke(testcontext.New(), lambda.FnRef{Name: functionName, RoleARN: &roleArn}, request)).
				To(Equal(response))
		})
		It("should convert response without role", func() {
			mockAssumeClients.EXPECT().Get(gomock.Nil()).Return(mockLambda)
			mockLambda.EXPECT().Invoke(gomock.Any(), gomock.Eq(createInput(functionName, request))).
				Return(createOutput(response, false), nil)
			Expect(sut.Invoke(testcontext.New(), lambda.FnRef{Name: functionName}, request)).
				To(Equal(response))
		})
		It("should error if invocation errors", func() {
			lambdaErr := errors.New("meh")
			mockAssumeClients.EXPECT().Get(gomock.Nil()).Return(mockLambda)
			mockLambda.EXPECT().Invoke(gomock.Any(), gomock.Eq(createInput(functionName, request))).
				Return(nil, lambdaErr)
			_, err := sut.Invoke(testcontext.New(), lambda.FnRef{Name: functionName}, request)
			Expect(err).To(MatchError(And(
				ContainSubstring("could not invoke lambda"),
				ContainSubstring("my-function"),
				ContainSubstring(lambdaErr.Error()),
			)))
		})
		It("should error if payload cannot be unmarshalled", func() {
			mockAssumeClients.EXPECT().Get(gomock.Nil()).Return(mockLambda)
			mockLambda.EXPECT().Invoke(gomock.Any(), gomock.Eq(createInput(functionName, request))).
				Return(createRawOutput([]byte("6")), nil)
			_, err := sut.Invoke(testcontext.New(), lambda.FnRef{Name: functionName}, request)
			Expect(err).To(MatchError(And(
				ContainSubstring("could not adapt response"),
				ContainSubstring("could not unmarshal payload to response"),
				ContainSubstring("cannot unmarshal number into Go value of type lambda.Response"),
			)))
		})
	})
})

func createInput(functionName string, request *lambda.Request) *awslambda.InvokeInput {
	payload, err := json.Marshal(request)
	Expect(err).NotTo(HaveOccurred())
	return &awslambda.InvokeInput{
		FunctionName:   &functionName,
		InvocationType: types.InvocationTypeRequestResponse,
		LogType:        types.LogTypeNone,
		Payload:        payload,
	}
}

func createCanInput(functionName string) *awslambda.InvokeInput {
	return &awslambda.InvokeInput{
		FunctionName:   &functionName,
		InvocationType: types.InvocationTypeDryRun,
		LogType:        types.LogTypeNone,
	}
}

func createOutput(response *lambda.Response, base64 bool) *awslambda.InvokeOutput {
	payload, err := json.Marshal(response)
	Expect(err).NotTo(HaveOccurred())
	if base64 {
		payload = base64Enc(payload)
	}
	return createRawOutput(payload)
}

func createRawOutput(payload []byte) *awslambda.InvokeOutput {
	return &awslambda.InvokeOutput{Payload: payload}
}

func base64Enc(src []byte) []byte {
	result := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(result, src)
	return result
}

var _ = Describe("NewLambdaService()", func() {
	It("should create service or error", func() {
		res, err := lambda.NewLambdaService(testcontext.New())
		Expect([]interface{}{res, err}).To(Or(
			ConsistOf(Not(BeNil()), Not(HaveOccurred())),
			ConsistOf(BeNil(), HaveOccurred()),
		))
	})
})
