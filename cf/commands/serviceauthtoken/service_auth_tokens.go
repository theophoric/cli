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

type ListServiceAuthTokens struct {
	ui            terminal.UI
	config        core_config.Reader
	authTokenRepo api.ServiceAuthTokenRepository
}

func init() {
	command_registry.Register(&ListServiceAuthTokens{})
}

func (cmd *ListServiceAuthTokens) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "service-auth-tokens",
		Description: T("List service auth tokens"),
		Usage: []string{
			T("CF_NAME service-auth-tokens"),
		},
	}
}

func (cmd *ListServiceAuthTokens) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	usageReq := requirements.NewUsageRequirement(command_registry.CliCommandUsagePresenter(cmd),
		T("No argument required"),
		func() bool {
			return len(fc.Args()) != 0
		},
	)

	maximumVersion, err := semver.Make("2.46.0")
	if err != nil {
		panic(err.Error())
	}

	reqs := []requirements.Requirement{
		usageReq,
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewMaxAPIVersionRequirement("service-auth-tokens", maximumVersion),
	}

	return reqs
}

func (cmd *ListServiceAuthTokens) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.authTokenRepo = deps.RepoLocator.GetServiceAuthTokenRepository()
	return cmd
}

func (cmd *ListServiceAuthTokens) Execute(c flags.FlagContext) {
	cmd.ui.Say(T("Getting service auth tokens as {{.CurrentUser}}...",
		map[string]interface{}{
			"CurrentUser": terminal.EntityNameColor(cmd.config.Username()),
		}))
	authTokens, apiErr := cmd.authTokenRepo.FindAll()
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}
	cmd.ui.Ok()
	cmd.ui.Say("")

	table := terminal.NewTable(cmd.ui, []string{T("label"), T("provider")})

	for _, authToken := range authTokens {
		table.Add(authToken.Label, authToken.Provider)
	}

	table.Print()
}
