package lambda

import (
	"context"
)

//go:generate go run github.com/golang/mock/mockgen -destination=../../../mocks/mocklambda/mock_facade.go -package=mocklambda -source ./facade.go Facade

// Facade for invocation.
type Facade interface {
	// Invoke the lambda from the provided arn.ARN using the provided Request.
	Invoke(ctx context.Context, ref FnRef, request *Request) (*Response, error)
	// CanInvoke checks, if the lambda function can be invoked. Returns an error if this is not the case.
	CanInvoke(ctx context.Context, ref FnRef) error
}
