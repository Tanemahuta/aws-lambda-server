package errorx_test

import (
	"github.com/Tanemahuta/aws-lambda-server/pkg/errorx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("Fns", func() {
	Context("Run()", func() {
		It("should not error if no Fn errors", func() {
			Expect(
				errorx.Fns{
					func() error { return nil },
					func() error { return nil },
				}.Run(),
			).NotTo(HaveOccurred())
		})
		It("should return the first error", func() {
			Expect(
				errorx.Fns{
					func() error { return errors.New("bla") },
					func() error { return nil },
				}.Run(),
			).To(MatchError("bla"))
		})
	})
})
