package server_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws"
	"github.com/Tanemahuta/aws-lambda-server/pkg/server"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Run()", func() {
	var (
		httpServer   *httptest.Server
		serverConfig server.Config
	)
	BeforeEach(func() {
		var lambdaStubs LambdaStubs
		data, err := os.ReadFile("testdata/lambda-stubs.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(yaml.Unmarshal(data, &lambdaStubs)).NotTo(HaveOccurred())
		serverConfig = server.Config{
			Filename: "../config/testdata/config.yaml",
			Listen:   ":8080",
			LambdaServiceFactory: func(context.Context) (aws.LambdaService, error) {
				return lambdaStubs, nil
			},
			RunFunc: func(ctx context.Context, _ string, handler http.Handler) error {
				httpServer = httptest.NewServer(handler)
				for key, stubs := range lambdaStubs {
					for idx := range stubs {
						stubs[idx].Request.Host = strings.TrimPrefix(httpServer.URL, "http://")
					}
					lambdaStubs[key] = stubs
				}
				return nil
			},
		}
	})
	AfterEach(func() {
		if httpServer != nil {
			httpServer.Close()
		}
	})
	It("should run correctly", func() {
		Expect(server.Run(context.Background(), serverConfig)).NotTo(HaveOccurred())
		response, err := http.Post(httpServer.URL+"/test", "text/plain", bytes.NewBufferString("test"))
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusAccepted))
		Expect(response.Header).To(HaveKeyWithValue("Test", []string{"test"}))
		Expect(io.ReadAll(response.Body)).To(BeEquivalentTo([]byte("test")))
	})
	It("should error from config.Read", func() {
		serverConfig.Filename += "2"
		Expect(server.Run(context.Background(), serverConfig)).To(MatchError(ContainSubstring("no such file or directory")))
	})
	It("should error from lambda factory", func() {
		serverConfig.LambdaServiceFactory = func(_ context.Context) (aws.LambdaService, error) {
			return nil, errors.New("meh")
		}
		Expect(server.Run(context.Background(), serverConfig)).To(MatchError(ContainSubstring("meh")))
	})
})

type LambdaStub struct {
	Request  aws.LambdaRequest  `json:"requests"`
	Response aws.LambdaResponse `json:"response"`
}

type LambdaStubs map[string][]LambdaStub

func (l LambdaStubs) Invoke(_ context.Context, arn arn.ARN, request *aws.LambdaRequest) (*aws.LambdaResponse, error) {
	for _, stub := range l[arn.String()] {
		if reflect.DeepEqual(&stub.Request, request) {
			return &stub.Response, nil
		}
	}
	defer GinkgoRecover()
	data, _ := yaml.Marshal(request)
	Fail(fmt.Sprintf("request for '%v' not found:\n%v", arn, string(data)))
	return nil, errors.Errorf("could not find request stub for lambda '%v': %v", arn, request)
}
