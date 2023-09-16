package config

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
)

// Function providing the lambda function and routes.
type Function struct {
	// Name of the function.
	Name string `json:"name,omitempty" yaml:"name,omitempty" validate:"required_without=ARN,excluded_with=ARN"`
	// ARN of the function to be invoked.
	// Deprecated: use Name.
	ARN *LambdaARN `json:"arn,omitempty" yaml:"arn,omitempty" validate:"required_without=Name,excluded_with=Name"`
	// InvocationRole
	InvocationRole *RoleARN `json:"invocationRole,omitempty" yaml:"invocationRole,omitempty" validate:"omitempty"`
	// Routes to be added for that function.
	Routes []Route `json:"routes" yaml:"routes" validate:"min=1,dive"`
}

// GetName of the function.
func (f *Function) GetName() string {
	switch {
	case len(f.Name) > 0:
		break
	case f.ARN != nil:
		return f.ARN.String()
	}
	return f.Name
}

// GetInvocationRoleARN of the function.
func (f *Function) GetInvocationRoleARN() *arn.ARN {
	if f.InvocationRole == nil {
		return nil
	}
	return &f.InvocationRole.ARN.ARN
}
