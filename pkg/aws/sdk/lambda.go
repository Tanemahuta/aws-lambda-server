package sdk

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

//go:generate go run github.com/golang/mock/mockgen -destination=../../../mocks/mocksdk/mock_lambda.go -package=mocksdk -source ./lambda.go Lambda

// Lambda sdk interface.
type Lambda interface {
	Invoke(ctx context.Context, params *lambda.InvokeInput, opts ...func(*lambda.Options)) (*lambda.InvokeOutput, error)
}
