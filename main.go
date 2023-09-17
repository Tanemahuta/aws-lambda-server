package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Tanemahuta/aws-lambda-server/buildinfo"
	"github.com/Tanemahuta/aws-lambda-server/pkg/aws/lambda"
	"github.com/go-logr/logr"

	"github.com/Tanemahuta/aws-lambda-server/pkg/server"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func main() {
	devel := false
	serverConfig := server.Config{
		Filename:             "/etc/aws-lambda-server/config.yaml",
		Listen:               ":8080",
		MetricsListen:        ":8081",
		LambdaServiceFactory: lambda.NewLambdaService,
		RunFunc: func(ctx context.Context, addr string, handler http.Handler) error {
			httpSrv := &http.Server{
				Addr:    addr,
				Handler: handler,
				BaseContext: func(net.Listener) context.Context {
					return ctx
				},
				ReadHeaderTimeout: time.Millisecond * 500, //nolint:gomnd // just no.
			}
			return httpSrv.ListenAndServe()
		},
	}
	flag.BoolVar(
		&devel, "devel", devel, "activate development logging",
	)
	flag.StringVar(
		&serverConfig.Filename, "config-file", serverConfig.Filename, "use config file (yaml or json)",
	)
	flag.StringVar(
		&serverConfig.Listen, "listen", serverConfig.Listen, "listener address for lambda requests",
	)
	flag.StringVar(
		&serverConfig.MetricsListen, "metrics-listen", serverConfig.MetricsListen, "listener address for metrics requests",
	)
	flag.Parse()
	zapLog, err := createLogger(devel)
	if err != nil {
		panic(err)
	}
	log := zapr.NewLogger(zapLog)
	log.Info("starting aws-lambda-server",
		"version", buildinfo.Version(),
		"commitSHA", buildinfo.CommitSHA(),
		"timestamp", buildinfo.Timestamp(),
	)
	ctx := logr.NewContext(context.Background(), log)
	if err = server.Run(ctx, serverConfig); err != nil {
		log.Error(err, "could not run server")
		os.Exit(1)
	}
}

func createLogger(debug bool) (*zap.Logger, error) {
	if debug {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}
