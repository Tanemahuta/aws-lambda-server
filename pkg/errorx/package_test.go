package errorx_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestErrorx(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "errorx Suite")
}
