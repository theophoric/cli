package route

import (
	"fmt"
	"strings"

	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/flags"

	"github.com/theophoric/cf-cli/cf/api"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
)

type ListRoutes struct {
	ui        terminal.UI
	routeRepo api.RouteRepository
	config    core_config.Reader
}

func init() {
	command_registry.Register(&ListRoutes{})
}

func (cmd *ListRoutes) MetaData() command_registry.CommandMetadata {
	fs := make(map[string]flags.FlagSet)
	fs["orglevel"] = &flags.BoolFlag{Name: "orglevel", Usage: T("List all the routes for all spaces of current organization")}

	return command_registry.CommandMetadata{
		Name:        "routes",
		ShortName:   "r",
		Description: T("List all routes in the current space or the current organization"),
		Usage: []string{
			"CF_NAME routes [--orglevel]",
		},
		Flags: fs,
	}
}

func (cmd *ListRoutes) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	usageReq := requirements.NewUsageRequirement(command_registry.CliCommandUsagePresenter(cmd),
		T("No argument required"),
		func() bool {
			return len(fc.Args()) != 0
		},
	)

	reqs := []requirements.Requirement{
		usageReq,
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedSpaceRequirement(),
	}

	return reqs
}

func (cmd *ListRoutes) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.routeRepo = deps.RepoLocator.GetRouteRepository()
	return cmd
}

func (cmd *ListRoutes) Execute(c flags.FlagContext) {
	orglevel := c.Bool("orglevel")

	if orglevel {
		cmd.ui.Say(T("Getting routes for org {{.OrgName}} as {{.Username}} ...\n",
			map[string]interface{}{
				"Username": terminal.EntityNameColor(cmd.config.Username()),
				"OrgName":  terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
			}))
	} else {
		cmd.ui.Say(T("Getting routes for org {{.OrgName}} / space {{.SpaceName}} as {{.Username}} ...\n",
			map[string]interface{}{
				"Username":  terminal.EntityNameColor(cmd.config.Username()),
				"OrgName":   terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
				"SpaceName": terminal.EntityNameColor(cmd.config.SpaceFields().Name),
			}))
	}

	table := cmd.ui.Table([]string{T("space"), T("host"), T("domain"), T("port"), T("path"), T("type"), T("apps"), T("service")})

	var routesFound bool
	cb := func(route models.Route) bool {
		routesFound = true
		appNames := []string{}
		for _, app := range route.Apps {
			appNames = append(appNames, app.Name)
		}

		var port string
		if route.Port != 0 {
			port = fmt.Sprintf("%d", route.Port)
		}

		table.Add(
			route.Space.Name,
			route.Host,
			route.Domain.Name,
			port,
			route.Path,
			strings.Join(route.Domain.RouterGroupTypes, ","),
			strings.Join(appNames, ","),
			route.ServiceInstance.Name,
		)
		return true
	}

	var err error
	if orglevel {
		err = cmd.routeRepo.ListAllRoutes(cb)
	} else {
		err = cmd.routeRepo.ListRoutes(cb)
	}

	table.Print()
	if err != nil {
		cmd.ui.Failed(T("Failed fetching routes.\n{{.Err}}", map[string]interface{}{"Err": err.Error()}))
	}

	if !routesFound {
		cmd.ui.Say(T("No routes found"))
	}
}
