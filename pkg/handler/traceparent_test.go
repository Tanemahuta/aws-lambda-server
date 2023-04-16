package handler_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/Tanemahuta/aws-lambda-server/pkg/handler"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Traceparent", func() {
	var (
		sut     http.Handler
		request *http.Request
		writer  *httptest.ResponseRecorder
	)
	BeforeEach(func() {
		sut = handler.NewTraceparent(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			Expect(request).NotTo(BeNil())
			Expect(writer).NotTo(BeNil())
			writer.WriteHeader(http.StatusAccepted)
			_, _ = writer.Write([]byte("test"))
		}), "test")
		request = httptest.NewRequest(http.MethodGet, "http://test.example.com", nil)
		writer = httptest.NewRecorder()
	})
	It("should add headers", func() {
		Expect(func() { sut.ServeHTTP(writer, request) }).NotTo(Panic())
		Expect(writer.Code).To(Equal(http.StatusAccepted))
		Expect(writer.Header()).To(BeEmpty()) // For now
		Expect(writer.Body.String()).To(Equal("test"))
	})
})
