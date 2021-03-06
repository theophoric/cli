package application

import (
	"github.com/theophoric/cf-cli/cf/api/applications"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/flags"

	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
)

type GetHealthCheck struct {
	ui      terminal.UI
	config  core_config.Reader
	appReq  requirements.ApplicationRequirement
	appRepo applications.ApplicationRepository
}

func init() {
	command_registry.Register(&GetHealthCheck{})
}

func (cmd *GetHealthCheck) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "get-health-check",
		Description: T("Get the health_check_type value of an app"),
		Usage: []string{
			T("CF_NAME get-health-check APP_NAME"),
		},
	}
}

func (cmd *GetHealthCheck) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(T("Incorrect Usage. Requires APP_NAME as argument\n\n") + command_registry.Commands.CommandUsage("get-health-check"))
	}

	cmd.ui.Say(T("Getting health_check_type value for ") + terminal.EntityNameColor(fc.Args()[0]))
	cmd.ui.Say("")

	cmd.appReq = requirementsFactory.NewApplicationRequirement(fc.Args()[0])

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedSpaceRequirement(),
		cmd.appReq,
	}

	return reqs
}

func (cmd *GetHealthCheck) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.appRepo = deps.RepoLocator.GetApplicationRepository()
	return cmd
}

func (cmd *GetHealthCheck) Execute(fc flags.FlagContext) {
	app := cmd.appReq.GetApplication()

	cmd.ui.Say(T("health_check_type is ") + terminal.HeaderColor(app.HealthCheckType))
}
