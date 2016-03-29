package plugin_repo_test

import (
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/models"

	testcmd "github.com/theophoric/cf-cli/testhelpers/commands"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
	testreq "github.com/theophoric/cf-cli/testhelpers/requirements"
	testterm "github.com/theophoric/cf-cli/testhelpers/terminal"

	. "github.com/theophoric/cf-cli/testhelpers/matchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("list-plugin-repo", func() {
	var (
		ui                  *testterm.FakeUI
		config              core_config.Repository
		requirementsFactory *testreq.FakeReqFactory
		deps                command_registry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.Config = config
		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("list-plugin-repos").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		requirementsFactory = &testreq.FakeReqFactory{}
		config = testconfig.NewRepositoryWithDefaults()
	})

	var callListPluginRepos = func(args ...string) bool {
		return testcmd.RunCliCommand("list-plugin-repos", args, requirementsFactory, updateCommandDependency, false)
	}

	It("lists all added plugin repo in a table", func() {
		config.SetPluginRepo(models.PluginRepo{
			Name: "repo1",
			Url:  "http://url1.com",
		})
		config.SetPluginRepo(models.PluginRepo{
			Name: "repo2",
			Url:  "http://url2.com",
		})

		callListPluginRepos()

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"repo1", "http://url1.com"},
			[]string{"repo2", "http://url2.com"},
		))

	})

})
