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

type SpaceQuota struct {
	ui             terminal.UI
	config         core_config.Reader
	spaceQuotaRepo space_quotas.SpaceQuotaRepository
}

func init() {
	command_registry.Register(&SpaceQuota{})
}

func (cmd *SpaceQuota) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "space-quota",
		Description: T("Show space quota info"),
		Usage: []string{
			T("CF_NAME space-quota SPACE_QUOTA_NAME"),
		},
	}
}

func (cmd *SpaceQuota) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(T("Incorrect Usage. Requires an argument\n\n") + command_registry.Commands.CommandUsage("space-quota"))
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedOrgRequirement(),
	}

	return reqs
}

func (cmd *SpaceQuota) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.spaceQuotaRepo = deps.RepoLocator.GetSpaceQuotaRepository()
	return cmd
}

func (cmd *SpaceQuota) Execute(c flags.FlagContext) {
	name := c.Args()[0]

	cmd.ui.Say(T("Getting space quota {{.Quota}} info as {{.Username}}...",
		map[string]interface{}{
			"Quota":    terminal.EntityNameColor(name),
			"Username": terminal.EntityNameColor(cmd.config.Username()),
		}))

	spaceQuota, apiErr := cmd.spaceQuotaRepo.FindByName(name)

	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
	cmd.ui.Say("")
	var megabytes string

	table := terminal.NewTable(cmd.ui, []string{"", ""})
	table.Add(T("total memory limit"), formatters.ByteSize(spaceQuota.MemoryLimit*formatters.MEGABYTE))
	if spaceQuota.InstanceMemoryLimit == -1 {
		megabytes = T("unlimited")
	} else {
		megabytes = formatters.ByteSize(spaceQuota.InstanceMemoryLimit * formatters.MEGABYTE)
	}

	servicesLimit := strconv.Itoa(spaceQuota.ServicesLimit)
	if servicesLimit == "-1" {
		servicesLimit = T("unlimited")
	}

	table.Add(T("instance memory limit"), megabytes)
	table.Add(T("routes"), fmt.Sprintf("%d", spaceQuota.RoutesLimit))
	table.Add(T("services"), servicesLimit)
	table.Add(T("non basic services"), formatters.Allowed(spaceQuota.NonBasicServicesAllowed))
	table.Add(T("app instance limit"), T(spaceQuota.FormattedAppInstanceLimit()))

	table.Print()

}
