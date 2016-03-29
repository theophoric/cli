package commands

import (
	"github.com/theophoric/cf-cli/cf"
	"github.com/theophoric/cf-cli/cf/api/authentication"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type Authenticate struct {
	ui            terminal.UI
	config        core_config.ReadWriter
	authenticator authentication.AuthenticationRepository
}

func init() {
	command_registry.Register(&Authenticate{})
}

func (cmd *Authenticate) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "auth",
		Description: T("Authenticate user non-interactively"),
		Usage: []string{
			T("CF_NAME auth USERNAME PASSWORD\n\n"),
			terminal.WarningColor(T("WARNING:\n   Providing your password as a command line option is highly discouraged\n   Your password may be visible to others and may be recorded in your shell history")),
		},
		Examples: []string{
			T("CF_NAME auth name@example.com \"my password\" (use quotes for passwords with a space)"),
			T("CF_NAME auth name@example.com \"\\\"password\\\"\" (escape quotes if used in password)"),
		},
	}
}

func (cmd *Authenticate) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 2 {
		cmd.ui.Failed(T("Incorrect Usage. Requires 'username password' as arguments\n\n") + command_registry.Commands.CommandUsage("auth"))
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewApiEndpointRequirement(),
	}

	return reqs
}

func (cmd *Authenticate) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.authenticator = deps.RepoLocator.GetAuthenticationRepository()
	return cmd
}

func (cmd *Authenticate) Execute(c flags.FlagContext) {
	cmd.config.ClearSession()
	cmd.authenticator.GetLoginPromptsAndSaveUAAServerURL()

	cmd.ui.Say(T("API endpoint: {{.ApiEndpoint}}",
		map[string]interface{}{"ApiEndpoint": terminal.EntityNameColor(cmd.config.ApiEndpoint())}))
	cmd.ui.Say(T("Authenticating..."))

	apiErr := cmd.authenticator.Authenticate(map[string]string{"username": c.Args()[0], "password": c.Args()[1]})
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
	cmd.ui.Say(T("Use '{{.Name}}' to view or set your target org and space",
		map[string]interface{}{"Name": terminal.CommandColor(cf.Name + " target")}))

	cmd.ui.NotifyUpdateIfNeeded(cmd.config)

	return
}
