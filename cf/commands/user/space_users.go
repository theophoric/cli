package user

import (
	"github.com/theophoric/cf-cli/cf"
	"github.com/theophoric/cf-cli/cf/actors/userprint"
	"github.com/theophoric/cf-cli/cf/api"
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

type SpaceUsers struct {
	ui          terminal.UI
	config      core_config.Reader
	spaceRepo   spaces.SpaceRepository
	userRepo    api.UserRepository
	orgReq      requirements.OrganizationRequirement
	pluginModel *[]plugin_models.GetSpaceUsers_Model
	pluginCall  bool
}

func init() {
	command_registry.Register(&SpaceUsers{})
}

func (cmd *SpaceUsers) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "space-users",
		Description: T("Show space users by role"),
		Usage: []string{
			T("CF_NAME space-users ORG SPACE"),
		},
	}
}

func (cmd *SpaceUsers) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 2 {
		cmd.ui.Failed(T("Incorrect Usage. Requires arguments\n\n") + command_registry.Commands.CommandUsage("space-users"))
	}

	cmd.orgReq = requirementsFactory.NewOrganizationRequirement(fc.Args()[0])

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		cmd.orgReq,
	}

	return reqs
}

func (cmd *SpaceUsers) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.userRepo = deps.RepoLocator.GetUserRepository()
	cmd.spaceRepo = deps.RepoLocator.GetSpaceRepository()
	cmd.pluginCall = pluginCall
	cmd.pluginModel = deps.PluginModels.SpaceUsers

	return cmd
}

func (cmd *SpaceUsers) Execute(c flags.FlagContext) {
	spaceName := c.Args()[1]
	org := cmd.orgReq.GetOrganization()

	space, err := cmd.spaceRepo.FindByNameInOrg(spaceName, org.Guid)
	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	printer := cmd.printer(org, space, cmd.config.Username())
	printer.PrintUsers(space.Guid, cmd.config.Username())
}

func (cmd *SpaceUsers) printer(org models.Organization, space models.Space, username string) userprint.UserPrinter {
	var roles = []string{models.SPACE_MANAGER, models.SPACE_DEVELOPER, models.SPACE_AUDITOR}

	if cmd.pluginCall {
		return userprint.NewSpaceUsersPluginPrinter(
			cmd.pluginModel,
			cmd.userLister(),
			roles,
		)
	}

	cmd.ui.Say(T("Getting users in org {{.TargetOrg}} / space {{.TargetSpace}} as {{.CurrentUser}}",
		map[string]interface{}{
			"TargetOrg":   terminal.EntityNameColor(org.Name),
			"TargetSpace": terminal.EntityNameColor(space.Name),
			"CurrentUser": terminal.EntityNameColor(username),
		}))

	return &userprint.SpaceUsersUiPrinter{
		Ui:         cmd.ui,
		UserLister: cmd.userLister(),
		Roles:      roles,
		RoleDisplayNames: map[string]string{
			models.SPACE_MANAGER:   T("SPACE MANAGER"),
			models.SPACE_DEVELOPER: T("SPACE DEVELOPER"),
			models.SPACE_AUDITOR:   T("SPACE AUDITOR"),
		},
	}
}

func (cmd *SpaceUsers) userLister() func(spaceGuid string, role string) ([]models.UserFields, error) {
	if cmd.config.IsMinApiVersion(cf.ListUsersInOrgOrSpaceWithoutUAAMinimumApiVersion) {
		return cmd.userRepo.ListUsersInSpaceForRoleWithNoUAA
	}
	return cmd.userRepo.ListUsersInSpaceForRole
}
