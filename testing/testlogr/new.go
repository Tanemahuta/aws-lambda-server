package testlogr

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/onsi/ginkgo/v2"
	"go.uber.org/zap/zaptest"
)

// New returns a new runtime singleton logr.Logger.
func New() logr.Logger {
	return zapr.NewLogger(zaptest.NewLogger(ginkgo.GinkgoT()))
}
