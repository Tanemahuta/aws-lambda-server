package config_test

import (
	"encoding/json"
	"time"

	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Duration", func() {
	var sut *config.Duration
	BeforeEach(func() {
		sut = &config.Duration{}
	})
	Context("UnmarshalJSON()", func() {
		It("should unmarshal correctly", func() {
			Expect(json.Unmarshal([]byte(`"123s"`), sut)).NotTo(HaveOccurred())
			Expect(sut.Duration).To(Equal(time.Second * 123))
		})
		It("should error on non-string", func() {
			Expect(json.Unmarshal([]byte(`1`), sut)).To(MatchError(ContainSubstring("number")))
		})
		It("should error on invalid duration", func() {
			Expect(json.Unmarshal([]byte(`"hoob"`), sut)).To(MatchError(ContainSubstring("invalid duration")))
		})
	})
	Context("UnmarshalYAML()", func() {
		It("should unmarshal correctly", func() {
			Expect(yaml.Unmarshal([]byte(`123s`), sut)).NotTo(HaveOccurred())
			Expect(sut.Duration).To(Equal(time.Second * 123))
		})
		It("should error on invalid duration", func() {
			Expect(yaml.Unmarshal([]byte(`"hoob"`), sut)).To(MatchError(ContainSubstring("invalid duration")))
		})
	})
})