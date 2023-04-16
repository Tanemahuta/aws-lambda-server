package mux_test

import (
	"github.com/Tanemahuta/aws-lambda-server/pkg/config"
	"github.com/Tanemahuta/aws-lambda-server/pkg/mux"
	gorilla "github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type testConfig struct {
	ignored  string         // not exported
	NoMethod string         // errors => no method found
	GetError string         // errors => method signature does not match
	Methods  []int          // errors => in[1] does not match
	Headers  map[string]int // errors => value type does not match in[1] slice elem type
}

var _ = Describe("ConfigureRoute()", func() {
	var route *gorilla.Route
	BeforeEach(func() {
		route = gorilla.NewRouter().NewRoute()
	})
	Context("RouteConfig", func() {
		It("should skip zero values", func() {
			configured, err := mux.ConfigureRoute(
				route,
				config.Route{
					Path: "/test/{id}",
				},
				(*gorilla.Route).GetError,
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(configured.GetPathTemplate()).To(Equal("/test/{id}"))
		})
		It("should use all values", func() {
			configured, err := mux.ConfigureRoute(
				route,
				config.Route{
					Name:          "test",
					Host:          "example.com",
					Methods:       []string{"GET"},
					PathPrefix:    "/test/",
					Headers:       map[string]string{"a": "b"},
					HeadersRegexp: map[string]string{"c": ".+"},
				},
				(*gorilla.Route).GetError,
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(configured.GetName()).To(Equal("test"))
			Expect(configured.GetHostTemplate()).To(Equal("example.com"))
			Expect(configured.GetMethods()).To(ConsistOf("GET"))
			Expect(configured.GetPathTemplate()).To(Equal("/test/"))
		})
	})
	Context("errors", func() {
		It("should error on invalid method", func() {
			_, err := mux.ConfigureRoute(route, testConfig{ignored: "b", NoMethod: "a"}, (*gorilla.Route).GetError)
			Expect(err).To(MatchError(ContainSubstring("could not find exported config function")))
		})
		It("should error on invalid method", func() {
			_, err := mux.ConfigureRoute(route, testConfig{GetError: "a"}, (*gorilla.Route).GetError)
			Expect(err).To(MatchError(ContainSubstring("expected two in parameters, but got 1")))
		})
		It("should error on invalid method", func() {
			_, err := mux.ConfigureRoute(route, testConfig{Methods: []int{1}}, (*gorilla.Route).GetError)
			Expect(err).To(MatchError(ContainSubstring("could not convert config value to function input")))
		})
		It("should error on invalid method", func() {
			_, err := mux.ConfigureRoute(route, testConfig{Headers: map[string]int{"a": 1}}, (*gorilla.Route).GetError)
			Expect(err).To(MatchError(ContainSubstring("could not convert map to slice")))
		})
	})
})
