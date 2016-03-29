package quota_test

import (
	"github.com/theophoric/cf-cli/cf"
	"github.com/theophoric/cf-cli/cf/api/quotas/fakes"
	"github.com/theophoric/cf-cli/cf/api/resources"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/errors"
	testcmd "github.com/theophoric/cf-cli/testhelpers/commands"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
	. "github.com/theophoric/cf-cli/testhelpers/matchers"
	testreq "github.com/theophoric/cf-cli/testhelpers/requirements"
	testterm "github.com/theophoric/cf-cli/testhelpers/terminal"

	"github.com/theophoric/cf-cli/cf/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("app Command", func() {
	var (
		ui                  *testterm.FakeUI
		requirementsFactory *testreq.FakeReqFactory
		quotaRepo           *fakes.FakeQuotaRepository
		quota               models.QuotaFields
		configRepo          core_config.Repository
		deps                command_registry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.Config = configRepo
		deps.RepoLocator = deps.RepoLocator.SetQuotaRepository(quotaRepo)
		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("update-quota").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		configRepo = testconfig.NewRepositoryWithDefaults()
		requirementsFactory = &testreq.FakeReqFactory{
			LoginSuccess:         true,
			MinAPIVersionSuccess: true,
		}
		quotaRepo = &fakes.FakeQuotaRepository{}
	})

	runCommand := func(args ...string) bool {
		return testcmd.RunCliCommand("update-quota", args, requirementsFactory, updateCommandDependency, false)
	}

	Describe("requirements", func() {
		Context("login requirement", func() {
			It("fails if not logged in", func() {
				requirementsFactory.LoginSuccess = false

				Expect(runCommand("cf-plays-dwarf-fortress")).To(BeFalse())
			})
		})

		Context("usage requirement", func() {
			It("fails with usage when no arguments are given", func() {
				passed := runCommand()
				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"Incorrect Usage", "Requires an argument"},
				))
				Expect(passed).To(BeFalse())
			})
		})

		Context("the minimum API version requirement", func() {
			BeforeEach(func() {
				requirementsFactory.LoginSuccess = true
				requirementsFactory.MinAPIVersionSuccess = false
			})

			It("fails when the -a option is provided", func() {
				Expect(runCommand("quota-name", "-a", "10")).To(BeFalse())

				Expect(requirementsFactory.MinAPIVersionRequiredVersion).To(Equal(cf.OrgAppInstanceLimitMinimumApiVersion))
				Expect(requirementsFactory.MinAPIVersionFeatureName).To(Equal("Option '-a'"))
			})

			It("does not fail when the -a option is not provided", func() {
				Expect(runCommand("quota-name", "-m", "10G")).To(BeTrue())
			})
		})
	})

	Describe("updating quota fields", func() {
		BeforeEach(func() {
			quota = models.QuotaFields{
				Guid:             "quota-guid",
				Name:             "quota-name",
				MemoryLimit:      1024,
				RoutesLimit:      111,
				ServicesLimit:    222,
				AppInstanceLimit: 333,
			}
		})

		JustBeforeEach(func() {
			quotaRepo.FindByNameReturns(quota, nil)
		})

		Context("when the -i flag is provided", func() {
			It("updates the instance memory limit", func() {
				runCommand("-i", "15G", "quota-name")
				Expect(quotaRepo.UpdateArgsForCall(0).Name).To(Equal("quota-name"))
				Expect(quotaRepo.UpdateArgsForCall(0).InstanceMemoryLimit).To(Equal(int64(15360)))
			})

			It("totally accepts -1 as a value because it means unlimited", func() {
				runCommand("-i", "-1", "quota-name")
				Expect(quotaRepo.UpdateArgsForCall(0).Name).To(Equal("quota-name"))
				Expect(quotaRepo.UpdateArgsForCall(0).InstanceMemoryLimit).To(Equal(int64(-1)))
			})

			It("fails with usage when the value cannot be parsed", func() {
				runCommand("-m", "blasé", "le-tired")
				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"Incorrect Usage"},
				))
			})
		})

		Context("when the -a flag is provided", func() {
			It("updated the total number of application instances limit", func() {
				runCommand("-a", "2", "quota-name")
				Expect(quotaRepo.UpdateCallCount()).To(Equal(1))
				Expect(quotaRepo.UpdateArgsForCall(0).AppInstanceLimit).To(Equal(2))
			})

			It("totally accepts -1 as a value because it means unlimited", func() {
				runCommand("-a", "-1", "quota-name")
				Expect(quotaRepo.UpdateCallCount()).To(Equal(1))
				Expect(quotaRepo.UpdateArgsForCall(0).AppInstanceLimit).To(Equal(resources.UnlimitedAppInstances))
			})

			It("does not override the value if it's not provided", func() {
				runCommand("-s", "5", "quota-name")
				Expect(quotaRepo.UpdateCallCount()).To(Equal(1))
				Expect(quotaRepo.UpdateArgsForCall(0).AppInstanceLimit).To(Equal(333))
			})
		})

		Context("when the -m flag is provided", func() {
			It("updates the memory limit", func() {
				runCommand("-m", "15G", "quota-name")
				Expect(quotaRepo.UpdateArgsForCall(0).Name).To(Equal("quota-name"))
				Expect(quotaRepo.UpdateArgsForCall(0).MemoryLimit).To(Equal(int64(15360)))
			})

			It("fails with usage when the value cannot be parsed", func() {
				runCommand("-m", "blasé", "le-tired")
				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"Incorrect Usage"},
				))
			})
		})

		Context("when the -n flag is provided", func() {
			It("updates the quota name", func() {
				runCommand("-n", "quota-new-name", "quota-name")

				Expect(quotaRepo.UpdateArgsForCall(0).Name).To(Equal("quota-new-name"))

				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"Updating quota", "quota-name", "as", "my-user"},
					[]string{"OK"},
				))
			})
		})

		It("updates the total allowed services", func() {
			runCommand("-s", "9000", "quota-name")
			Expect(quotaRepo.UpdateArgsForCall(0).ServicesLimit).To(Equal(9000))
		})

		It("updates the total allowed routes", func() {
			runCommand("-r", "9001", "quota-name")
			Expect(quotaRepo.UpdateArgsForCall(0).RoutesLimit).To(Equal(9001))
		})

		Context("update paid service plans", func() {
			BeforeEach(func() {
				quota.NonBasicServicesAllowed = false
			})

			It("changes to paid service plan when --allow flag is provided", func() {
				runCommand("--allow-paid-service-plans", "quota-name")
				Expect(quotaRepo.UpdateArgsForCall(0).NonBasicServicesAllowed).To(BeTrue())
			})

			It("shows an error when both --allow and --disallow flags are provided", func() {
				runCommand("--allow-paid-service-plans", "--disallow-paid-service-plans", "quota-name")

				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"FAILED"},
					[]string{"Both flags are not permitted"},
				))
			})

			Context("when paid services are allowed", func() {
				BeforeEach(func() {
					quota.NonBasicServicesAllowed = true
				})
				It("changes to non-paid service plan when --disallow flag is provided", func() {
					quotaRepo.FindByNameReturns(quota, nil) // updating an existing quota

					runCommand("--disallow-paid-service-plans", "quota-name")
					Expect(quotaRepo.UpdateArgsForCall(0).NonBasicServicesAllowed).To(BeFalse())
				})
			})
		})

		It("shows an error when updating fails", func() {
			quotaRepo.UpdateReturns(errors.New("I accidentally a quota"))
			runCommand("-m", "1M", "dead-serious")
			Expect(ui.Outputs).To(ContainSubstrings([]string{"FAILED"}))
		})

		It("shows a message explaining the update", func() {
			quota.Name = "i-love-ui"
			quotaRepo.FindByNameReturns(quota, nil)

			runCommand("-m", "50G", "i-love-ui")
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Updating quota", "i-love-ui", "as", "my-user"},
				[]string{"OK"},
			))
		})

	})

	It("shows the user an error when finding the quota fails", func() {
		quotaRepo.FindByNameReturns(models.QuotaFields{}, errors.New("i can't believe it's not quotas!"))

		runCommand("-m", "50Somethings", "what-could-possibly-go-wrong?")
		Expect(ui.Outputs).To(ContainSubstrings([]string{"FAILED"}))
	})
})
