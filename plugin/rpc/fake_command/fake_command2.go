package fake_command

import (
	"errors"
	"fmt"

	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/flags"
)

type FakeCommand2 struct {
	Data string
	req  fakeReq2
	ui   terminal.UI
}

func init() {
	command_registry.Register(FakeCommand2{Data: "FakeCommand2 data", req: fakeReq2{}})
}

func (cmd FakeCommand2) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "fake-command2",
		Description: "Description for fake-command2 with bad requirement",
		Usage: []string{
			"Usage of fake-command",
		},
	}
}

func (cmd FakeCommand2) Requirements(_ requirements.Factory, _ flags.FlagContext) []requirements.Requirement {
	reqs := []requirements.Requirement{cmd.req}
	return reqs
}

func (cmd FakeCommand2) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.req.ui = deps.Ui
	cmd.ui = deps.Ui
	cmd.ui.Say("SetDependency() called, pluginCall " + fmt.Sprintf("%t", pluginCall))

	return cmd
}

func (cmd FakeCommand2) Execute(c flags.FlagContext) {
	cmd.ui.Say("Command Executed")
}

type fakeReq2 struct {
	ui terminal.UI
}

func (f fakeReq2) Execute() error {
	f.ui.Say("Requirement executed and failed")
	return errors.New("Requirement executed and failed")
}
