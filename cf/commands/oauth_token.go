package commands

import (
	"github.com/theophoric/cf-cli/cf/api/authentication"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
	"github.com/theophoric/cf-cli/plugin/models"
)

type OAuthToken struct {
	ui          terminal.UI
	config      core_config.ReadWriter
	authRepo    authentication.AuthenticationRepository
	pluginModel *plugin_models.GetOauthToken_Model
	pluginCall  bool
}

func init() {
	command_registry.Register(&OAuthToken{})
}

func (cmd *OAuthToken) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "oauth-token",
		Description: T("Retrieve and display the OAuth token for the current session"),
		Usage: []string{
			T("CF_NAME oauth-token"),
		},
	}
}

func (cmd *OAuthToken) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}

	return reqs
}

func (cmd *OAuthToken) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.authRepo = deps.RepoLocator.GetAuthenticationRepository()
	cmd.pluginCall = pluginCall
	cmd.pluginModel = deps.PluginModels.OauthToken
	return cmd
}

func (cmd *OAuthToken) Execute(c flags.FlagContext) {
	token, err := cmd.authRepo.RefreshAuthToken()
	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	if cmd.pluginCall {
		cmd.pluginModel.Token = token
	} else {
		cmd.ui.Say(token)
	}
}
