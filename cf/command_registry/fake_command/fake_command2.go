package fake_command

import (
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/flags"
)

type FakeCommand2 struct {
	Data string
}

func (cmd FakeCommand2) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "fake-command2",
		ShortName:   "fc2",
		Description: "Description for fake-command2",
		Usage: []string{
			"Usage of fake-command2",
		},
	}
}

func (cmd FakeCommand2) Requirements(_ requirements.Factory, _ flags.FlagContext) []requirements.Requirement {
	return []requirements.Requirement{}
}

func (cmd FakeCommand2) SetDependency(deps command_registry.Dependency, _ bool) command_registry.Command {
	return cmd
}

func (cmd FakeCommand2) Execute(c flags.FlagContext) {
}
