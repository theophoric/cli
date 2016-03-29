package application

import (
	"strconv"
	"strings"

	"github.com/theophoric/cf-cli/cf/command_registry"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/flags"
	"github.com/theophoric/cf-cli/plugin/models"

	"github.com/theophoric/cf-cli/cf/api"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	"github.com/theophoric/cf-cli/cf/formatters"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/cf/ui_helpers"
)

type ListApps struct {
	ui             terminal.UI
	config         core_config.Reader
	appSummaryRepo api.AppSummaryRepository

	pluginAppModels *[]plugin_models.GetAppsModel
	pluginCall      bool
}

func init() {
	command_registry.Register(&ListApps{})
}

func (cmd *ListApps) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "apps",
		ShortName:   "a",
		Description: T("List all apps in the target space"),
		Usage: []string{
			"CF_NAME apps",
		},
	}
}

func (cmd *ListApps) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
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

func (cmd *ListApps) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.appSummaryRepo = deps.RepoLocator.GetAppSummaryRepository()
	cmd.pluginAppModels = deps.PluginModels.AppsSummary
	cmd.pluginCall = pluginCall
	return cmd
}

func (cmd *ListApps) Execute(c flags.FlagContext) {
	cmd.ui.Say(T("Getting apps in org {{.OrgName}} / space {{.SpaceName}} as {{.Username}}...",
		map[string]interface{}{
			"OrgName":   terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
			"SpaceName": terminal.EntityNameColor(cmd.config.SpaceFields().Name),
			"Username":  terminal.EntityNameColor(cmd.config.Username())}))

	apps, apiErr := cmd.appSummaryRepo.GetSummariesInCurrentSpace()

	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
	cmd.ui.Say("")

	if len(apps) == 0 {
		cmd.ui.Say(T("No apps found"))
		return
	}

	table := terminal.NewTable(cmd.ui, []string{
		T("name"),
		T("requested state"),
		T("instances"),
		T("memory"),
		T("disk"),
		T("app ports"),
		T("urls"),
	})

	for _, application := range apps {
		var urls []string
		for _, route := range application.Routes {
			urls = append(urls, route.URL())
		}

		appPorts := make([]string, len(application.AppPorts))
		for i, p := range application.AppPorts {
			appPorts[i] = strconv.Itoa(p)
		}

		table.Add(
			application.Name,
			ui_helpers.ColoredAppState(application.ApplicationFields),
			ui_helpers.ColoredAppInstances(application.ApplicationFields),
			formatters.ByteSize(application.Memory*formatters.MEGABYTE),
			formatters.ByteSize(application.DiskQuota*formatters.MEGABYTE),
			strings.Join(appPorts, ", "),
			strings.Join(urls, ", "),
		)
	}

	table.Print()

	if cmd.pluginCall {
		cmd.populatePluginModel(apps)
	}
}

func (cmd *ListApps) populatePluginModel(apps []models.Application) {
	for _, app := range apps {
		appModel := plugin_models.GetAppsModel{}
		appModel.Name = app.Name
		appModel.Guid = app.Guid
		appModel.TotalInstances = app.InstanceCount
		appModel.RunningInstances = app.RunningInstances
		appModel.Memory = app.Memory
		appModel.State = app.State
		appModel.DiskQuota = app.DiskQuota
		appModel.AppPorts = app.AppPorts

		*(cmd.pluginAppModels) = append(*(cmd.pluginAppModels), appModel)

		for _, route := range app.Routes {
			r := plugin_models.GetAppsRouteSummary{}
			r.Host = route.Host
			r.Guid = route.Guid
			r.Domain.Guid = route.Domain.Guid
			r.Domain.Name = route.Domain.Name
			r.Domain.OwningOrganizationGuid = route.Domain.OwningOrganizationGuid
			r.Domain.Shared = route.Domain.Shared

			(*(cmd.pluginAppModels))[len(*(cmd.pluginAppModels))-1].Routes = append((*(cmd.pluginAppModels))[len(*(cmd.pluginAppModels))-1].Routes, r)
		}

	}
}
