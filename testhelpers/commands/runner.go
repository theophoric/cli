package commands

import (
	"fmt"
	"os"

	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/flags"
	testreq "github.com/theophoric/cf-cli/testhelpers/requirements"
	testterm "github.com/theophoric/cf-cli/testhelpers/terminal"
)

type RunCommandResult int

const (
	RunCommandResultSuccess            = iota
	RunCommandResultFailed             = iota
	RunCommandResultRequirementsFailed = iota
)

func RunCliCommand(cmdName string, args []string, requirementsFactory *testreq.FakeReqFactory, updateFunc func(bool), pluginCall bool) (passedRequirements bool) {
	updateFunc(pluginCall)
	cmd := command_registry.Commands.FindCommand(cmdName)
	context := flags.NewFlagContext(cmd.MetaData().Flags)
	context.SkipFlagParsing(cmd.MetaData().SkipFlagParsing)
	err := context.Parse(args...)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	defer func() {
		errMsg := recover()

		if errMsg != nil && errMsg != testterm.QuietPanic {
			panic(errMsg)
		}
	}()
	requirements := cmd.Requirements(requirementsFactory, context)

	for _, requirement := range requirements {
		if err = requirement.Execute(); err != nil {
			return false
		}
	}

	passedRequirements = true

	cmd.Execute(context)

	return
}

func RunCliCommandWithoutDependency(cmdName string, args []string, requirementsFactory *testreq.FakeReqFactory) (passedRequirements bool) {
	cmd := command_registry.Commands.FindCommand(cmdName)
	context := flags.NewFlagContext(cmd.MetaData().Flags)
	context.SkipFlagParsing(cmd.MetaData().SkipFlagParsing)
	err := context.Parse(args...)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	defer func() {
		errMsg := recover()

		if errMsg != nil && errMsg != testterm.QuietPanic {
			panic(errMsg)
		}
	}()

	requirements := cmd.Requirements(requirementsFactory, context)

	for _, requirement := range requirements {
		if err = requirement.Execute(); err != nil {
			return false
		}
	}

	cmd.Execute(context)

	return true
}
