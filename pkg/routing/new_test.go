package routing_test

import (
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/pkg/routing"
	"github.com/Tanemahuta/aws-lambda-server/testing"
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
		cfg, err = config.Read("../config/testdata/config.yaml")
		Expect(err).NotTo(HaveOccurred())
		stubs = testing.DefaultLambdaStubs()
	})
	It("should compile example router", func() {
		decoratorInvoked := false
		Expect(routing.New(stubs, cfg.Functions, func(handler http.Handler, _ string) http.Handler {
			decoratorInvoked = true
			return handler
		})).NotTo(BeNil())
		Expect(decoratorInvoked).To(BeTrue())
	})
	It("should error in case the lambda cannot be invoked", func() {
		delete(stubs, cfg.Functions[0].ARN.String())
		handler, err := routing.New(stubs, cfg.Functions)
		Expect(handler).To(BeNil())
		Expect(err).To(HaveOccurred())
	})
	It("should error in case compilation fails", func() {
		handler, err := routing.New(stubs, []config.Function{
			{ARN: cfg.Functions[0].ARN, Routes: []config.Route{
				{Methods: []string{http.MethodGet}, Path: "/{"},
			}},
		})
		Expect(handler).To(BeNil())
		Expect(err).To(HaveOccurred())
	})
})
