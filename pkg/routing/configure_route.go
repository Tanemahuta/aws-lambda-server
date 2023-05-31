package routing

import (
	"reflect"

	"github.com/pkg/errors"
)

// ConfigureRoute configures a route using the provided configuration and validates it using the provided function.
func ConfigureRoute[R, C any](route R, config C, validateFn func(r R) error) (R, error) {
	cfgVal := reflect.ValueOf(config)
	cfgTpe := cfgVal.Type()
	routeVal := reflect.ValueOf(route)
	routeTpe := reflect.TypeOf((*R)(nil)).Elem()
	for idx := 0; idx < cfgVal.NumField(); idx++ {
		fld := cfgTpe.Field(idx)
		nextRoute, err := applyConfigFieldToRoute(fld, cfgVal.Field(idx), routeVal, routeTpe)
		if err != nil {
			return route, errors.Wrapf(err, "could not apply field '%v'", fld.Name)
		}
		//nolint:errcheck // this is correct.
		route = nextRoute.Interface().(R)
	}
	return route, validateFn(route)
}

func applyConfigFieldToRoute(fld reflect.StructField, val reflect.Value, rVal reflect.Value, rTpe reflect.Type) (
	reflect.Value, error,
) {
	if !fld.IsExported() || val.IsZero() {
		return rVal, nil // Skip unexported fields and zero config values
	}
	configFn, ok := rVal.Type().MethodByName(fld.Name)
	if !ok || !configFn.IsExported() {
		return reflect.Value{}, errors.New("could not find exported config function")
	}
	configFnTpe := configFn.Type
	//nolint:gomnd // gosh.
	if numIn := configFnTpe.NumIn(); numIn != 2 {
		return reflect.Value{}, errors.Errorf("expected two in parameters, but got %v", numIn)
	}
	if numOut := configFnTpe.NumOut(); numOut != 1 {
		return reflect.Value{}, errors.Errorf("expected single out parameters, but got %v", numOut)
	}
	if outTpe := configFnTpe.Out(0); !outTpe.AssignableTo(rTpe) {
		return reflect.Value{}, errors.Errorf("expected out parameter to be assignable to %v (%v)", rTpe, outTpe)
	}
	convVal, err := adaptConfigValue(val, configFnTpe.In(1))
	if err != nil {
		return reflect.Value{}, errors.Wrap(err, "could not convert config value to function input")
	}
	return func() reflect.Value {
		defer func() {
			if recovered := recover(); recovered != nil {
				err = errors.Errorf("config func invocation failed: %v", recovered)
			}
		}()
		args := []reflect.Value{rVal}
		if configFn.Type.IsVariadic() && convVal.Kind() == reflect.Slice {
			for idx := 0; idx < convVal.Len(); idx++ {
				args = append(args, convVal.Index(idx))
			}
		} else {
			args = append(args, convVal)
		}
		return configFn.Func.Call(args)[0]
	}(), err
}

func adaptConfigValue(val reflect.Value, in reflect.Type) (reflect.Value, error) {
	tpe := val.Type()
	// Assignable
	if tpe.AssignableTo(in) {
		return val, nil
	}
	// Convertible
	if tpe.ConvertibleTo(in) {
		return val.Convert(in), nil
	}
	// Map to slice
	if tpe.Kind() == reflect.Map && in.Kind() == reflect.Slice {
		elemTpe, keyTpe, valTpe := in.Elem(), tpe.Key(), tpe.Elem()
		// Only works in case key and value are of elem type
		if !keyTpe.AssignableTo(elemTpe) || !valTpe.AssignableTo(elemTpe) {
			return reflect.Value{}, errors.Errorf(
				"could not convert map to slice, since key type '%v' and value type '%v' are not assignable to '%v'",
				keyTpe, valTpe, elemTpe,
			)
		}
		// Create a new slice of the provided type
		//nolint:gomnd // gosh.
		size := val.Len() * 2
		result := reflect.MakeSlice(in, size, size)
		for iter, idx := val.MapRange(), 0; iter.Next(); idx += 2 {
			result.Index(idx).Set(iter.Key())
			result.Index(idx + 1).Set(iter.Value())
		}
		return result, nil
	}
	return reflect.Value{}, errors.Errorf("could not adapt %v to %v", tpe, in)
}
