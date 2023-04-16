package server

import (
	"context"
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws"
)

// Config for the server.
type Config struct {
	// Filename of the config file.
	Filename string
	// Listen address.
	Listen string
	// LambdaServiceFactory to be used.
	LambdaServiceFactory func(context.Context) (aws.LambdaService, error)
	// RunFunc which runs the server.
	RunFunc func(context.Context, string, http.Handler) error
}
