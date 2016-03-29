package spacequota_test

import (
	"github.com/theophoric/cf-cli/cf/api/space_quotas/fakes"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/errors"
	"github.com/theophoric/cf-cli/cf/models"
	testcmd "github.com/theophoric/cf-cli/testhelpers/commands"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
	testreq "github.com/theophoric/cf-cli/testhelpers/requirements"
	testterm "github.com/theophoric/cf-cli/testhelpers/terminal"

	. "github.com/theophoric/cf-cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("quotas command", func() {
	var (
		ui                  *testterm.FakeUI
		quotaRepo           *fakes.FakeSpaceQuotaRepository
		requirementsFactory *testreq.FakeReqFactory
		configRepo          core_config.Repository
		deps                command_registry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.RepoLocator = deps.RepoLocator.SetSpaceQuotaRepository(quotaRepo)
		deps.Config = configRepo
		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("space-quota").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		quotaRepo = &fakes.FakeSpaceQuotaRepository{}
		configRepo = testconfig.NewRepositoryWithDefaults()
		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true}
	})

	runCommand := func(args ...string) bool {
		return testcmd.RunCliCommand("space-quota", args, requirementsFactory, updateCommandDependency, false)
	}

	Describe("requirements", func() {
		It("requires the user to be logged in", func() {
			requirementsFactory.LoginSuccess = false
			Expect(runCommand("foo")).ToNot(HavePassedRequirements())
		})

		It("requires the user to target an org", func() {
			requirementsFactory.TargetedOrgSuccess = false
			Expect(runCommand("bar")).ToNot(HavePassedRequirements())
		})

		It("fails when a quota name is not provided", func() {
			requirementsFactory.LoginSuccess = true
			requirementsFactory.TargetedOrgSuccess = true
			Expect(runCommand()).ToNot(HavePassedRequirements())
		})
	})

	Context("when logged in", func() {
		JustBeforeEach(func() {
			requirementsFactory.LoginSuccess = true
			requirementsFactory.TargetedOrgSuccess = true
			Expect(runCommand("quota-name")).To(HavePassedRequirements())
		})

		Context("when quotas exist", func() {
			BeforeEach(func() {
				quotaRepo.FindByNameReturns(
					models.SpaceQuota{
						Name:                    "quota-name",
						MemoryLimit:             1024,
						InstanceMemoryLimit:     -1,
						RoutesLimit:             111,
						ServicesLimit:           222,
						NonBasicServicesAllowed: true,
						OrgGuid:                 "my-org-guid",
						AppInstanceLimit:        5,
					}, nil)
			})

			It("lists the specific quota info", func() {
				Expect(quotaRepo.FindByNameArgsForCall(0)).To(Equal("quota-name"))
				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"Getting space quota quota-name info as", "my-user"},
					[]string{"OK"},
					[]string{"total memory limit", "1G"},
					[]string{"instance memory limit", "unlimited"},
					[]string{"routes", "111"},
					[]string{"service", "222"},
					[]string{"non basic services", "allowed"},
					[]string{"app instance limit", "5"},
				))
			})

			Context("when the services are unlimited", func() {
				BeforeEach(func() {
					quotaRepo.FindByNameReturns(
						models.SpaceQuota{
							Name:                    "quota-name",
							MemoryLimit:             1024,
							InstanceMemoryLimit:     14,
							RoutesLimit:             111,
							ServicesLimit:           -1,
							NonBasicServicesAllowed: true,
							OrgGuid:                 "my-org-guid",
						}, nil)

				})

				It("replaces -1 with unlimited", func() {
					Expect(quotaRepo.FindByNameArgsForCall(0)).To(Equal("quota-name"))
					Expect(ui.Outputs).To(ContainSubstrings(
						[]string{"Getting space quota quota-name info as", "my-user"},
						[]string{"OK"},
						[]string{"total memory limit", "1G"},
						[]string{"instance memory limit", "14M"},
						[]string{"routes", "111"},
						[]string{"service", "unlimited"},
						[]string{"non basic services", "allowed"},
					))
				})
			})

			Context("when the app instances are unlimitied", func() {
				BeforeEach(func() {
					quotaRepo.FindByNameReturns(
						models.SpaceQuota{
							Name:                    "quota-name",
							MemoryLimit:             1024,
							InstanceMemoryLimit:     -1,
							RoutesLimit:             111,
							ServicesLimit:           222,
							NonBasicServicesAllowed: true,
							OrgGuid:                 "my-org-guid",
							AppInstanceLimit:        -1,
						}, nil)
				})

				It("replaces -1 with unlimited", func() {
					Expect(quotaRepo.FindByNameArgsForCall(0)).To(Equal("quota-name"))
					Expect(ui.Outputs).To(ContainSubstrings(
						[]string{"Getting space quota quota-name info as", "my-user"},
						[]string{"OK"},
						[]string{"total memory limit", "1G"},
						[]string{"instance memory limit", "unlimited"},
						[]string{"routes", "111"},
						[]string{"service", "222"},
						[]string{"non basic services", "allowed"},
						[]string{"app instance limit", "unlimited"},
					))
				})
			})

			Context("when the app instances are not defined", func() {
				BeforeEach(func() {
					quotaRepo.FindByNameReturns(
						models.SpaceQuota{
							Name:                    "quota-name",
							MemoryLimit:             1024,
							InstanceMemoryLimit:     -1,
							RoutesLimit:             111,
							ServicesLimit:           222,
							NonBasicServicesAllowed: true,
							OrgGuid:                 "my-org-guid",
							AppInstanceLimit:        -1,
						}, nil)
				})

				It("replaces -1 with unlimited", func() {
					Expect(quotaRepo.FindByNameArgsForCall(0)).To(Equal("quota-name"))
					Expect(ui.Outputs).To(ContainSubstrings(
						[]string{"Getting space quota quota-name info as", "my-user"},
						[]string{"OK"},
						[]string{"total memory limit", "1G"},
						[]string{"instance memory limit", "unlimited"},
						[]string{"routes", "111"},
						[]string{"service", "222"},
						[]string{"non basic services", "allowed"},
						[]string{"app instance limit", "unlimited"},
					))
				})
			})
		})
		Context("when an error occurs fetching quotas", func() {
			BeforeEach(func() {
				quotaRepo.FindByNameReturns(models.SpaceQuota{}, errors.New("I haz a borken!"))
			})

			It("prints an error", func() {
				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"Getting space quota quota-name info as", "my-user"},
					[]string{"FAILED"},
				))
			})
		})
	})

})
