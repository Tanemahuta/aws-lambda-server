package config

// LambdaARN type alias for validation.
type LambdaARN struct{ ARN }

// Function providing the lambda function and routes.
type Function struct {
	// ARN of the function to be invoked.
	ARN LambdaARN `json:"arn" yaml:"arn" validate:"required"`
	// Routes to be added for that function.
	Routes []Route `json:"routes" yaml:"routes" validate:"min=1,dive"`
}
