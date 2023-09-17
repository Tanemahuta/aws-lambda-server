package testcontext_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestContext(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "testing/testcontext suite")
}
