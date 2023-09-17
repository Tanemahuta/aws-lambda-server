package sdk

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

// LambdaClientProps for the provided aws.Config.
func LambdaClientProps(config aws.Config) ClientProps[Lambda] {
	return ClientProps[Lambda]{
		Config: config,
		NewClient: func(config aws.Config, provider aws.CredentialsProvider) Lambda {
			return lambda.NewFromConfig(config, func(options *lambda.Options) {
				options.Credentials = provider
			})
		},
	}
}
