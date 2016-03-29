package commands

import (
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/configuration/core_config"
	. "github.com/theophoric/cf-cli/cf/i18n"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type Logout struct {
	ui     terminal.UI
	config core_config.ReadWriter
}

func init() {
	command_registry.Register(&Logout{})
}

func (cmd *Logout) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "logout",
		ShortName:   "lo",
		Description: T("Log user out"),
		Usage: []string{
			T("CF_NAME logout"),
		},
	}
}

func (cmd *Logout) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) []requirements.Requirement {
	reqs := []requirements.Requirement{}
	return reqs
}

func (cmd *Logout) SetDependency(deps command_registry.Dependency, _ bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	return cmd
}

func (cmd *Logout) Execute(c flags.FlagContext) {
	cmd.ui.Say(T("Logging out..."))
	cmd.config.ClearSession()
	cmd.ui.Ok()
}
