package mux_test

import (
	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/pkg/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("New()", func() {
	var cfg *config.Server
	BeforeEach(func() {
		var err error
		cfg, err = config.Read("../config/testdata/config.yaml")
		Expect(err).NotTo(HaveOccurred())
	})
	It("should compile example router", func() {
		Expect(mux.New(nil, cfg.Functions)).NotTo(BeNil())
	})
})
