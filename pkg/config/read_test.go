package config_test

import (
	"net/http"

	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Read()", func() {
	It("should read YAML correctly", func() {
		Expect(config.Read("testdata/config.yaml")).To(Equal(&config.Server{
			Functions: []config.Function{
				{
					ARN: config.LambdaARN{ARN: config.ARN{ARN: arn.ARN{
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
		Expect(config.Read("testdata/config.json")).To(Equal(&config.Server{
			Functions: []config.Function{
				{
					ARN: config.LambdaARN{ARN: config.ARN{ARN: arn.ARN{
						Partition: "aws",
						Service:   "lambda",
						Region:    "eu-central-1",
						AccountID: "123456789012",
						Resource:  "function:my-function",
					}}},
					Routes: []config.Route{
						{Methods: []string{http.MethodPost}, Path: "/test"},
					},
				},
			},
		}))
	})
	It("should error on invalid file", func() {
		_, err := config.Read("testdata/config.txt")
		Expect(err).To(MatchError(ContainSubstring(".txt")))
	})
	It("should error on unknown file", func() {
		_, err := config.Read("testdata/config.bla")
		Expect(err).To(HaveOccurred())
	})
})
