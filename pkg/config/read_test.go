package config_test

import (
	"net/http"
	"time"

	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/testing/testcontext"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Read()", func() {
	It("should read YAML correctly", func() {
		Expect(config.Read(testcontext.New(), "testdata/config.yaml")).To(Equal(&config.Server{
			HTTP: config.HTTP{
				ReadHeaderTimeout: config.Duration{Duration: time.Minute * 1},
				ReadTimeout:       config.Duration{Duration: time.Minute * 2},
				WriteTimeout:      config.Duration{Duration: time.Minute * 3},
				EnableTraceparent: true,
			},
			DisableValidation: true,
			Functions: []config.Function{
				{
					Name: "test-function",
					InvocationRole: &config.RoleARN{ARN: config.ARN{ARN: arn.ARN{
						Partition: "aws",
						Service:   "iam",
						AccountID: "123456789012",
						Resource:  "role/test-role",
					}}},
					Routes: []config.Route{
						{
							Methods: []string{http.MethodPost}, Path: "/test",
							Headers: map[string]string{}, HeadersRegexp: map[string]string{},
						},
					},
					Timeout: config.Duration{Duration: time.Minute * 4},
				},
			},
		}))
	})
	It("should read deprecated YAML correctly", func() {
		Expect(config.Read(testcontext.New(), "testdata/deprecated-config.yaml")).To(Equal(&config.Server{
			Functions: []config.Function{
				{
					ARN: &config.LambdaARN{ARN: config.ARN{ARN: arn.ARN{
						Partition: "aws",
						Service:   "lambda",
						Region:    "eu-central-1",
						AccountID: "123456789012",
						Resource:  "function:my-function",
					}}},
					Routes: []config.Route{
						{
							Methods: []string{http.MethodPost}, Path: "/test",
							Headers: map[string]string{}, HeadersRegexp: map[string]string{},
						},
					},
				},
			},
		}))
	})
	It("should read JSON correctly", func() {
		Expect(config.Read(testcontext.New(), "testdata/config.json")).To(Equal(&config.Server{
			Functions: []config.Function{
				{
					Name: "test-function",
					Routes: []config.Route{
						{Methods: []string{http.MethodPost}, Path: "/test"},
					},
				},
			},
		}))
	})
	It("should error on invalid file", func() {
		_, err := config.Read(testcontext.New(), "testdata/config.txt")
		Expect(err).To(MatchError(ContainSubstring(".txt")))
	})
	It("should error on unknown file", func() {
		_, err := config.Read(testcontext.New(), "testdata/config.bla")
		Expect(err).To(HaveOccurred())
	})
})
