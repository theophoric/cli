package environmentvariablegroup

import (
	"github.com/theophoric/cf-cli/cf/api/environment_variable_groups"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	cf_errors "github.com/theophoric/cf-cli/cf/errors"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type SetStagingEnvironmentVariableGroup struct {
	ui                           terminal.UI
	config                       core_config.ReadWriter
	environmentVariableGroupRepo environment_variable_groups.EnvironmentVariableGroupsRepository
}

func init() {
	command_registry.Register(&SetStagingEnvironmentVariableGroup{})
}

func (cmd *SetStagingEnvironmentVariableGroup) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "set-staging-environment-variable-group",
		Description: T("Pass parameters as JSON to create a staging environment variable group"),
		ShortName:   "ssevg",
		Usage: []string{
			T(`CF_NAME set-staging-environment-variable-group '{"name":"value","name":"value"}'`),
		},
	}
}

func (cmd *SetStagingEnvironmentVariableGroup) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(T("Incorrect Usage. Requires an argument\n\n") + command_registry.Commands.CommandUsage("set-staging-environment-variable-group"))
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}
	return reqs
}

func (cmd *SetStagingEnvironmentVariableGroup) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.environmentVariableGroupRepo = deps.RepoLocator.GetEnvironmentVariableGroupsRepository()
	return cmd
}

func (cmd *SetStagingEnvironmentVariableGroup) Execute(c flags.FlagContext) {
	cmd.ui.Say(T("Setting the contents of the staging environment variable group as {{.Username}}...", map[string]interface{}{
		"Username": terminal.EntityNameColor(cmd.config.Username())}))

	err := cmd.environmentVariableGroupRepo.SetStaging(c.Args()[0])
	if err != nil {
		suggestionText := ""

		httpError, ok := err.(cf_errors.HttpError)
		if ok && httpError.ErrorCode() == cf_errors.MessageParseError {
			suggestionText = T(`

Your JSON string syntax is invalid.  Proper syntax is this:  cf set-staging-environment-variable-group '{"name":"value","name":"value"}'`)
		}
		cmd.ui.Failed(err.Error() + suggestionText)
	}

	cmd.ui.Ok()
}
