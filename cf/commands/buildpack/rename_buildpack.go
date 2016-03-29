package buildpack

import (
	"github.com/theophoric/cf-cli/cf/api"
	"github.com/theophoric/cf-cli/cf/command_registry"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type RenameBuildpack struct {
	ui            terminal.UI
	buildpackRepo api.BuildpackRepository
}

func init() {
	command_registry.Register(&RenameBuildpack{})
}

func (cmd *RenameBuildpack) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "rename-buildpack",
		Description: T("Rename a buildpack"),
		Usage: []string{
			T("CF_NAME rename-buildpack BUILDPACK_NAME NEW_BUILDPACK_NAME"),
		},
	}
}

func (cmd *RenameBuildpack) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	if len(fc.Args()) != 2 {
		cmd.ui.Failed(T("Incorrect Usage. Requires BUILDPACK_NAME, NEW_BUILDPACK_NAME as arguments\n\n") + command_registry.Commands.CommandUsage("rename-buildpack"))
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}

	return reqs
}

func (cmd *RenameBuildpack) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.buildpackRepo = deps.RepoLocator.GetBuildpackRepository()
	return cmd
}

func (cmd *RenameBuildpack) Execute(c flags.FlagContext) {
	buildpackName := c.Args()[0]
	newBuildpackName := c.Args()[1]

	cmd.ui.Say(T("Renaming buildpack {{.OldBuildpackName}} to {{.NewBuildpackName}}...", map[string]interface{}{"OldBuildpackName": terminal.EntityNameColor(buildpackName), "NewBuildpackName": terminal.EntityNameColor(newBuildpackName)}))

	buildpack, apiErr := cmd.buildpackRepo.FindByName(buildpackName)

	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
	}

	buildpack.Name = newBuildpackName
	buildpack, apiErr = cmd.buildpackRepo.Update(buildpack)
	if apiErr != nil {
		cmd.ui.Failed(T("Error renaming buildpack {{.Name}}\n{{.Error}}", map[string]interface{}{
			"Name":  terminal.EntityNameColor(buildpack.Name),
			"Error": apiErr.Error(),
		}))
	}

	cmd.ui.Ok()
}
