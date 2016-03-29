package environmentvariablegroup_test

import (
	test_environmentVariableGroups "github.com/theophoric/cf-cli/cf/api/environment_variable_groups/fakes"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	cf_errors "github.com/theophoric/cf-cli/cf/errors"
	testcmd "github.com/theophoric/cf-cli/testhelpers/commands"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
	. "github.com/theophoric/cf-cli/testhelpers/matchers"
	testreq "github.com/theophoric/cf-cli/testhelpers/requirements"
	testterm "github.com/theophoric/cf-cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("set-staging-environment-variable-group command", func() {
	var (
		ui                           *testterm.FakeUI
		requirementsFactory          *testreq.FakeReqFactory
		environmentVariableGroupRepo *test_environmentVariableGroups.FakeEnvironmentVariableGroupsRepository
		configRepo                   core_config.Repository
		deps                         command_registry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.RepoLocator = deps.RepoLocator.SetEnvironmentVariableGroupsRepository(environmentVariableGroupRepo)
		deps.Config = configRepo
		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("set-staging-environment-variable-group").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		configRepo = testconfig.NewRepositoryWithDefaults()
		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true}
		environmentVariableGroupRepo = &test_environmentVariableGroups.FakeEnvironmentVariableGroupsRepository{}
	})

	runCommand := func(args ...string) bool {
		return testcmd.RunCliCommand("set-staging-environment-variable-group", args, requirementsFactory, updateCommandDependency, false)
	}

	Describe("requirements", func() {
		It("requires the user to be logged in", func() {
			requirementsFactory.LoginSuccess = false
			Expect(runCommand()).ToNot(HavePassedRequirements())
		})

		It("fails with usage when it does not receive any arguments", func() {
			requirementsFactory.LoginSuccess = true
			runCommand()
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Incorrect Usage", "Requires an argument"},
			))
		})
	})

	Describe("when logged in", func() {
		BeforeEach(func() {
			requirementsFactory.LoginSuccess = true
		})

		It("Sets the staging environment variable group", func() {
			runCommand(`{"abc":"123", "def": "456"}`)

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Setting the contents of the staging environment variable group as my-user..."},
				[]string{"OK"},
			))
			Expect(environmentVariableGroupRepo.SetStagingArgsForCall(0)).To(Equal(`{"abc":"123", "def": "456"}`))
		})

		It("Fails with a reasonable message when invalid JSON is passed", func() {
			environmentVariableGroupRepo.SetStagingReturns(cf_errors.NewHttpError(400, cf_errors.MessageParseError, "Request invalid due to parse error"))
			runCommand(`{"abc":"123", "invalid : "json"}`)
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Setting the contents of the staging environment variable group as my-user..."},
				[]string{"FAILED"},
				[]string{`Your JSON string syntax is invalid.  Proper syntax is this:  cf set-staging-environment-variable-group '{"name":"value","name":"value"}'`},
			))
		})
	})
})
