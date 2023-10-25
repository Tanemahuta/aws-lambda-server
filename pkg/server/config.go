package server

import (
	"context"
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws/lambda"
	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
)

// Config for the server.
type Config struct {
	// Filename of the config file.
	Filename string
	// Listen address for requests.
	Listen string
	// MetricsListen address for metrics/health checks.
	MetricsListen string
	// LambdaServiceFactory to be used.
	LambdaServiceFactory func(context.Context, *config.AWS) (lambda.Facade, error)
	// RunFunc which runs the server.
	RunFunc func(context.Context, string, http.Handler, *config.HTTP) error
}
