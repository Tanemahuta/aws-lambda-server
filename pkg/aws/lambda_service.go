package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
)

// LambdaService for invocation.
type LambdaService interface {
	// Invoke the lambda from the provided arn.ARN using the provided LambdaRequest.
	Invoke(ctx context.Context, arn arn.ARN, request *LambdaRequest) (*LambdaResponse, error)
}
