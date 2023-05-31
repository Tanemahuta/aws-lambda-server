package testing

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type LambdaStub struct {
	Request  aws.LambdaRequest  `json:"requests"`
	Response aws.LambdaResponse `json:"response"`
}

var _ aws.LambdaService = LambdaStubs{}

func NewLambdaStubs(yamlData []byte) LambdaStubs {
	var result LambdaStubs
	gomega.Expect(yaml.Unmarshal(yamlData, &result)).NotTo(gomega.HaveOccurred())
	return result
}

type LambdaStubs map[string][]LambdaStub

func (l LambdaStubs) CanInvoke(_ context.Context, arn arn.ARN) error {
	if _, ok := l[arn.String()]; !ok {
		return errors.Errorf("lambda %v not stubbed", arn)
	}
	return nil
}

func (l LambdaStubs) Invoke(_ context.Context, arn arn.ARN, request *aws.LambdaRequest) (*aws.LambdaResponse, error) {
	for _, stub := range l[arn.String()] {
		if reflect.DeepEqual(&stub.Request, request) {
			return &stub.Response, nil
		}
	}
	defer ginkgo.GinkgoRecover()
	data, _ := yaml.Marshal(request)
	ginkgo.Fail(fmt.Sprintf("request for '%v' not found:\n%v", arn, string(data)))
	return nil, errors.Errorf("could not find request stub for lambda '%v': %v", arn, request)
}
