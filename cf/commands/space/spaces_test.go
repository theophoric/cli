package space_test

import (
	testapi "github.com/theophoric/cf-cli/cf/api/fakes"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/trace/fakes"
	"github.com/theophoric/cf-cli/flags"
	"github.com/theophoric/cf-cli/plugin/models"
	testcmd "github.com/theophoric/cf-cli/testhelpers/commands"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
	testreq "github.com/theophoric/cf-cli/testhelpers/requirements"
	testterm "github.com/theophoric/cf-cli/testhelpers/terminal"

	"github.com/theophoric/cf-cli/cf/commands/space"
	. "github.com/theophoric/cf-cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("spaces command", func() {
	var (
		ui                  *testterm.FakeUI
		requirementsFactory *testreq.FakeReqFactory
		configRepo          core_config.Repository
		spaceRepo           *testapi.FakeSpaceRepository

		deps command_registry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.Config = configRepo
		deps.RepoLocator = deps.RepoLocator.SetSpaceRepository(spaceRepo)
		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("spaces").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		deps = command_registry.NewDependency(new(fakes.FakePrinter))
		ui = &testterm.FakeUI{}
		spaceRepo = &testapi.FakeSpaceRepository{}
		requirementsFactory = &testreq.FakeReqFactory{}
		configRepo = testconfig.NewRepositoryWithDefaults()
	})

	runCommand := func(args ...string) bool {
		return testcmd.RunCliCommand("spaces", args, requirementsFactory, updateCommandDependency, false)
	}

	Describe("requirements", func() {
		It("fails when not logged in", func() {
			requirementsFactory.TargetedOrgSuccess = true

			Expect(runCommand()).To(BeFalse())
		})

		It("fails when an org is not targeted", func() {
			requirementsFactory.LoginSuccess = true

			Expect(runCommand()).To(BeFalse())
		})

		Context("when arguments are provided", func() {
			var cmd command_registry.Command
			var flagContext flags.FlagContext

			BeforeEach(func() {
				cmd = &space.ListSpaces{}
				cmd.SetDependency(deps, false)
				flagContext = flags.NewFlagContext(cmd.MetaData().Flags)
			})

			It("should fail with usage", func() {
				flagContext.Parse("blahblah")

				reqs := cmd.Requirements(requirementsFactory, flagContext)

				err := testcmd.RunRequirements(reqs)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage"))
				Expect(err.Error()).To(ContainSubstring("No argument required"))
			})
		})
	})

	listSpacesStub := func(spaces []models.Space) func(func(models.Space) bool) error {
		return func(cb func(models.Space) bool) error {
			var keepGoing bool
			for _, s := range spaces {
				keepGoing = cb(s)
				if !keepGoing {
					return nil
				}
			}
			return nil
		}
	}

	Describe("when invoked by a plugin", func() {
		var (
			pluginModels []plugin_models.GetSpaces_Model
		)

		BeforeEach(func() {
			pluginModels = []plugin_models.GetSpaces_Model{}
			deps.PluginModels.Spaces = &pluginModels

			space := models.Space{}
			space.Name = "space1"
			space.Guid = "123"
			space2 := models.Space{}
			space2.Name = "space2"
			space2.Guid = "456"
			spaceRepo.ListSpacesStub = listSpacesStub([]models.Space{space, space2})

			requirementsFactory.TargetedOrgSuccess = true
			requirementsFactory.LoginSuccess = true

		})

		It("populates the plugin models upon execution", func() {
			testcmd.RunCliCommand("spaces", []string{}, requirementsFactory, updateCommandDependency, true)
			runCommand()
			Expect(pluginModels[0].Name).To(Equal("space1"))
			Expect(pluginModels[0].Guid).To(Equal("123"))
			Expect(pluginModels[1].Name).To(Equal("space2"))
			Expect(pluginModels[1].Guid).To(Equal("456"))
		})
	})

	Context("when logged in and an org is targeted", func() {
		BeforeEach(func() {
			space := models.Space{}
			space.Name = "space1"
			space2 := models.Space{}
			space2.Name = "space2"
			space3 := models.Space{}
			space3.Name = "space3"
			spaceRepo.ListSpacesStub = listSpacesStub([]models.Space{space, space2, space3})
			requirementsFactory.LoginSuccess = true
			requirementsFactory.TargetedOrgSuccess = true
		})

		It("lists all of the spaces", func() {
			runCommand()

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Getting spaces in org", "my-org", "my-user"},
				[]string{"space1"},
				[]string{"space2"},
				[]string{"space3"},
			))
		})

		Context("when there are no spaces", func() {
			BeforeEach(func() {
				spaceRepo.ListSpacesStub = listSpacesStub([]models.Space{})
			})

			It("politely tells the user", func() {
				runCommand()
				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"Getting spaces in org", "my-org", "my-user"},
					[]string{"No spaces found"},
				))
			})
		})
	})
})
