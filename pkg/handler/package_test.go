package handler_test

import (
	"context"
	"testing"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "handler Suite")
}

var _ aws.LambdaService = LambdaServiceFn(nil)

type LambdaServiceFn func(ctx context.Context, arn arn.ARN, request *aws.LambdaRequest) (*aws.LambdaResponse, error)

func (l LambdaServiceFn) CanInvoke(context.Context, arn.ARN) error {
	return nil
}

func (l LambdaServiceFn) Invoke(ctx context.Context, arn arn.ARN, request *aws.LambdaRequest) (
	*aws.LambdaResponse, error,
) {
	return l(ctx, arn, request)
}
