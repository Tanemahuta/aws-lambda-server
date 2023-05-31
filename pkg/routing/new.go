package routing

import (
	"context"
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws"
	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/pkg/handler"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Decorator for a http request.
type Decorator func(decorated http.Handler, functionArn string) http.Handler

// New creates a new http.Handler for the aws.LambdaService using the Server.
func New(invoker aws.LambdaService, funcs []config.Function, decorators ...Decorator) (http.Handler, error) {
	result := mux.NewRouter()
	for fIdx, functionRoute := range funcs {
		if functionRoute.ARN.AccountID != "000000000000" { // This is a test account
			if err := invoker.CanInvoke(context.TODO(), functionRoute.ARN.ARN.ARN); err != nil {
				return nil, err
			}
		}
		var routeHandler http.Handler = &handler.Lambda{Invoker: invoker, ARN: functionRoute.ARN.ARN.ARN}
		for _, decorator := range decorators {
			routeHandler = decorator(routeHandler, functionRoute.ARN.ARN.String())
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
