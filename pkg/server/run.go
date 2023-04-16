package server

import (
	"context"
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws"
	"github.com/Tanemahuta/aws-lambda-server/pkg/errorx"

	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/pkg/mux"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

func Run(ctx context.Context, serverConfig Config) error {
	log := logr.FromContextOrDiscard(ctx)
	log.Info("reading config from", "filename", serverConfig.Filename)
	var (
		routerConfig  *config.Server
		lambdaService aws.LambdaService
		handler       http.Handler
		err           error
	)
	return errorx.Fns{
		func() error {
			log.Info("reading config file", "filename", serverConfig.Filename)
			routerConfig, err = config.Read(serverConfig.Filename)
			return errors.Wrapf(err, "could not read routerConfig '%v'", serverConfig.Filename)
		},
		func() error {
			log.Info("validating config", "filename", serverConfig.Filename)
			return config.Validate(routerConfig)
		},
		func() error {
			log.Info("creating lambda service")
			lambdaService, err = serverConfig.LambdaServiceFactory(ctx)
			return errors.Wrapf(err, "could not create lambda service '%v'", serverConfig.Filename)
		},
		func() error {
			log.Info("creating server router")
			handler, err = mux.New(lambdaService, routerConfig.Functions)
			return errors.Wrapf(err, "could not create lambda service '%v'", serverConfig.Filename)
		},
		func() error {
			log.Info("handling requests")
			return serverConfig.RunFunc(ctx, serverConfig.Listen, handler)
		},
	}.Run()
}
