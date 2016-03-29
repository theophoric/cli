package application_test

import (
	testapi "github.com/theophoric/cf-cli/cf/api/fakes"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/trace/fakes"
	"github.com/theophoric/cf-cli/plugin/models"
	testcmd "github.com/theophoric/cf-cli/testhelpers/commands"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
	testreq "github.com/theophoric/cf-cli/testhelpers/requirements"
	testterm "github.com/theophoric/cf-cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/theophoric/cf-cli/cf/commands/application"
	"github.com/theophoric/cf-cli/flags"
	. "github.com/theophoric/cf-cli/testhelpers/matchers"
)

var _ = Describe("list-apps command", func() {
	var (
		ui                  *testterm.FakeUI
		configRepo          core_config.Repository
		appSummaryRepo      *testapi.FakeAppSummaryRepo
		requirementsFactory *testreq.FakeReqFactory
		deps                command_registry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.Config = configRepo
		deps.RepoLocator = deps.RepoLocator.SetAppSummaryRepository(appSummaryRepo)
		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("apps").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		appSummaryRepo = &testapi.FakeAppSummaryRepo{}
		configRepo = testconfig.NewRepositoryWithDefaults()
		requirementsFactory = &testreq.FakeReqFactory{
			LoginSuccess:         true,
			TargetedSpaceSuccess: true,
		}

		app1Routes := []models.RouteSummary{
			models.RouteSummary{
				Host: "app1",
				Domain: models.DomainFields{
					Name:                   "cfapps.io",
					Shared:                 true,
					OwningOrganizationGuid: "org-123",
					Guid: "domain-guid",
				},
			},
			models.RouteSummary{
				Host: "app1",
				Domain: models.DomainFields{
					Name: "example.com",
				},
			}}

		app2Routes := []models.RouteSummary{
			models.RouteSummary{
				Host:   "app2",
				Domain: models.DomainFields{Name: "cfapps.io"},
			}}

		app := models.Application{}
		app.Name = "Application-1"
		app.Guid = "Application-1-guid"
		app.State = "started"
		app.RunningInstances = 1
		app.InstanceCount = 1
		app.Memory = 512
		app.DiskQuota = 1024
		app.Routes = app1Routes
		app.AppPorts = []int{8080, 9090}

		app2 := models.Application{}
		app2.Name = "Application-2"
		app2.Guid = "Application-2-guid"
		app2.State = "started"
		app2.RunningInstances = 1
		app2.InstanceCount = 2
		app2.Memory = 256
		app2.DiskQuota = 1024
		app2.Routes = app2Routes

		appSummaryRepo.GetSummariesInCurrentSpaceApps = []models.Application{app, app2}

		deps = command_registry.NewDependency(new(fakes.FakePrinter))
	})

	runCommand := func(args ...string) bool {
		return testcmd.RunCliCommand("apps", args, requirementsFactory, updateCommandDependency, false)
	}

	Describe("requirements", func() {
		var cmd command_registry.Command
		var flagContext flags.FlagContext

		BeforeEach(func() {
			cmd = &application.ListApps{}
			cmd.SetDependency(deps, false)
			flagContext = flags.NewFlagContext(cmd.MetaData().Flags)

		})

		It("requires the user to be logged in", func() {
			requirementsFactory.LoginSuccess = false
			reqs := cmd.Requirements(requirementsFactory, flagContext)

			Expect(testcmd.RunRequirements(reqs)).To(HaveOccurred())
		})

		It("requires the user to have a space targeted", func() {
			requirementsFactory.TargetedSpaceSuccess = false
			reqs := cmd.Requirements(requirementsFactory, flagContext)

			Expect(testcmd.RunRequirements(reqs)).To(HaveOccurred())
		})

		It("should fail with usage when provided any arguments", func() {
			requirementsFactory.LoginSuccess = true
			requirementsFactory.TargetedSpaceSuccess = true

			flagContext.Parse("blahblah")

			reqs := cmd.Requirements(requirementsFactory, flagContext)

			err := testcmd.RunRequirements(reqs)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage"))
			Expect(err.Error()).To(ContainSubstring("No argument required"))
		})

		It("succeeds with all", func() {
			requirementsFactory.LoginSuccess = true
			requirementsFactory.TargetedSpaceSuccess = true

			reqs := cmd.Requirements(requirementsFactory, flagContext)

			Expect(testcmd.RunRequirements(reqs)).NotTo(HaveOccurred())
		})
	})

	Describe("when invoked by a plugin", func() {
		var (
			pluginAppModels []plugin_models.GetAppsModel
		)

		BeforeEach(func() {
			pluginAppModels = []plugin_models.GetAppsModel{}
			deps.PluginModels.AppsSummary = &pluginAppModels
		})

		It("populates the plugin models upon execution", func() {
			testcmd.RunCliCommand("apps", []string{}, requirementsFactory, updateCommandDependency, true)

			Expect(pluginAppModels[0].Name).To(Equal("Application-1"))
			Expect(pluginAppModels[0].Guid).To(Equal("Application-1-guid"))
			Expect(pluginAppModels[1].Name).To(Equal("Application-2"))
			Expect(pluginAppModels[1].Guid).To(Equal("Application-2-guid"))
			Expect(pluginAppModels[0].State).To(Equal("started"))
			Expect(pluginAppModels[0].TotalInstances).To(Equal(1))
			Expect(pluginAppModels[0].RunningInstances).To(Equal(1))
			Expect(pluginAppModels[0].Memory).To(Equal(int64(512)))
			Expect(pluginAppModels[0].DiskQuota).To(Equal(int64(1024)))
			Expect(pluginAppModels[0].AppPorts).To(Equal([]int{8080, 9090}))
			Expect(pluginAppModels[0].Routes[0].Host).To(Equal("app1"))
			Expect(pluginAppModels[0].Routes[1].Host).To(Equal("app1"))
			Expect(pluginAppModels[0].Routes[0].Domain.Name).To(Equal("cfapps.io"))
			Expect(pluginAppModels[0].Routes[0].Domain.Shared).To(BeTrue())
			Expect(pluginAppModels[0].Routes[0].Domain.OwningOrganizationGuid).To(Equal("org-123"))
			Expect(pluginAppModels[0].Routes[0].Domain.Guid).To(Equal("domain-guid"))
		})
	})

	Context("when the user is logged in and a space is targeted", func() {
		It("lists apps in a table", func() {
			runCommand()

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Getting apps in", "my-org", "my-space", "my-user"},
				[]string{"OK"},
				[]string{"name", "requested state", "instances", "memory", "disk", "app ports", "urls"},
				[]string{"Application-1", "started", "1/1", "512M", "1G", "8080, 9090", "app1.cfapps.io", "app1.example.com"},
				[]string{"Application-2", "started", "1/2", "256M", "1G", "app2.cfapps.io"},
			))
		})

		Context("when an app's running instances is unknown", func() {
			It("dipslays a '?' for running instances", func() {
				appRoutes := []models.RouteSummary{
					models.RouteSummary{
						Host:   "app1",
						Domain: models.DomainFields{Name: "cfapps.io"},
					}}
				app := models.Application{}
				app.Name = "Application-1"
				app.Guid = "Application-1-guid"
				app.State = "started"
				app.RunningInstances = -1
				app.InstanceCount = 2
				app.Memory = 512
				app.DiskQuota = 1024
				app.Routes = appRoutes

				appSummaryRepo.GetSummariesInCurrentSpaceApps = []models.Application{app}

				runCommand()

				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"Getting apps in", "my-org", "my-space", "my-user"},
					[]string{"OK"},
					[]string{"Application-1", "started", "?/2", "512M", "1G", "app1.cfapps.io"},
				))
			})
		})

		Context("when there are no apps", func() {
			It("tells the user that there are no apps", func() {
				appSummaryRepo.GetSummariesInCurrentSpaceApps = []models.Application{}

				runCommand()
				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"Getting apps in", "my-org", "my-space", "my-user"},
					[]string{"OK"},
					[]string{"No apps found"},
				))
			})
		})
	})
})
