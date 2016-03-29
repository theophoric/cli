package commands

import (
	"fmt"
	"strings"

	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/formatters"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/cf/v3/models"
	"github.com/theophoric/cf-cli/cf/v3/repository"
	"github.com/theophoric/cf-cli/flags"

	. "github.com/theophoric/cf-cli/cf/i18n"
)

type V3Apps struct {
	ui         terminal.UI
	config     core_config.ReadWriter
	repository repository.Repository
}

func init() {
	command_registry.Register(&V3Apps{})
}

func (c *V3Apps) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "v3apps",
		Description: T("List all apps in the target space"),
		Usage: []string{
			"CF_NAME v3apps",
		},
		Hidden: true,
	}
}

func (c *V3Apps) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	usageReq := requirements.NewUsageRequirement(command_registry.CliCommandUsagePresenter(c),
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

func (c *V3Apps) SetDependency(deps command_registry.Dependency, _ bool) command_registry.Command {
	c.ui = deps.Ui
	c.config = deps.Config
	c.repository = deps.RepoLocator.GetV3Repository()

	return c
}

func (c *V3Apps) Execute(fc flags.FlagContext) {
	applications, err := c.repository.GetApplications()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	processes := make([][]models.V3Process, len(applications))
	routes := make([][]models.V3Route, len(applications))

	for i, app := range applications {
		ps, err := c.repository.GetProcesses(app.Links.Processes.Href)
		if err != nil {
			c.ui.Failed(err.Error())
		}
		processes[i] = ps

		rs, err := c.repository.GetRoutes(app.Links.Routes.Href)
		if err != nil {
			c.ui.Failed(err.Error())
		}
		routes[i] = rs
	}

	table := terminal.NewTable(c.ui, []string{T("name"), T("requested state"), T("instances"), T("memory"), T("disk"), T("urls")})

	for i := range applications {
		c.addRow(table, applications[i], processes[i], routes[i])
	}

	table.Print()
}

func (c *V3Apps) addRow(
	table terminal.Table,
	application models.V3Application,
	processes []models.V3Process,
	routes []models.V3Route,
) {
	var webProcess models.V3Process
	for i := range processes {
		if processes[i].Type == "web" {
			webProcess = processes[i]
		}
	}

	var appRoutes []string
	for _, route := range routes {
		appRoutes = append(appRoutes, route.Host+route.Path)
	}

	table.Add(
		application.Name,
		strings.ToLower(application.DesiredState),
		fmt.Sprintf("%d", application.TotalDesiredInstances),
		formatters.ByteSize(webProcess.MemoryInMB*formatters.MEGABYTE),
		formatters.ByteSize(webProcess.DiskInMB*formatters.MEGABYTE),
		strings.Join(appRoutes, ", "),
	)
}
