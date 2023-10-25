package lambda

import (
	"context"
	"encoding/json"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws/sdk"
	usercfg "github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/pkg/errorx"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

var _ Facade = &SdkFacade{}

type SdkFacade struct {
	Clients sdk.AssumeClients[sdk.Lambda]
}

func (s *SdkFacade) CanInvoke(ctx context.Context, ref FnRef) error {
	log := logr.FromContextOrDiscard(ctx)
	log.Info("checking if lambda can be invoked", "ref", ref)
	_, err := s.Clients.Get(ref.RoleARN).Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   &ref.Name,
		LogType:        types.LogTypeNone,
		InvocationType: types.InvocationTypeDryRun,
	})
	return err
}

func (s *SdkFacade) Invoke(ctx context.Context, ref FnRef, req *Request) (*Response, error) {
	var (
		payload  []byte
		response *lambda.InvokeOutput
		result   *Response
		err      error
	)
	log := logr.FromContextOrDiscard(ctx)
	err = errorx.Fns{
		func() error {
			payload, err = json.Marshal(req)
			return errors.Wrap(err, "could not marshal request")
		},
		func() error {
			log.Info("invoking lambda function", "ref", ref)
			response, err = s.handleResponse(s.Clients.Get(ref.RoleARN).Invoke(ctx,
				&lambda.InvokeInput{
					FunctionName:   &ref.Name,
					InvocationType: types.InvocationTypeRequestResponse,
					LogType:        types.LogTypeNone,
					Payload:        payload,
				},
			))
			if err != nil && response != nil && response.LogResult != nil {
				log.Error(err, *response.LogResult)
			}
			return errors.Wrapf(err, "could not invoke lambda %v with role %v", ref.Name, ref.RoleARN)
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

func (s *SdkFacade) adaptResponse(log logr.Logger, response *lambda.InvokeOutput) (*Response, error) {
	log.Info("converting response", "response", response)
	var result Response
	if err := json.Unmarshal(response.Payload, &result); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal payload to response")
	}
	return &result, nil
}

func (s *SdkFacade) handleResponse(response *lambda.InvokeOutput, err error) (*lambda.InvokeOutput, error) {
	if err != nil {
		return nil, err
	}
	response.Payload = HandleBase64(response.Payload)
	if response.FunctionError != nil {
		err = errors.Errorf("error '%v' details '%v'", *response.FunctionError, string(response.Payload))
	}
	return response, err
}

// NewLambdaService from aws-sdk.
func NewLambdaService(ctx context.Context, userCfg *usercfg.AWS) (Facade, error) {
	var (
		cfg    aws.Config
		result Facade
		err    error
	)
	err = errorx.Fns{
		func() error {
			cfg, err = config.LoadDefaultConfig(ctx)
			return err
		},
		func() error {
			return userCfg.Apply(&cfg)
		},
		func() error {
			result = &SdkFacade{Clients: sdk.NewAssumeClients[sdk.Lambda](sdk.LambdaClientProps(cfg))}
			return nil
		},
	}.Run()
	return result, err
}
