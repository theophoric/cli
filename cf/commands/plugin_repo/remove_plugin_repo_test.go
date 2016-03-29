package plugin_repo_test

import (
	testcmd "github.com/theophoric/cf-cli/testhelpers/commands"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
	testreq "github.com/theophoric/cf-cli/testhelpers/requirements"
	testterm "github.com/theophoric/cf-cli/testhelpers/terminal"

	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/models"
	. "github.com/theophoric/cf-cli/testhelpers/matchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("delte-plugin-repo", func() {
	var (
		ui                  *testterm.FakeUI
		config              core_config.Repository
		requirementsFactory *testreq.FakeReqFactory
		deps                command_registry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.Config = config
		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("remove-plugin-repo").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		requirementsFactory = &testreq.FakeReqFactory{}
		config = testconfig.NewRepositoryWithDefaults()
	})

	var callRemovePluginRepo = func(args ...string) bool {
		return testcmd.RunCliCommand("remove-plugin-repo", args, requirementsFactory, updateCommandDependency, false)
	}

	Context("When repo name is valid", func() {
		BeforeEach(func() {
			config.SetPluginRepo(models.PluginRepo{
				Name: "repo1",
				Url:  "http://someserver1.com:1234",
			})

			config.SetPluginRepo(models.PluginRepo{
				Name: "repo2",
				Url:  "http://server2.org:8080",
			})
		})

		It("deletes the repo from the config", func() {
			callRemovePluginRepo("repo1")
			Expect(len(config.PluginRepos())).To(Equal(1))
			Expect(config.PluginRepos()[0].Name).To(Equal("repo2"))
			Expect(config.PluginRepos()[0].Url).To(Equal("http://server2.org:8080"))
		})
	})

	Context("When named repo doesn't exist", func() {
		BeforeEach(func() {
			config.SetPluginRepo(models.PluginRepo{
				Name: "repo1",
				Url:  "http://someserver1.com:1234",
			})

			config.SetPluginRepo(models.PluginRepo{
				Name: "repo2",
				Url:  "http://server2.org:8080",
			})
		})

		It("doesn't change the config the config", func() {
			callRemovePluginRepo("fake-repo")

			Expect(len(config.PluginRepos())).To(Equal(2))
			Expect(config.PluginRepos()[0].Name).To(Equal("repo1"))
			Expect(config.PluginRepos()[0].Url).To(Equal("http://someserver1.com:1234"))
			Expect(ui.Outputs).To(ContainSubstrings([]string{"fake-repo", "does not exist as a repo"}))
		})
	})
})
