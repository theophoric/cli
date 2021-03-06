package net_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/errors"
	. "github.com/theophoric/cf-cli/cf/net"
	"github.com/theophoric/cf-cli/cf/trace/fakes"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
	testterm "github.com/theophoric/cf-cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var failingUAARequest = func(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusBadRequest)
	jsonResponse := `{ "error": "foo", "error_description": "The foo is wrong..." }`
	fmt.Fprintln(writer, jsonResponse)
}

var _ = Describe("UAA Gateway", func() {
	var gateway Gateway
	var config core_config.Reader

	BeforeEach(func() {
		config = testconfig.NewRepository()
		gateway = NewUAAGateway(config, &testterm.FakeUI{}, new(fakes.FakePrinter))
	})

	It("parses error responses", func() {
		ts := httptest.NewTLSServer(http.HandlerFunc(failingUAARequest))
		defer ts.Close()
		gateway.SetTrustedCerts(ts.TLS.Certificates)

		request, apiErr := gateway.NewRequest("GET", ts.URL, "TOKEN", nil)
		_, apiErr = gateway.PerformRequest(request)

		Expect(apiErr).NotTo(BeNil())
		Expect(apiErr.Error()).To(ContainSubstring("The foo is wrong"))
		Expect(apiErr.(errors.HttpError).ErrorCode()).To(ContainSubstring("foo"))
	})
})
