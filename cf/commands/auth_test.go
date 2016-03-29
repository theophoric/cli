package commands_test

import (
	"errors"

	"github.com/theophoric/cf-cli/cf"
	authenticationfakes "github.com/theophoric/cf-cli/cf/api/authentication/fakes"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/trace/fakes"
	testcmd "github.com/theophoric/cf-cli/testhelpers/commands"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
	testreq "github.com/theophoric/cf-cli/testhelpers/requirements"
	testterm "github.com/theophoric/cf-cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/theophoric/cf-cli/testhelpers/matchers"
)

var _ = Describe("auth command", func() {
	var (
		ui                  *testterm.FakeUI
		config              core_config.Repository
		authRepo            *authenticationfakes.FakeAuthenticationRepository
		requirementsFactory *testreq.FakeReqFactory
		deps                command_registry.Dependency
		fakeLogger          *fakes.FakePrinter
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.Config = config
		deps.RepoLocator = deps.RepoLocator.SetAuthenticationRepository(authRepo)
		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("auth").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		config = testconfig.NewRepositoryWithDefaults()
		requirementsFactory = &testreq.FakeReqFactory{}
		authRepo = &authenticationfakes.FakeAuthenticationRepository{}
		authRepo.AuthenticateStub = func(credentials map[string]string) error {
			config.SetAccessToken("my-access-token")
			config.SetRefreshToken("my-refresh-token")
			return nil
		}

		fakeLogger = new(fakes.FakePrinter)
		deps = command_registry.NewDependency(fakeLogger)
	})

	Describe("requirements", func() {
		It("fails with usage when given too few arguments", func() {
			testcmd.RunCliCommand("auth", []string{}, requirementsFactory, updateCommandDependency, false)

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Incorrect Usage", "Requires", "arguments"},
			))
		})

		It("fails if the user has not set an api endpoint", func() {
			Expect(testcmd.RunCliCommand("auth", []string{"username", "password"}, requirementsFactory, updateCommandDependency, false)).To(BeFalse())
		})
	})

	Context("when an api endpoint is targeted", func() {
		BeforeEach(func() {
			requirementsFactory.ApiEndpointSuccess = true
			config.SetApiEndpoint("foo.example.org/authenticate")
		})

		It("authenticates successfully", func() {
			requirementsFactory.ApiEndpointSuccess = true
			testcmd.RunCliCommand("auth", []string{"foo@example.com", "password"}, requirementsFactory, updateCommandDependency, false)

			Expect(ui.FailedWithUsage).To(BeFalse())
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"foo.example.org/authenticate"},
				[]string{"OK"},
			))

			Expect(authRepo.AuthenticateArgsForCall(0)).To(Equal(map[string]string{
				"username": "foo@example.com",
				"password": "password",
			}))
		})

		It("prompts users to upgrade if CLI version < min cli version requirement", func() {
			config.SetMinCliVersion("5.0.0")
			config.SetMinRecommendedCliVersion("5.5.0")
			cf.Version = "4.5.0"

			testcmd.RunCliCommand("auth", []string{"foo@example.com", "password"}, requirementsFactory, updateCommandDependency, false)

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"To upgrade your CLI"},
				[]string{"5.0.0"},
			))
		})

		It("gets the UAA endpoint and saves it to the config file", func() {
			requirementsFactory.ApiEndpointSuccess = true
			testcmd.RunCliCommand("auth", []string{"foo@example.com", "password"}, requirementsFactory, updateCommandDependency, false)
			Expect(authRepo.GetLoginPromptsAndSaveUAAServerURLCallCount()).To(Equal(1))
		})

		Describe("when authentication fails", func() {
			BeforeEach(func() {
				authRepo.AuthenticateReturns(errors.New("Error authenticating."))
				testcmd.RunCliCommand("auth", []string{"username", "password"}, requirementsFactory, updateCommandDependency, false)
			})

			It("does not prompt the user when provided username and password", func() {
				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{config.ApiEndpoint()},
					[]string{"Authenticating..."},
					[]string{"FAILED"},
					[]string{"Error authenticating"},
				))
			})

			It("clears the user's session", func() {
				Expect(config.AccessToken()).To(BeEmpty())
				Expect(config.RefreshToken()).To(BeEmpty())
				Expect(config.SpaceFields()).To(Equal(models.SpaceFields{}))
				Expect(config.OrganizationFields()).To(Equal(models.OrganizationFields{}))
			})
		})
	})
})
