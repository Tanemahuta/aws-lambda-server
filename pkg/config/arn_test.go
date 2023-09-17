package config_test

import (
	"encoding/json"

	"gopkg.in/yaml.v3"

	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LambdaARN", func() {
	var sut *config.ARN
	BeforeEach(func() {
		sut = &config.ARN{}
	})
	Context("UnmarshalJSON", func() {
		It("should unmarshal correctly", func() {
			Expect(json.Unmarshal([]byte(`"arn:aws:lambda:us-west-2:123456789012:function:my-function"`), sut)).
				NotTo(HaveOccurred())
			Expect(sut.ARN.Service).To(Equal("lambda"))
			Expect(sut.ARN.Partition).To(Equal("aws"))
			Expect(sut.ARN.AccountID).To(Equal("123456789012"))
			Expect(sut.ARN.Region).To(Equal("us-west-2"))
			Expect(sut.ARN.Resource).To(Equal("function:my-function"))
		})
		It("should error on non-string", func() {
			Expect(json.Unmarshal([]byte(`3`), sut)).To(MatchError(ContainSubstring("string")))
		})
		It("should error on parser error", func() {
			Expect(json.Unmarshal([]byte(`"arn:aws:lambda:us-west-2:xy"`), sut)).
				To(MatchError(ContainSubstring("not enough sections")))
		})
	})
	Context("UnmarshalYAML", func() {
		It("should unmarshal correctly", func() {
			Expect(yaml.Unmarshal([]byte(`"arn:aws:lambda:us-west-2:123456789012:function:my-function"`), sut)).
				NotTo(HaveOccurred())
			Expect(sut.ARN.Service).To(Equal("lambda"))
			Expect(sut.ARN.Partition).To(Equal("aws"))
			Expect(sut.ARN.AccountID).To(Equal("123456789012"))
			Expect(sut.ARN.Region).To(Equal("us-west-2"))
			Expect(sut.ARN.Resource).To(Equal("function:my-function"))
		})
		It("should error on parser error", func() {
			Expect(yaml.Unmarshal([]byte(`"arn:aws:lambda:us-west-2:xy"`), sut)).
				To(MatchError(ContainSubstring("not enough sections")))
		})
	})
})
