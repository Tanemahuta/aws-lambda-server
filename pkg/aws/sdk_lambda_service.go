package aws

import (
	"context"
	"encoding/json"

	"github.com/Tanemahuta/aws-lambda-server/pkg/errorx"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

type SdkLambdaService func(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (
	*lambda.InvokeOutput, error,
)

func (s SdkLambdaService) Invoke(ctx context.Context, arn arn.ARN, request *LambdaRequest) (*LambdaResponse, error) {
	var (
		payload  []byte
		response *lambda.InvokeOutput
		result   *LambdaResponse
		err      error
	)
	log := logr.FromContextOrDiscard(ctx)
	err = errorx.Fns{
		func() error {
			payload, err = json.Marshal(request)
			return errors.Wrap(err, "could not marshal request")
		},
		func() error {
			functionName := arn.String()
			log.Info("invoking lambda function", "arn", arn)
			response, err = s(ctx, &lambda.InvokeInput{
				FunctionName:   &functionName,
				InvocationType: types.InvocationTypeRequestResponse,
				LogType:        types.LogTypeTail,
				Payload:        payload,
			})
			s.handleResponse(response)
			return errors.Wrapf(err, "could not invoke lambda %v", arn)
		},
		func() error {
			result, err = s.adaptResponse(log.V(1), response)
			return errors.Wrap(err, "could not adapt response")
		},
	}.Run()
	if err != nil {
		log.Error(err, "lambda invocation failed", "response", response)
	}
	return result, err
}

func (s SdkLambdaService) adaptResponse(log logr.Logger, response *lambda.InvokeOutput) (*LambdaResponse, error) {
	log.Info("converting response", "response", response)
	var result LambdaResponse
	if err := json.Unmarshal(response.Payload, &result); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal payload to response")
	}
	return &result, nil
}

func (s SdkLambdaService) handleResponse(response *lambda.InvokeOutput) {
	if response == nil {
		return
	}
	response.LogResult = HandleBase64String(response.LogResult)
	response.Payload = HandleBase64(response.Payload)
}

// NewLambdaService from aws-sdk.
func NewLambdaService(ctx context.Context) (LambdaService, error) {
	var (
		cfg    aws.Config
		result LambdaService
		err    error
	)
	err = errorx.Fns{
		func() error {
			cfg, err = config.LoadDefaultConfig(ctx)
			return err
		},
		func() error {
			result = SdkLambdaService(lambda.NewFromConfig(cfg).Invoke)
			return nil
		},
	}.Run()
	return result, err
}
