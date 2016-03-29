package space

import (
	"github.com/theophoric/cf-cli/cf/api/spaces"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
	"github.com/theophoric/cf-cli/plugin/models"
)

type ListSpaces struct {
	ui        terminal.UI
	config    core_config.Reader
	spaceRepo spaces.SpaceRepository

	pluginModel *[]plugin_models.GetSpaces_Model
	pluginCall  bool
}

func init() {
	command_registry.Register(&ListSpaces{})
}

func (cmd *ListSpaces) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "spaces",
		Description: T("List all spaces in an org"),
		Usage: []string{
			T("CF_NAME spaces"),
		},
	}

}

func (cmd *ListSpaces) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	usageReq := requirements.NewUsageRequirement(command_registry.CliCommandUsagePresenter(cmd),
		T("No argument required"),
		func() bool {
			return len(fc.Args()) != 0
		},
	)

	reqs := []requirements.Requirement{
		usageReq,
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedOrgRequirement(),
	}

	return reqs
}

func (cmd *ListSpaces) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.spaceRepo = deps.RepoLocator.GetSpaceRepository()
	cmd.pluginCall = pluginCall
	cmd.pluginModel = deps.PluginModels.Spaces
	return cmd
}

func (cmd *ListSpaces) Execute(c flags.FlagContext) {
	cmd.ui.Say(T("Getting spaces in org {{.TargetOrgName}} as {{.CurrentUser}}...\n",
		map[string]interface{}{
			"TargetOrgName": terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
			"CurrentUser":   terminal.EntityNameColor(cmd.config.Username()),
		}))

	foundSpaces := false
	table := cmd.ui.Table([]string{T("name")})
	apiErr := cmd.spaceRepo.ListSpaces(func(space models.Space) bool {
		table.Add(space.Name)
		foundSpaces = true

		if cmd.pluginCall {
			s := plugin_models.GetSpaces_Model{}
			s.Name = space.Name
			s.Guid = space.Guid
			*(cmd.pluginModel) = append(*(cmd.pluginModel), s)
		}

		return true
	})
	table.Print()

	if apiErr != nil {
		cmd.ui.Failed(T("Failed fetching spaces.\n{{.ErrorDescription}}",
			map[string]interface{}{
				"ErrorDescription": apiErr.Error(),
			}))
		return
	}

	if !foundSpaces {
		cmd.ui.Say(T("No spaces found"))
	}
}
