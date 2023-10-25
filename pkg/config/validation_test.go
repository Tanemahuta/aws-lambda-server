package config_test

import (
	"path"

	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/testing/testcontext"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate()", func() {
	Context("Server", func() {
		It("should not fail on minimal configurations", func() {
			validArn := config.LambdaARN{
				ARN: config.ARN{ARN: arn.ARN{
					Partition: "aws",
					Service:   "lambda",
					Region:    "eu-central-1",
					AccountID: "123456789012",
					Resource:  "function:my-function",
				},
				}}
			Expect(config.Validate(&config.Server{
				Functions: []config.Function{
					{ARN: &validArn, Routes: []config.Route{{Path: "/"}}},
				},
			})).NotTo(HaveOccurred())
			Expect(config.Validate(&config.Server{
				Functions: []config.Function{
					{ARN: &validArn, Routes: []config.Route{{PathPrefix: "/"}}},
				},
			})).NotTo(HaveOccurred())
		})
		It("should validate example config", func() {
			cfg, err := config.Read(testcontext.New(), path.Join("testdata", "config.yaml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Validate(cfg)).NotTo(HaveOccurred())
		})
		It("should error if functions are empty", func() {
			Expect(config.Validate(&config.Server{})).To(MatchError(ContainSubstring("Functions")))
		})
		It("should error on ARN and name not set", func() {
			Expect(config.Validate(&config.Server{
				Functions: []config.Function{
					{
						Routes: []config.Route{
							{PathPrefix: "/"},
						},
					},
				},
			})).To(MatchError(ContainSubstring("ARN")))
		})
		It("should error on invalid ARN", func() {
			Expect(config.Validate(&config.Server{
				Functions: []config.Function{
					{
						ARN: &config.LambdaARN{ARN: config.ARN{ARN: arn.ARN{
							Partition: "aws",
							Service:   "lambda",
							Region:    "eu-central-1",
							AccountID: "123456789012",
							Resource:  "layer:my-layer",
						}}},
						Routes: []config.Route{
							{PathPrefix: "/"},
						},
					},
				},
			})).To(MatchError(ContainSubstring("ARN")))
		})
		It("should error on Path,PathPrefix not set", func() {
			Expect(config.Validate(&config.Server{
				Functions: []config.Function{
					{
						ARN:    &config.LambdaARN{ARN: config.ARN{ARN: arn.ARN{Region: "eu-central-1"}}},
						Routes: []config.Route{{}},
					},
				},
			})).To(MatchError(And(
				ContainSubstring("Path"),
				ContainSubstring("PathPrefix"),
			)))
		})
	})
})
