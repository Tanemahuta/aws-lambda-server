package testcontext_test

import (
	"github.com/Tanemahuta/aws-lambda-server/testing/testcontext"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("New()", func() {
	It("should create a new Context with logger", func() {
		ctx := testcontext.New()
		Expect(ctx).NotTo(BeNil())
		log, err := logr.FromContext(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(log.V(1).Enabled()).To(BeTrue())
	})
})
