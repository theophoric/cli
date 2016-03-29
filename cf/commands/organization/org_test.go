package organization_test

import (
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/trace/fakes"
	"github.com/theophoric/cf-cli/plugin/models"
	testcmd "github.com/theophoric/cf-cli/testhelpers/commands"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
	testreq "github.com/theophoric/cf-cli/testhelpers/requirements"
	testterm "github.com/theophoric/cf-cli/testhelpers/terminal"

	. "github.com/theophoric/cf-cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("org command", func() {
	var (
		ui                  *testterm.FakeUI
		configRepo          core_config.Repository
		requirementsFactory *testreq.FakeReqFactory
		deps                command_registry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.Config = configRepo
		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("org").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		requirementsFactory = &testreq.FakeReqFactory{}
		configRepo = testconfig.NewRepositoryWithDefaults()

		deps = command_registry.NewDependency(new(fakes.FakePrinter))
	})

	runCommand := func(args ...string) bool {
		return testcmd.RunCliCommand("org", args, requirementsFactory, updateCommandDependency, false)
	}

	Describe("requirements", func() {
		It("fails when not logged in", func() {
			Expect(runCommand("whoops")).To(BeFalse())
		})

		It("fails with usage when not provided exactly one arg", func() {
			requirementsFactory.LoginSuccess = true
			runCommand("too", "much")
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Incorrect Usage", "Requires an argument"},
			))
		})
	})

	Context("when logged in, and provided the name of an org", func() {
		BeforeEach(func() {
			developmentSpaceFields := models.SpaceFields{}
			developmentSpaceFields.Name = "development"
			developmentSpaceFields.Guid = "dev-space-guid-1"
			stagingSpaceFields := models.SpaceFields{}
			stagingSpaceFields.Name = "staging"
			stagingSpaceFields.Guid = "staging-space-guid-1"
			domainFields := models.DomainFields{}
			domainFields.Name = "cfapps.io"
			domainFields.Guid = "1111"
			domainFields.OwningOrganizationGuid = "my-org-guid"
			domainFields.Shared = true
			cfAppDomainFields := models.DomainFields{}
			cfAppDomainFields.Name = "cf-app.com"
			cfAppDomainFields.Guid = "2222"
			cfAppDomainFields.OwningOrganizationGuid = "my-org-guid"
			cfAppDomainFields.Shared = false

			org := models.Organization{}
			org.Name = "my-org"
			org.Guid = "my-org-guid"
			org.QuotaDefinition = models.QuotaFields{
				Name:                    "cantina-quota",
				MemoryLimit:             512,
				InstanceMemoryLimit:     256,
				RoutesLimit:             2,
				ServicesLimit:           5,
				NonBasicServicesAllowed: true,
				AppInstanceLimit:        7,
			}
			org.Spaces = []models.SpaceFields{developmentSpaceFields, stagingSpaceFields}
			org.Domains = []models.DomainFields{domainFields, cfAppDomainFields}
			org.SpaceQuotas = []models.SpaceQuota{
				{Name: "space-quota-1", Guid: "space-quota-1-guid", MemoryLimit: 512, InstanceMemoryLimit: -1},
				{Name: "space-quota-2", Guid: "space-quota-2-guid", MemoryLimit: 256, InstanceMemoryLimit: 128},
			}

			requirementsFactory.LoginSuccess = true
			requirementsFactory.Organization = org
		})

		It("shows the org with the given name", func() {
			runCommand("my-org")

			Expect(requirementsFactory.OrganizationName).To(Equal("my-org"))
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Getting info for org", "my-org", "my-user"},
				[]string{"OK"},
				[]string{"my-org"},
				[]string{"domains:", "cfapps.io", "cf-app.com"},
				[]string{"quota: ", "cantina-quota", "512M", "256M instance memory limit", "2 routes", "5 services", "paid services allowed", "7 app instance limit"},
				[]string{"spaces:", "development", "staging"},
				[]string{"space quotas:", "space-quota-1", "space-quota-2"},
			))
		})

		Context("when the guid flag is provided", func() {
			It("shows only the org guid", func() {
				runCommand("--guid", "my-org")

				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"my-org-guid"},
				))

				Expect(ui.Outputs).ToNot(ContainSubstrings(
					[]string{"Getting info for org", "my-org", "my-user"},
				))
			})
		})

		Context("when invoked by a plugin", func() {
			var (
				pluginModel plugin_models.GetOrg_Model
			)
			BeforeEach(func() {
				pluginModel = plugin_models.GetOrg_Model{}
				deps.PluginModels.Organization = &pluginModel
			})

			It("populates the plugin model", func() {
				testcmd.RunCliCommand("org", []string{"my-org"}, requirementsFactory, updateCommandDependency, true)

				Expect(pluginModel.Name).To(Equal("my-org"))
				Expect(pluginModel.Guid).To(Equal("my-org-guid"))
				// quota
				Expect(pluginModel.QuotaDefinition.Name).To(Equal("cantina-quota"))
				Expect(pluginModel.QuotaDefinition.MemoryLimit).To(Equal(int64(512)))
				Expect(pluginModel.QuotaDefinition.InstanceMemoryLimit).To(Equal(int64(256)))
				Expect(pluginModel.QuotaDefinition.RoutesLimit).To(Equal(2))
				Expect(pluginModel.QuotaDefinition.ServicesLimit).To(Equal(5))
				Expect(pluginModel.QuotaDefinition.NonBasicServicesAllowed).To(BeTrue())

				// domains
				Expect(pluginModel.Domains).To(HaveLen(2))
				Expect(pluginModel.Domains[0].Name).To(Equal("cfapps.io"))
				Expect(pluginModel.Domains[0].Guid).To(Equal("1111"))
				Expect(pluginModel.Domains[0].OwningOrganizationGuid).To(Equal("my-org-guid"))
				Expect(pluginModel.Domains[0].Shared).To(BeTrue())
				Expect(pluginModel.Domains[1].Name).To(Equal("cf-app.com"))
				Expect(pluginModel.Domains[1].Guid).To(Equal("2222"))
				Expect(pluginModel.Domains[1].OwningOrganizationGuid).To(Equal("my-org-guid"))
				Expect(pluginModel.Domains[1].Shared).To(BeFalse())

				// spaces
				Expect(pluginModel.Spaces).To(HaveLen(2))
				Expect(pluginModel.Spaces[0].Name).To(Equal("development"))
				Expect(pluginModel.Spaces[0].Guid).To(Equal("dev-space-guid-1"))
				Expect(pluginModel.Spaces[1].Name).To(Equal("staging"))
				Expect(pluginModel.Spaces[1].Guid).To(Equal("staging-space-guid-1"))

				// space quotas
				Expect(pluginModel.SpaceQuotas).To(HaveLen(2))
				Expect(pluginModel.SpaceQuotas[0].Name).To(Equal("space-quota-1"))
				Expect(pluginModel.SpaceQuotas[0].Guid).To(Equal("space-quota-1-guid"))
				Expect(pluginModel.SpaceQuotas[0].MemoryLimit).To(Equal(int64(512)))
				Expect(pluginModel.SpaceQuotas[0].InstanceMemoryLimit).To(Equal(int64(-1)))
				Expect(pluginModel.SpaceQuotas[1].Name).To(Equal("space-quota-2"))
				Expect(pluginModel.SpaceQuotas[1].Guid).To(Equal("space-quota-2-guid"))
				Expect(pluginModel.SpaceQuotas[1].MemoryLimit).To(Equal(int64(256)))
				Expect(pluginModel.SpaceQuotas[1].InstanceMemoryLimit).To(Equal(int64(128)))
			})

		})
	})
})
