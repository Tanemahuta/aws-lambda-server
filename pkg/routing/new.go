package routing

import (
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws/lambda"
	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/pkg/handler"
	"github.com/Tanemahuta/aws-lambda-server/testing/testcontext"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Decorator for a http request.
type Decorator func(decorated http.Handler, ref lambda.FnRef) http.Handler

// New creates a new http.Handler for the aws.Facade using the Server.
func New(invoker lambda.Facade, cfg *config.Server, decorators ...Decorator) (http.Handler, error) {
	result := mux.NewRouter()
	for fIdx, functionRoute := range cfg.Functions {
		fnRef := lambda.FnRef{Name: functionRoute.GetName(), RoleARN: functionRoute.GetInvocationRoleARN()}
		if !cfg.DisableValidation {
			if err := invoker.CanInvoke(testcontext.New(), fnRef); err != nil {
				return nil, err
			}
		}
		var routeHandler http.Handler = &handler.Lambda{Invoker: invoker, FnRef: fnRef}
		for _, decorator := range decorators {
			routeHandler = decorator(routeHandler, fnRef)
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
