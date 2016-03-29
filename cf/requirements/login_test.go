package requirements_test

import (
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/requirements"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
)

var _ = Describe("LoginRequirement", func() {
	BeforeEach(func() {
	})

	It("succeeds when given a config with an API endpoint and authentication", func() {
		config := testconfig.NewRepositoryWithAccessToken(core_config.TokenInfo{Username: "my-user"})
		config.SetApiEndpoint("api.example.com")
		req := NewLoginRequirement(config)
		err := req.Execute()
		Expect(err).NotTo(HaveOccurred())
	})

	It("fails when given a config with only an API endpoint", func() {
		config := testconfig.NewRepository()
		config.SetApiEndpoint("api.example.com")
		req := NewLoginRequirement(config)
		err := req.Execute()
		Expect(err).To(HaveOccurred())

		Expect(err.Error()).To(ContainSubstring("Not logged in."))
	})

	It("fails when given a config with neither an API endpoint nor authentication", func() {
		config := testconfig.NewRepository()
		req := NewLoginRequirement(config)
		err := req.Execute()
		Expect(err).To(HaveOccurred())

		Expect(err.Error()).To(ContainSubstring("No API endpoint"))
		Expect(err.Error()).ToNot(ContainSubstring("Not logged in."))
	})
})
