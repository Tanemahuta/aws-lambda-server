package testing

import (
	"context"
	"fmt"
	"strings"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws/lambda"
	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type LambdaStub struct {
	Request  lambda.Request  `json:"requests"`
	Response lambda.Response `json:"response"`
}

var _ lambda.Facade = LambdaStubs{}

func NewLambdaStubs(yamlData []byte) LambdaStubs {
	var result LambdaStubs
	gomega.Expect(yaml.Unmarshal(yamlData, &result)).NotTo(gomega.HaveOccurred())
	return result
}

type LambdaStubs map[string][]LambdaStub

func (l LambdaStubs) CanInvoke(_ context.Context, fnRef lambda.FnRef) error {
	if _, ok := l[fnRef.Name]; !ok {
		return errors.Errorf("lambda %v not stubbed", fnRef.Name)
	}
	return nil
}

func (l LambdaStubs) Invoke(_ context.Context, fnRef lambda.FnRef, request *lambda.Request) (*lambda.Response, error) {
	stubs, ok := l[fnRef.Name]
	if !ok {
		return nil, l.fail("no request for lambda '%v' stubbed", fnRef.Name)
	}
	diffs := make(map[int]dyff.Report)
	for idx := range stubs {
		report := l.requestsMatch(&stubs[idx].Request, request)
		if len(report.Diffs) == 0 {
			return &stubs[idx].Response, nil
		}
		diffs[idx] = report
	}
	defer ginkgo.GinkgoRecover()
	data, _ := yaml.Marshal(request)
	var sb strings.Builder
	for idx, report := range diffs {
		sb.WriteString(fmt.Sprintf("# request %v:\n", idx))
		_ = (&dyff.HumanReport{Report: report, DoNotInspectCerts: true, OmitHeader: true}).WriteReport(&sb)
		sb.WriteString("\n")
	}
	return nil, l.fail("could not find request stub for lambda '%v': %v\n%v",
		fnRef.Name, string(data), sb.String())
}

func (l LambdaStubs) fail(msg string, keyValues ...interface{}) error {
	ginkgo.Fail(fmt.Sprintf(msg, keyValues...))
	return errors.Errorf(msg, keyValues...)
}

func (l LambdaStubs) requestsMatch(lhs, rhs *lambda.Request) dyff.Report {
	lhsFile, rhsFile := l.createDoc(lhs), l.createDoc(rhs)
	report, _ := dyff.CompareInputFiles(lhsFile, rhsFile)
	return report
}

func (l LambdaStubs) createDoc(req *lambda.Request) ytbx.InputFile {
	data, _ := yaml.Marshal(req)
	docs, _ := ytbx.LoadYAMLDocuments(data)
	return ytbx.InputFile{Documents: docs}
}
