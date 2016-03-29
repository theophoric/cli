package password_test

import (
	"net/http"
	"net/http/httptest"

	testapi "github.com/theophoric/cf-cli/cf/api/fakes"
	"github.com/theophoric/cf-cli/testhelpers/cloud_controller_gateway"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
	testnet "github.com/theophoric/cf-cli/testhelpers/net"

	. "github.com/theophoric/cf-cli/cf/api/password"
	. "github.com/theophoric/cf-cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudControllerPasswordRepository", func() {
	It("updates your password", func() {
		req := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
			Method:   "PUT",
			Path:     "/Users/my-user-guid/password",
			Matcher:  testnet.RequestBodyMatcher(`{"password":"new-password","oldPassword":"old-password"}`),
			Response: testnet.TestResponse{Status: http.StatusOK},
		})

		passwordUpdateServer, handler, repo := createPasswordRepo(req)
		defer passwordUpdateServer.Close()

		apiErr := repo.UpdatePassword("old-password", "new-password")
		Expect(handler).To(HaveAllRequestsCalled())
		Expect(apiErr).NotTo(HaveOccurred())
	})
})

func createPasswordRepo(req testnet.TestRequest) (passwordServer *httptest.Server, handler *testnet.TestHandler, repo PasswordRepository) {
	passwordServer, handler = testnet.NewServer([]testnet.TestRequest{req})

	configRepo := testconfig.NewRepositoryWithDefaults()
	configRepo.SetUaaEndpoint(passwordServer.URL)
	gateway := cloud_controller_gateway.NewTestCloudControllerGateway(configRepo)
	repo = NewCloudControllerPasswordRepository(configRepo, gateway)
	return
}
