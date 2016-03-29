package application

import (
	"github.com/theophoric/cf-cli/cf/api/app_files"
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type Files struct {
	ui           terminal.UI
	config       core_config.Reader
	appFilesRepo app_files.AppFilesRepository
	appReq       requirements.DEAApplicationRequirement
}

func init() {
	command_registry.Register(&Files{})
}

func (cmd *Files) MetaData() command_registry.CommandMetadata {
	fs := make(map[string]flags.FlagSet)
	fs["i"] = &flags.IntFlag{ShortName: "i", Usage: T("Instance")}

	return command_registry.CommandMetadata{
		Name:        "files",
		ShortName:   "f",
		Description: T("Print out a list of files in a directory or the contents of a specific file of an app running on the DEA backend"),
		Usage: []string{
			T(`CF_NAME files APP_NAME [PATH] [-i INSTANCE]
			
TIP:
  To list and inspect files of an app running on the Diego backend, use 'CF_NAME ssh'`),
		},
		Flags: fs,
	}
}

func (cmd *Files) Requirements(requirementsFactory requirements.Factory, c flags.FlagContext) []requirements.Requirement {
	if len(c.Args()) < 1 || len(c.Args()) > 2 {
		cmd.ui.Failed(T("Incorrect Usage. Requires an argument\n\n") + command_registry.Commands.CommandUsage("files"))
	}

	cmd.appReq = requirementsFactory.NewDEAApplicationRequirement(c.Args()[0])

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedSpaceRequirement(),
		cmd.appReq,
	}

	return reqs
}

func (cmd *Files) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.appFilesRepo = deps.RepoLocator.GetAppFilesRepository()
	return cmd
}

func (cmd *Files) Execute(c flags.FlagContext) {
	app := cmd.appReq.GetApplication()

	var instance int
	if c.IsSet("i") {
		instance = c.Int("i")
		if instance < 0 {
			cmd.ui.Failed(T("Invalid instance: {{.Instance}}\nInstance must be a positive integer",
				map[string]interface{}{
					"Instance": instance,
				}))
		}
		if instance >= app.InstanceCount {
			cmd.ui.Failed(T("Invalid instance: {{.Instance}}\nInstance must be less than {{.InstanceCount}}",
				map[string]interface{}{
					"Instance":      instance,
					"InstanceCount": app.InstanceCount,
				}))
		}
	}

	cmd.ui.Say(T("Getting files for app {{.AppName}} in org {{.OrgName}} / space {{.SpaceName}} as {{.Username}}...",
		map[string]interface{}{
			"AppName":   terminal.EntityNameColor(app.Name),
			"OrgName":   terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
			"SpaceName": terminal.EntityNameColor(cmd.config.SpaceFields().Name),
			"Username":  terminal.EntityNameColor(cmd.config.Username())}))

	path := "/"
	if len(c.Args()) > 1 {
		path = c.Args()[1]
	}

	list, apiErr := cmd.appFilesRepo.ListFiles(app.Guid, instance, path)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
	cmd.ui.Say("")

	if list == "" {
		cmd.ui.Say("No files found")
	} else {
		cmd.ui.Say("%s", list)
	}
}
