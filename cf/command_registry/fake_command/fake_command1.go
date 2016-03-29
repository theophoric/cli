package fake_command

import (
	"fmt"

	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/flags"
)

type FakeCommand1 struct {
	Data string
}

func init() {
	command_registry.Register(FakeCommand1{Data: "FakeCommand1 data"})
}

func (cmd FakeCommand1) MetaData() command_registry.CommandMetadata {
	fs := make(map[string]flags.FlagSet)
	fs["f"] = &flags.BoolFlag{ShortName: "f", Usage: "Usage for BoolFlag"}
	fs["boolFlag"] = &flags.BoolFlag{Name: "BoolFlag", Usage: "Usage for BoolFlag"}
	fs["intFlag"] = &flags.IntFlag{Name: "intFlag", Usage: "Usage for intFlag"}

	return command_registry.CommandMetadata{
		Name:        "fake-command",
		ShortName:   "fc1",
		Description: "Description for fake-command",
		Usage: []string{
			"CF_NAME Usage of fake-command",
		},
		Flags: fs,
	}
}

func (cmd FakeCommand1) Requirements(_ requirements.Factory, _ flags.FlagContext) []requirements.Requirement {
	return []requirements.Requirement{}
}

func (cmd FakeCommand1) SetDependency(deps command_registry.Dependency, _ bool) command_registry.Command {
	return cmd
}

func (cmd FakeCommand1) Execute(c flags.FlagContext) {
	fmt.Println("This is fake-command")
}
