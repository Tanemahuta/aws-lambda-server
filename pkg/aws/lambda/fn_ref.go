package lambda

import "github.com/aws/aws-sdk-go-v2/aws/arn"

// FnRef provides the necessary parameters for invoking a lambda function.
type FnRef struct {
	// Name of the function
	Name string
	// RoleARN to be assumed (optional)
	RoleARN *arn.ARN
}
