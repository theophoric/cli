package spacequota

import (
	"github.com/theophoric/cf-cli/cf/api/space_quotas"
	"github.com/theophoric/cf-cli/cf/api/spaces"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type UnsetSpaceQuota struct {
	ui        terminal.UI
	config    core_config.Reader
	quotaRepo space_quotas.SpaceQuotaRepository
	spaceRepo spaces.SpaceRepository
}

func init() {
	command_registry.Register(&UnsetSpaceQuota{})
}

func (cmd *UnsetSpaceQuota) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "unset-space-quota",
		Description: T("Unassign a quota from a space"),
		Usage: []string{
			T("CF_NAME unset-space-quota SPACE QUOTA\n\n"),
		},
	}
}

func (cmd *UnsetSpaceQuota) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 2 {
		cmd.ui.Failed(T("Incorrect Usage. Requires SPACE and QUOTA as arguments\n\n") + command_registry.Commands.CommandUsage("unset-space-quota"))
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedOrgRequirement(),
	}

	return reqs
}

func (cmd *UnsetSpaceQuota) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.spaceRepo = deps.RepoLocator.GetSpaceRepository()
	cmd.quotaRepo = deps.RepoLocator.GetSpaceQuotaRepository()
	return cmd
}

func (cmd *UnsetSpaceQuota) Execute(c flags.FlagContext) {
	spaceName := c.Args()[0]
	quotaName := c.Args()[1]

	space, apiErr := cmd.spaceRepo.FindByName(spaceName)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	quota, apiErr := cmd.quotaRepo.FindByName(quotaName)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Say(T("Unassigning space quota {{.QuotaName}} from space {{.SpaceName}} as {{.Username}}...",
		map[string]interface{}{
			"QuotaName": terminal.EntityNameColor(quota.Name),
			"SpaceName": terminal.EntityNameColor(space.Name),
			"Username":  terminal.EntityNameColor(cmd.config.Username())}))

	apiErr = cmd.quotaRepo.UnassignQuotaFromSpace(space.Guid, quota.Guid)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}
