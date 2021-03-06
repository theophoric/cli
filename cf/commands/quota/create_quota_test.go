package quota_test

import (
	"github.com/theophoric/cf-cli/cf"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/theophoric/cf-cli/cf/api/quotas/fakes"
	"github.com/theophoric/cf-cli/cf/api/resources"
	"github.com/theophoric/cf-cli/cf/errors"
	testcmd "github.com/theophoric/cf-cli/testhelpers/commands"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
	testreq "github.com/theophoric/cf-cli/testhelpers/requirements"
	testterm "github.com/theophoric/cf-cli/testhelpers/terminal"
)

var _ = Describe("create-quota command", func() {
	var (
		ui                  *testterm.FakeUI
		quotaRepo           *fakes.FakeQuotaRepository
		requirementsFactory *testreq.FakeReqFactory
		configRepo          core_config.Repository
		deps                command_registry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.Config = configRepo
		deps.RepoLocator = deps.RepoLocator.SetQuotaRepository(quotaRepo)
		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("create-quota").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		configRepo = testconfig.NewRepositoryWithDefaults()
		quotaRepo = &fakes.FakeQuotaRepository{}
		requirementsFactory = &testreq.FakeReqFactory{}
	})

	runCommand := func(args ...string) bool {
		return testcmd.RunCliCommand("create-quota", args, requirementsFactory, updateCommandDependency, false)
	}

	Context("when the user is not logged in", func() {
		BeforeEach(func() {
			requirementsFactory.LoginSuccess = false
		})

		It("fails requirements", func() {
			Expect(runCommand("my-quota", "-m", "50G")).To(BeFalse())
		})
	})

	Context("the minimum API version requirement", func() {
		BeforeEach(func() {
			requirementsFactory.LoginSuccess = true
			requirementsFactory.MinAPIVersionSuccess = false
		})

		It("fails when the -a option is provided", func() {
			Expect(runCommand("my-quota", "-a", "10")).To(BeFalse())

			Expect(requirementsFactory.MinAPIVersionRequiredVersion).To(Equal(cf.OrgAppInstanceLimitMinimumApiVersion))
			Expect(requirementsFactory.MinAPIVersionFeatureName).To(Equal("Option '-a'"))
		})

		It("does not fail when the -a option is not provided", func() {
			Expect(runCommand("my-quota", "-m", "10G")).To(BeTrue())
		})
	})

	Context("when the user is logged in", func() {
		BeforeEach(func() {
			requirementsFactory.LoginSuccess = true
			requirementsFactory.MinAPIVersionSuccess = true
		})

		It("fails requirements when called without a quota name", func() {
			runCommand()
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Incorrect Usage", "Requires an argument"},
			))
		})

		It("creates a quota with a given name", func() {
			runCommand("my-quota")
			Expect(quotaRepo.CreateArgsForCall(0).Name).To(Equal("my-quota"))
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Creating quota", "my-quota", "my-user", "..."},
				[]string{"OK"},
			))
		})

		Context("when the -i flag is not provided", func() {
			It("defaults the memory limit to unlimited", func() {
				runCommand("my-quota")

				Expect(quotaRepo.CreateArgsForCall(0).InstanceMemoryLimit).To(Equal(int64(-1)))
			})
		})

		Context("when the -m flag is provided", func() {
			It("sets the memory limit", func() {
				runCommand("-m", "50G", "erryday makin fitty jeez")
				Expect(quotaRepo.CreateArgsForCall(0).MemoryLimit).To(Equal(int64(51200)))
			})

			It("alerts the user when parsing the memory limit fails", func() {
				runCommand("whoops", "12")

				Expect(ui.Outputs).To(ContainSubstrings([]string{"FAILED"}))
			})
		})

		Context("when the -i flag is provided", func() {
			It("sets the memory limit", func() {
				runCommand("-i", "50G", "erryday makin fitty jeez")
				Expect(quotaRepo.CreateArgsForCall(0).InstanceMemoryLimit).To(Equal(int64(51200)))
			})

			It("alerts the user when parsing the memory limit fails", func() {
				runCommand("-i", "whoops", "wit mah hussle", "12")

				Expect(ui.Outputs).To(ContainSubstrings([]string{"FAILED"}))
			})

			Context("and the provided value is -1", func() {
				It("sets the memory limit", func() {
					runCommand("-i", "-1", "yo")
					Expect(quotaRepo.CreateArgsForCall(0).InstanceMemoryLimit).To(Equal(int64(-1)))
				})
			})
		})

		Context("when the -a flag is provided", func() {
			It("sets the app limit", func() {
				runCommand("my-quota", "-a", "10")

				Expect(quotaRepo.CreateArgsForCall(0).AppInstanceLimit).To(Equal(10))
			})

			It("defaults to unlimited", func() {
				runCommand("my-quota")

				Expect(quotaRepo.CreateArgsForCall(0).AppInstanceLimit).To(Equal(resources.UnlimitedAppInstances))
			})
		})

		It("sets the route limit", func() {
			runCommand("-r", "12", "ecstatic")

			Expect(quotaRepo.CreateArgsForCall(0).RoutesLimit).To(Equal(12))
		})

		It("sets the service instance limit", func() {
			runCommand("-s", "42", "black star")
			Expect(quotaRepo.CreateArgsForCall(0).ServicesLimit).To(Equal(42))
		})

		Context("when requesting to allow paid service plans", func() {
			It("creates the quota with paid service plans allowed", func() {
				runCommand("--allow-paid-service-plans", "my-for-profit-quota")
				Expect(quotaRepo.CreateArgsForCall(0).NonBasicServicesAllowed).To(BeTrue())
			})

			It("defaults to not allowing paid service plans", func() {
				runCommand("my-pro-bono-quota")
				Expect(quotaRepo.CreateArgsForCall(0).NonBasicServicesAllowed).To(BeFalse())
			})
		})

		Context("when creating a quota returns an error", func() {
			It("alerts the user when creating the quota fails", func() {
				quotaRepo.CreateReturns(errors.New("WHOOP THERE IT IS"))
				runCommand("my-quota")

				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"Creating quota", "my-quota"},
					[]string{"FAILED"},
				))
			})

			It("warns the user when quota already exists", func() {
				quotaRepo.CreateReturns(errors.NewHttpError(400, errors.QuotaDefinitionNameTaken, "Quota Definition is taken: quota-sct"))
				runCommand("Banana")

				Expect(ui.Outputs).ToNot(ContainSubstrings(
					[]string{"FAILED"},
				))
				Expect(ui.WarnOutputs).To(ContainSubstrings([]string{"already exists"}))
			})

		})
	})
})
