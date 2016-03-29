package commands_test

import (
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/models"
	testcmd "github.com/theophoric/cf-cli/testhelpers/commands"
	testconfig "github.com/theophoric/cf-cli/testhelpers/configuration"
	testterm "github.com/theophoric/cf-cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("logout command", func() {

	var (
		config core_config.Repository
		ui     *testterm.FakeUI
		deps   command_registry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.Config = config
		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("logout").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		org := models.OrganizationFields{}
		org.Name = "MyOrg"

		space := models.SpaceFields{}
		space.Name = "MySpace"

		config = testconfig.NewRepository()
		config.SetAccessToken("MyAccessToken")
		config.SetOrganizationFields(org)
		config.SetSpaceFields(space)
		ui = &testterm.FakeUI{}

		testcmd.RunCliCommand("logout", []string{}, nil, updateCommandDependency, false)
	})

	It("clears access token from the config", func() {
		Expect(config.AccessToken()).To(Equal(""))
	})

	It("clears organization fields from the config", func() {
		Expect(config.OrganizationFields()).To(Equal(models.OrganizationFields{}))
	})

	It("clears space fields from the config", func() {
		Expect(config.SpaceFields()).To(Equal(models.SpaceFields{}))
	})
})
