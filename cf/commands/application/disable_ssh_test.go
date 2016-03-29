package application_test

import (
	"errors"

	testApplication "github.com/theophoric/cf-cli/cf/api/applications/fakes"
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

var _ = Describe("disable-ssh command", func() {
	var (
		ui                  *testterm.FakeUI
		requirementsFactory *testreq.FakeReqFactory
		appRepo             *testApplication.FakeApplicationRepository
		configRepo          core_config.Repository
		deps                command_registry.Dependency
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		configRepo = testconfig.NewRepositoryWithDefaults()
		requirementsFactory = &testreq.FakeReqFactory{}
		appRepo = &testApplication.FakeApplicationRepository{}
	})

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.Config = configRepo
		deps.RepoLocator = deps.RepoLocator.SetApplicationRepository(appRepo)
		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("disable-ssh").SetDependency(deps, pluginCall))
	}

	runCommand := func(args ...string) bool {
		return testcmd.RunCliCommand("disable-ssh", args, requirementsFactory, updateCommandDependency, false)
	}

	Describe("requirements", func() {
		It("fails with usage when called without enough arguments", func() {
			requirementsFactory.LoginSuccess = true

			runCommand()
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Incorrect Usage", "Requires", "argument"},
			))

		})

		It("fails requirements when not logged in", func() {
			Expect(runCommand("my-app", "none")).To(BeFalse())
		})

		It("fails if a space is not targeted", func() {
			requirementsFactory.LoginSuccess = true
			requirementsFactory.TargetedSpaceSuccess = false
			Expect(runCommand("my-app", "none")).To(BeFalse())
		})
	})

	Describe("disable-ssh", func() {
		var (
			app models.Application
		)

		BeforeEach(func() {
			requirementsFactory.LoginSuccess = true
			requirementsFactory.TargetedSpaceSuccess = true

			app = models.Application{}
			app.Name = "my-app"
			app.Guid = "my-app-guid"
			app.EnableSsh = true

			requirementsFactory.Application = app
		})

		Context("when enable_ssh is already set to the false", func() {
			BeforeEach(func() {
				app.EnableSsh = false
				requirementsFactory.Application = app
			})

			It("notifies the user", func() {
				runCommand("my-app")

				Expect(ui.Outputs).To(ContainSubstrings([]string{"ssh support is already disabled for 'my-app'"}))
			})
		})

		Context("Updating enable_ssh when not already set to false", func() {
			Context("Update successfully", func() {
				BeforeEach(func() {
					app = models.Application{}
					app.Name = "my-app"
					app.Guid = "my-app-guid"
					app.EnableSsh = false

					appRepo.UpdateReturns(app, nil)
				})

				It("updates the app's enable_ssh", func() {
					runCommand("my-app")

					Expect(appRepo.UpdateCallCount()).To(Equal(1))
					appGUID, params := appRepo.UpdateArgsForCall(0)
					Expect(appGUID).To(Equal("my-app-guid"))
					Expect(*params.EnableSsh).To(BeFalse())
					Expect(ui.Outputs).To(ContainSubstrings([]string{"Disabling ssh support for 'my-app'"}))
					Expect(ui.Outputs).To(ContainSubstrings([]string{"OK"}))
				})
			})

			Context("Update fails", func() {
				It("notifies user of any api error", func() {
					appRepo.UpdateReturns(models.Application{}, errors.New("Error updating app."))
					runCommand("my-app")

					Expect(appRepo.UpdateCallCount()).To(Equal(1))
					Expect(ui.Outputs).To(ContainSubstrings(
						[]string{"FAILED"},
						[]string{"Error disabling ssh support"},
					))

				})

				It("notifies user when updated result is not in the desired state", func() {
					app = models.Application{}
					app.Name = "my-app"
					app.Guid = "my-app-guid"
					app.EnableSsh = true
					appRepo.UpdateReturns(app, nil)

					runCommand("my-app")

					Expect(appRepo.UpdateCallCount()).To(Equal(1))
					Expect(ui.Outputs).To(ContainSubstrings(
						[]string{"FAILED"},
						[]string{"ssh support is not disabled for my-app"},
					))

				})
			})
		})
	})
})
