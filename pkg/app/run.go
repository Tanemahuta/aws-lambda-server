package app

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

func Run(ctx context.Context, appConfig Config) error {
	log := logr.FromContextOrDiscard(ctx)
	var (
		routerConfig  *config.Server
		lambdaService lambda.Facade
		requestRouter http.Handler
		err           error
	)
	return errorx.Fns{
		func() error {
			log.Info("reading config file", "filename", appConfig.Filename)
			routerConfig, err = config.Read(ctx, appConfig.Filename)
			return errors.Wrapf(err, "could not read routerConfig '%v'", appConfig.Filename)
		},
		func() error {
			log.Info("validating config", "filename", appConfig.Filename)
			return config.Validate(routerConfig)
		},
		func() error {
			log.Info("creating lambda service")
			lambdaService, err = appConfig.LambdaServiceFactory(ctx, routerConfig.AWS)
			return errors.Wrapf(err, "could not create lambda service '%v'", appConfig.Filename)
		},
		func() error {
			log.Info("creating app router")
			requestRouter, err = routing.New(lambdaService, routerConfig, routing.MetricsDecorators...)
			return errors.Wrapf(err, "could not create app router '%v'", appConfig.Filename)
		},
		func() error {
			group, runCtx := errgroup.WithContext(ctx)
			group.Go(func() error {
				log.Info("handling requests")
				return appConfig.RunFunc(runCtx, appConfig.Listen, requestRouter, &routerConfig.HTTP)
			})
			group.Go(func() error {
				log.Info("handling metrics")
				metricsRouter := mux.NewRouter()
				metricsRouter.NewRoute().Methods(http.MethodGet).Path("/metrics").Handler(promhttp.Handler())
				metricsRouter.NewRoute().Methods(http.MethodGet).Path("/healthz").HandlerFunc(ping)
				metricsRouter.NewRoute().Methods(http.MethodGet).Path("/readyz").HandlerFunc(ping)
				return appConfig.RunFunc(runCtx, appConfig.MetricsListen, metricsRouter, &routerConfig.HTTP)
			})
			return group.Wait()
		},
	}.Run()
}
