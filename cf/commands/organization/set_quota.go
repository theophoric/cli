package organization

import (
	"github.com/theophoric/cf-cli/cf/api/quotas"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type SetQuota struct {
	ui        terminal.UI
	config    core_config.Reader
	quotaRepo quotas.QuotaRepository
	orgReq    requirements.OrganizationRequirement
}

func init() {
	command_registry.Register(&SetQuota{})
}

func (cmd *SetQuota) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "set-quota",
		Description: T("Assign a quota to an org"),
		Usage: []string{
			T("CF_NAME set-quota ORG QUOTA\n\n"),
			T("TIP:\n"),
			T("   View allowable quotas with 'CF_NAME quotas'"),
		},
	}
}

func (cmd *SetQuota) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 2 {
		cmd.ui.Failed(T("Incorrect Usage. Requires ORG_NAME, QUOTA as arguments\n\n") + command_registry.Commands.CommandUsage("set-quota"))
	}

	cmd.orgReq = requirementsFactory.NewOrganizationRequirement(fc.Args()[0])

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		cmd.orgReq,
	}

	return reqs
}

func (cmd *SetQuota) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.quotaRepo = deps.RepoLocator.GetQuotaRepository()
	return cmd
}

func (cmd *SetQuota) Execute(c flags.FlagContext) {
	org := cmd.orgReq.GetOrganization()
	quotaName := c.Args()[1]
	quota, apiErr := cmd.quotaRepo.FindByName(quotaName)

	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Say(T("Setting quota {{.QuotaName}} to org {{.OrgName}} as {{.Username}}...",
		map[string]interface{}{
			"QuotaName": terminal.EntityNameColor(quota.Name),
			"OrgName":   terminal.EntityNameColor(org.Name),
			"Username":  terminal.EntityNameColor(cmd.config.Username())}))

	apiErr = cmd.quotaRepo.AssignQuotaToOrg(org.Guid, quota.Guid)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}
