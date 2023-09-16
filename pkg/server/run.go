package server

import (
	"context"
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws/lambda"
	"github.com/Tanemahuta/aws-lambda-server/pkg/errorx"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"

	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/pkg/routing"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

func Run(ctx context.Context, serverConfig Config) error {
	log := logr.FromContextOrDiscard(ctx)
	log.Info("reading config from", "filename", serverConfig.Filename)
	var (
		routerConfig  *config.Server
		lambdaService lambda.Facade
		requestRouter http.Handler
		err           error
	)
	return errorx.Fns{
		func() error {
			log.Info("reading config file", "filename", serverConfig.Filename)
			routerConfig, err = config.Read(ctx, serverConfig.Filename)
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
			requestRouter, err = routing.New(lambdaService, routerConfig, routing.MetricsDecorators...)
			return errors.Wrapf(err, "could not create server router '%v'", serverConfig.Filename)
		},
		func() error {
			group, runCtx := errgroup.WithContext(ctx)
			group.Go(func() error {
				log.Info("handling requests")
				return serverConfig.RunFunc(runCtx, serverConfig.Listen, requestRouter)
			})
			group.Go(func() error {
				log.Info("handling metrics")
				metricsRouter := mux.NewRouter()
				metricsRouter.NewRoute().Methods(http.MethodGet).Path("/metrics").Handler(promhttp.Handler())
				metricsRouter.NewRoute().Methods(http.MethodGet).Path("/healthz").HandlerFunc(ping)
				metricsRouter.NewRoute().Methods(http.MethodGet).Path("/readyz").HandlerFunc(ping)
				return serverConfig.RunFunc(runCtx, serverConfig.MetricsListen, metricsRouter)
			})
			return group.Wait()
		},
	}.Run()
}
