package mux

import (
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws"
	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/pkg/handler"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Decorator func(http.Handler, string) http.Handler

// New creates a new http.Handler for the aws.LambdaService using the Server.
func New(invoker aws.LambdaService, funcs []config.Function, decorators ...Decorator) (http.Handler, error) {
	result := mux.NewRouter()
	for fIdx, functionRoute := range funcs {
		var routeHandler http.Handler = &handler.Lambda{Invoker: invoker, ARN: functionRoute.ARN.ARN.ARN}
		for _, decorator := range decorators {
			routeHandler = decorator(routeHandler, "invoke "+functionRoute.ARN.ARN.String())
		}
		for rIdx, routeCfg := range functionRoute.Routes {
			route, err := ConfigureRoute(result.NewRoute(), routeCfg, (*mux.Route).GetError)
			if err != nil {
				return nil, errors.Wrapf(err, "could not apply route configuration %v of function %v", rIdx, fIdx)
			}
			route.Handler(routeHandler)
		}
	}
	return result, nil
}
