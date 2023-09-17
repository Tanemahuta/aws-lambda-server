package routing_test

import (
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws/lambda"
	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/pkg/routing"
	"github.com/Tanemahuta/aws-lambda-server/testing"
	"github.com/Tanemahuta/aws-lambda-server/testing/testcontext"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("New()", func() {
	var (
		stubs testing.LambdaStubs
		cfg   *config.Server
	)
	BeforeEach(func() {
		var err error
		cfg, err = config.Read(testcontext.New(), "../config/testdata/config.yaml")
		cfg.DisableValidation = false
		Expect(err).NotTo(HaveOccurred())
		stubs = testing.DefaultLambdaStubs()
	})
	It("should compile example router", func() {
		decoratorInvoked := false
		Expect(routing.New(stubs, cfg, func(handler http.Handler, _ lambda.FnRef) http.Handler {
			decoratorInvoked = true
			return handler
		})).NotTo(BeNil())
		Expect(decoratorInvoked).To(BeTrue())
	})
	It("should error in case the lambda cannot be invoked", func() {
		delete(stubs, cfg.Functions[0].Name)
		handler, err := routing.New(stubs, cfg)
		Expect(handler).To(BeNil())
		Expect(err).To(HaveOccurred())
	})
	It("should error in case compilation fails", func() {
		handler, err := routing.New(stubs, &config.Server{
			Functions: []config.Function{{Name: cfg.Functions[0].Name, Routes: []config.Route{
				{Methods: []string{http.MethodGet}, Path: "/{"},
			}}},
		})
		Expect(handler).To(BeNil())
		Expect(err).To(HaveOccurred())
	})
})
