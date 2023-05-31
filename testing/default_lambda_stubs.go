package testing

import (
	_ "embed"
)

//go:embed testdata/lambda-stubs.yaml
var lambdaStubsData []byte

func DefaultLambdaStubs() LambdaStubs {
	return NewLambdaStubs(lambdaStubsData)
}
