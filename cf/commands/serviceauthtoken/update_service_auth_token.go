package serviceauthtoken

import (
	"github.com/blang/semver"
	"github.com/theophoric/cf-cli/cf/api"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type UpdateServiceAuthTokenFields struct {
	ui            terminal.UI
	config        core_config.Reader
	authTokenRepo api.ServiceAuthTokenRepository
}

func init() {
	command_registry.Register(&UpdateServiceAuthTokenFields{})
}

func (cmd *UpdateServiceAuthTokenFields) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "update-service-auth-token",
		Description: T("Update a service auth token"),
		Usage: []string{
			T("CF_NAME update-service-auth-token LABEL PROVIDER TOKEN"),
		},
	}
}

func (cmd *UpdateServiceAuthTokenFields) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 3 {
		cmd.ui.Failed(T("Incorrect Usage. Requires LABEL, PROVIDER and TOKEN as arguments\n\n") + command_registry.Commands.CommandUsage("update-service-auth-token"))
	}

	maximumVersion, err := semver.Make("2.46.0")
	if err != nil {
		panic(err.Error())
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewMaxAPIVersionRequirement("update-service-auth-token", maximumVersion),
	}

	return reqs
}

func (cmd *UpdateServiceAuthTokenFields) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.authTokenRepo = deps.RepoLocator.GetServiceAuthTokenRepository()
	return cmd
}

func (cmd *UpdateServiceAuthTokenFields) Execute(c flags.FlagContext) {
	cmd.ui.Say(T("Updating service auth token as {{.CurrentUser}}...", map[string]interface{}{"CurrentUser": terminal.EntityNameColor(cmd.config.Username())}))

	serviceAuthToken, apiErr := cmd.authTokenRepo.FindByLabelAndProvider(c.Args()[0], c.Args()[1])
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	serviceAuthToken.Token = c.Args()[2]

	apiErr = cmd.authTokenRepo.Update(serviceAuthToken)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}
