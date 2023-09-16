package testcontext

import (
	"context"

	"github.com/Tanemahuta/aws-lambda-server/testing/testlogr"
	"github.com/go-logr/logr"
)

// New creates a new context.Context with a logr.Logger.
func New() context.Context {
	return logr.NewContext(context.TODO(), testlogr.New())
}
