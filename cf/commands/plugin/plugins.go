package plugin

import (
	"fmt"

	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/plugin_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
	"github.com/theophoric/cf-cli/utils"
)

type Plugins struct {
	ui     terminal.UI
	config plugin_config.PluginConfiguration
}

func init() {
	command_registry.Register(&Plugins{})
}

func (cmd *Plugins) MetaData() command_registry.CommandMetadata {
	fs := make(map[string]flags.FlagSet)
	fs["checksum"] = &flags.BoolFlag{Name: "checksum", Usage: T("Compute and show the sha1 value of the plugin binary file")}

	return command_registry.CommandMetadata{
		Name:        "plugins",
		Description: T("List all available plugin commands"),
		Usage: []string{
			T("CF_NAME plugins"),
		},
		Flags: fs,
	}
}

func (cmd *Plugins) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	usageReq := requirements.NewUsageRequirement(command_registry.CliCommandUsagePresenter(cmd),
		T("No argument required"),
		func() bool {
			return len(fc.Args()) != 0
		},
	)

	reqs := []requirements.Requirement{
		usageReq,
	}
	return reqs
}

func (cmd *Plugins) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.PluginConfig
	return cmd
}

func (cmd *Plugins) Execute(c flags.FlagContext) {
	var version string

	cmd.ui.Say(T("Listing Installed Plugins..."))

	plugins := cmd.config.Plugins()

	var table terminal.Table
	if c.Bool("checksum") {
		cmd.ui.Say(T("Computing sha1 for installed plugins, this may take a while ..."))
		table = terminal.NewTable(cmd.ui, []string{T("Plugin Name"), T("Version"), T("Command Name"), "sha1", T("Command Help")})
	} else {
		table = terminal.NewTable(cmd.ui, []string{T("Plugin Name"), T("Version"), T("Command Name"), T("Command Help")})
	}

	for pluginName, metadata := range plugins {
		if metadata.Version.Major == 0 && metadata.Version.Minor == 0 && metadata.Version.Build == 0 {
			version = "N/A"
		} else {
			version = fmt.Sprintf("%d.%d.%d", metadata.Version.Major, metadata.Version.Minor, metadata.Version.Build)
		}

		for _, command := range metadata.Commands {
			args := []string{pluginName, version}

			if command.Alias != "" {
				args = append(args, command.Name+", "+command.Alias)
			} else {
				args = append(args, command.Name)
			}

			if c.Bool("checksum") {
				checksum := utils.NewSha1Checksum(metadata.Location)
				sha1, err := checksum.ComputeFileSha1()
				if err != nil {
					args = append(args, "n/a")
				} else {
					args = append(args, fmt.Sprintf("%x", sha1))
				}
			}

			args = append(args, command.HelpText)
			table.Add(args...)
		}
	}

	cmd.ui.Ok()
	cmd.ui.Say("")

	table.Print()
}
