package routergroups

import (
	"github.com/theophoric/cf-cli/cf/api"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type RouterGroups struct {
	ui             terminal.UI
	routingApiRepo api.RoutingApiRepository
	config         core_config.Reader
}

func init() {
	command_registry.Register(&RouterGroups{})
}

func (cmd *RouterGroups) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "router-groups",
		Description: T("List router groups"),
		Usage: []string{
			"CF_NAME router-groups",
		},
	}
}

func (cmd *RouterGroups) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	usageReq := requirements.NewUsageRequirement(command_registry.CliCommandUsagePresenter(cmd),
		T("No argument required"),
		func() bool {
			return len(fc.Args()) != 0
		},
	)

	reqs := []requirements.Requirement{
		usageReq,
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewRoutingAPIRequirement(),
	}
	return reqs
}

func (cmd *RouterGroups) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.routingApiRepo = deps.RepoLocator.GetRoutingApiRepository()
	return cmd
}

func (cmd *RouterGroups) Execute(c flags.FlagContext) {
	cmd.ui.Say(T("Getting router groups as {{.Username}} ...\n",
		map[string]interface{}{"Username": terminal.EntityNameColor(cmd.config.Username())}))

	table := cmd.ui.Table([]string{T("name"), T("type")})

	noRouterGroups := true
	cb := func(group models.RouterGroup) bool {
		noRouterGroups = false
		table.Add(group.Name, group.Type)
		return true
	}

	apiErr := cmd.routingApiRepo.ListRouterGroups(cb)
	if apiErr != nil {
		cmd.ui.Failed(T("Failed fetching router groups.\n{{.Err}}", map[string]interface{}{"Err": apiErr.Error()}))
		return
	}

	if noRouterGroups {
		cmd.ui.Say(T("No router groups found"))
	}

	table.Print()
}
