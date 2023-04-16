package aws_test

import (
	"encoding/json"
	"strconv"

	"github.com/Tanemahuta/aws-lambda-server/pkg/aws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Body", func() {
	var sut *aws.Body
	BeforeEach(func() {
		sut = &aws.Body{}
	})
	Context("UnmarshalJSON", func() {
		It("should unmarshal byte slice", func() {
			Expect(json.Unmarshal([]byte(`[116,101,115,116]`), sut)).NotTo(HaveOccurred())
			Expect(sut.Formatted).To(BeFalse())
			Expect(sut.String()).To(Equal("test"))
		})
		It("should unmarshal string", func() {
			Expect(json.Unmarshal([]byte(`"test"`), sut)).NotTo(HaveOccurred())
			Expect(sut.Formatted).To(BeFalse())
			Expect(sut.String()).To(Equal("test"))
		})
		It("should unmarshal formatted", func() {
			Expect(json.Unmarshal([]byte(`[1,"x"]`), sut)).NotTo(HaveOccurred())
			Expect(sut.Formatted).To(BeTrue())
			Expect(sut.String()).To(Equal(`[1,"x"]`))
		})
	})
	Context("MarshalJSON", func() {
		It("should marshal formatted", func() {
			sut.Formatted = true
			sut.Data = ([]byte)(`{"x":"y"}`)
			Expect(json.Marshal(sut)).To(Equal(sut.Data))
		})
		It("should marshal unformatted", func() {
			sut.Formatted = false
			sut.Data = ([]byte)(`{"x":"y"}`)
			Expect(json.Marshal(sut)).To(BeEquivalentTo(strconv.Quote(string(sut.Data))))
		})
	})
	Context("UnmarshalYAML", func() {
		It("should unmarshal byte slice", func() {
			Expect(yaml.Unmarshal([]byte(`[116,101,115,116]`), sut)).NotTo(HaveOccurred())
			Expect(sut.Formatted).To(BeFalse())
			Expect(sut.String()).To(Equal("test"))
		})
		It("should unmarshal string", func() {
			Expect(yaml.Unmarshal([]byte(`"test"`), sut)).NotTo(HaveOccurred())
			Expect(sut.Formatted).To(BeFalse())
			Expect(sut.String()).To(Equal("test"))
		})
		It("should unmarshal formatted", func() {
			Expect(yaml.Unmarshal([]byte(`[1,"x"]`), sut)).NotTo(HaveOccurred())
			Expect(sut.Formatted).To(BeTrue())
			Expect(sut.String()).To(Equal(`[1,"x"]`))
		})
	})
	Context("String()", func() {
		It("should return contents", func() {
			sut.Data = []byte("test")
			Expect(sut.String()).To(Equal("test"))
		})
	})
})
