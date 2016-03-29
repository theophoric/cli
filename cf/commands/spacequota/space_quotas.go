package spacequota

import (
	"fmt"
	"strconv"

	"github.com/theophoric/cf-cli/cf/api/space_quotas"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/formatters"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type ListSpaceQuotas struct {
	ui             terminal.UI
	config         core_config.Reader
	spaceQuotaRepo space_quotas.SpaceQuotaRepository
}

func init() {
	command_registry.Register(&ListSpaceQuotas{})
}

func (cmd *ListSpaceQuotas) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "space-quotas",
		Description: T("List available space resource quotas"),
		Usage: []string{
			T("CF_NAME space-quotas"),
		},
	}
}

func (cmd *ListSpaceQuotas) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
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

func (cmd *ListSpaceQuotas) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.spaceQuotaRepo = deps.RepoLocator.GetSpaceQuotaRepository()
	return cmd
}

func (cmd *ListSpaceQuotas) Execute(c flags.FlagContext) {
	cmd.ui.Say(T("Getting space quotas as {{.Username}}...", map[string]interface{}{"Username": terminal.EntityNameColor(cmd.config.Username())}))

	quotas, apiErr := cmd.spaceQuotaRepo.FindByOrg(cmd.config.OrganizationFields().Guid)

	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
	cmd.ui.Say("")

	table := terminal.NewTable(cmd.ui, []string{
		T("name"),
		T("total memory limit"),
		T("instance memory limit"),
		T("routes"),
		T("service instances"),
		T("paid service plans"),
		T("app instance limit"),
	})

	var megabytes string

	for _, quota := range quotas {
		if quota.InstanceMemoryLimit == -1 {
			megabytes = T("unlimited")
		} else {
			megabytes = formatters.ByteSize(quota.InstanceMemoryLimit * formatters.MEGABYTE)
		}

		servicesLimit := strconv.Itoa(quota.ServicesLimit)
		if servicesLimit == "-1" {
			servicesLimit = T("unlimited")
		}

		table.Add(
			quota.Name,
			formatters.ByteSize(quota.MemoryLimit*formatters.MEGABYTE),
			megabytes,
			fmt.Sprintf("%d", quota.RoutesLimit),
			fmt.Sprintf(servicesLimit),
			formatters.Allowed(quota.NonBasicServicesAllowed),
			T(quota.FormattedAppInstanceLimit()),
		)
	}

	table.Print()
}
