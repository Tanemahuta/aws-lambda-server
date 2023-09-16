package testlogr_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestLogr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "testing/testlogr suite")
}
