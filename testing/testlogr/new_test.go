package testlogr_test

import (
	"github.com/Tanemahuta/aws-lambda-server/testing/testlogr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("New()", func() {
	It("should create a new usable logger with debug", func() {
		Expect(testlogr.New().V(1).Enabled()).To(BeTrue())
	})
})
