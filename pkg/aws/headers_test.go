package aws_test

import (
	_ "embed"
	"encoding/json"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

//go:embed testdata/headers.json
var headerJSONData []byte

//go:embed testdata/headers.yaml
var headerYamlData []byte

var _ = Describe("Headers", func() {
	var sut aws.Headers
	BeforeEach(func() {
		sut = aws.Headers{}
	})
	Context("UnmarshalJSON()", func() {
		It("should unmarshal different types correctly", func() {
			Expect(json.Unmarshal(headerJSONData, &sut)).NotTo(HaveOccurred())
			Expect(sut.Header).To(And(
				HaveKeyWithValue("Bool", ConsistOf("true")),
				HaveKeyWithValue("Double", ConsistOf("0.0")),
				HaveKeyWithValue("Int", ConsistOf("0")),
				HaveKeyWithValue("String", ConsistOf("a")),
			))
		})
		It("should error if not an object", func() {
			Expect(json.Unmarshal([]byte(`["a"]`), &sut)).To(HaveOccurred())
		})
	})
	Context("UnmarshalYAML()", func() {
		It("should unmarshal different types correctly", func() {
			Expect(yaml.Unmarshal(headerYamlData, &sut)).NotTo(HaveOccurred())
			Expect(sut.Header).To(And(
				HaveKeyWithValue("Bool", ConsistOf("true")),
				HaveKeyWithValue("Double", ConsistOf("0.0")),
				HaveKeyWithValue("Int", ConsistOf("0")),
				HaveKeyWithValue("String", ConsistOf("a")),
			))
		})
		It("should error if not an object", func() {
			Expect(json.Unmarshal([]byte(`- a\n- b`), &sut)).To(HaveOccurred())
		})
	})

})
