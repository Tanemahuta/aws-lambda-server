package config

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/arn"

	"github.com/go-playground/validator/v10"
)

// Validate the Server config.
func Validate(config any) error {
	validate := validator.New()
	validate.RegisterStructValidation(validateARN, ARN{})
	validate.RegisterStructValidation(
		validateSpecARN[LambdaARN]("lambda", "function:"),
		LambdaARN{},
	)
	return validate.Struct(config)
}

func validateARN(sl validator.StructLevel) {
	doValidateARN(sl, sl.Current().Interface().(ARN).ARN)
}

func doValidateARN(sl validator.StructLevel, arn arn.ARN) {
	if len(arn.Partition) == 0 {
		sl.ReportError(arn, "Partition", "Partition", "partition-empty", "")
	}
	if len(arn.Service) == 0 {
		sl.ReportError(arn, "Service", "Service", "not-empty", "")
	}
}

type arnProvider interface {
	wrapped() arn.ARN
}

func validateSpecARN[T arnProvider](service, resourcePrefix string) validator.StructLevelFunc {
	return func(sl validator.StructLevel) {
		tgt := sl.Current().Interface().(T).wrapped()
		doValidateARN(sl, tgt)
		if len(tgt.Partition) == 0 {
			sl.ReportError(tgt, "Partition", "Partition", "partition-empty", "")
		}
		if tgt.Service != service {
			sl.ReportError(tgt, "Service", "Service", service, "")
		}
		if len(tgt.AccountID) != 12 { //nolint:gomnd // meh.
			sl.ReportError(tgt, "AccountID", "AccountID", "account-id", "")
		}
		if !strings.HasPrefix(tgt.Resource, resourcePrefix) {
			sl.ReportError(tgt, "Resource", "Resource", resourcePrefix, "")
		}
	}
}
