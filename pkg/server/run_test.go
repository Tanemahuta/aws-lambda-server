package server_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws"
	"github.com/Tanemahuta/aws-lambda-server/pkg/metrics"
	"github.com/Tanemahuta/aws-lambda-server/pkg/server"
	"github.com/Tanemahuta/aws-lambda-server/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("Run()", func() {
	var (
		lambdaServer, metricsServer *httptest.Server
		serverConfig                server.Config
	)
	BeforeEach(func() {
		lambdaStubs := testing.DefaultLambdaStubs()
		serverConfig = server.Config{
			Filename:      "../config/testdata/config.yaml",
			Listen:        ":8080",
			MetricsListen: ":8081",
			LambdaServiceFactory: func(context.Context) (aws.LambdaService, error) {
				return lambdaStubs, nil
			},
			RunFunc: func(ctx context.Context, listenAddr string, handler http.Handler) error {
				switch listenAddr {
				case serverConfig.Listen:
					lambdaServer = httptest.NewServer(handler)
					for key, stubs := range lambdaStubs {
						for idx := range stubs {
							stubs[idx].Request.Host = strings.TrimPrefix(lambdaServer.URL, "http://")
						}
						lambdaStubs[key] = stubs
					}
				case serverConfig.MetricsListen:
					metricsServer = httptest.NewServer(handler)
				}
				return nil
			},
		}
	})
	AfterEach(func() {
		if lambdaServer != nil {
			lambdaServer.Close()
		}
		if metricsServer != nil {
			lambdaServer.Close()
		}
		metrics.HTTPRequestsTotal.Reset()
		metrics.HTTPRequestsDuration.Reset()
		metrics.HTTPRequestsSize.Reset()
		metrics.HTTPResponsesSize.Reset()
	})
	When("running the server", func() {
		var (
			response *http.Response
			err      error
		)
		BeforeEach(func() {
			Expect(server.Run(context.Background(), serverConfig)).NotTo(HaveOccurred())
		})
		When("handling a valid request", func() {
			BeforeEach(func() {
				response, err = http.Post(lambdaServer.URL+"/test", "text/plain", bytes.NewBufferString("test"))
				Expect(err).NotTo(HaveOccurred())
			})
			It("convert the response", func() {
				Expect(response.StatusCode).To(Equal(http.StatusAccepted))
				Expect(response.Header).To(HaveKeyWithValue("Test", []string{"test"}))
				Expect(io.ReadAll(response.Body)).To(BeEquivalentTo([]byte("test")))
			})
			It("should add metrics", func() {
				Expect(metrics.Collect(metrics.HTTPRequestsTotal)).To(HaveKeyWithValue(
					"code=202,functionArn=arn:aws:lambda:eu-central-1:123456789012:function:my-function,method=post",
					BeNumerically("==", 1),
				))
				Expect(metrics.Collect(metrics.HTTPRequestsDuration)).To(HaveKeyWithValue(
					"code=202,functionArn=arn:aws:lambda:eu-central-1:123456789012:function:my-function,method=post",
					BeNumerically("==", 1),
				))
				Expect(metrics.Collect(metrics.HTTPRequestsSize)).To(HaveKeyWithValue(
					"code=202,functionArn=arn:aws:lambda:eu-central-1:123456789012:function:my-function,method=post",
					BeNumerically("==", 1),
				))
				Expect(metrics.Collect(metrics.HTTPResponsesSize)).To(HaveKeyWithValue(
					"code=202,functionArn=arn:aws:lambda:eu-central-1:123456789012:function:my-function,method=post",
					BeNumerically("==", 1),
				))
			})
		})
		It("should serve metrics", func() {
			response, err = http.Get(metricsServer.URL + "/metrics")
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(io.ReadAll(response.Body)).NotTo(BeEmpty())
		})
		It("should serve readyz", func() {
			response, err = http.Get(metricsServer.URL + "/readyz")
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
		It("should serve healthz", func() {
			response, err = http.Get(metricsServer.URL + "/healthz")
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
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
