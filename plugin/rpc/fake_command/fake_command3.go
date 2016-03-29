package fake_command

import (
	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/flags"
)

type FakeCommand3 struct {
	Data string
}

func init() {
	command_registry.Register(FakeCommand3{Data: "FakeCommand3 data"})
}

func (cmd FakeCommand3) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "fake-command3",
		Description: "Description for fake-command3",
		Usage: []string{
			"Usage of fake-command3",
		},
	}
}

func (cmd FakeCommand3) Requirements(_ requirements.Factory, _ flags.FlagContext) []requirements.Requirement {
	reqs := []requirements.Requirement{}
	return reqs
}

func (cmd FakeCommand3) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	return cmd
}

func (cmd FakeCommand3) Execute(c flags.FlagContext) {
	panic("this is a test panic for cli_rpc_server_test (panic recovery)")
}
