// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/plugin/rpc"
)

type FakeCommandRunner struct {
	CommandStub        func([]string, command_registry.Dependency, bool) error
	commandMutex       sync.RWMutex
	commandArgsForCall []struct {
		arg1 []string
		arg2 command_registry.Dependency
		arg3 bool
	}
	commandReturns struct {
		result1 error
	}
}

func (fake *FakeCommandRunner) Command(arg1 []string, arg2 command_registry.Dependency, arg3 bool) error {
	fake.commandMutex.Lock()
	fake.commandArgsForCall = append(fake.commandArgsForCall, struct {
		arg1 []string
		arg2 command_registry.Dependency
		arg3 bool
	}{arg1, arg2, arg3})
	fake.commandMutex.Unlock()
	if fake.CommandStub != nil {
		return fake.CommandStub(arg1, arg2, arg3)
	} else {
		return fake.commandReturns.result1
	}
}

func (fake *FakeCommandRunner) CommandCallCount() int {
	fake.commandMutex.RLock()
	defer fake.commandMutex.RUnlock()
	return len(fake.commandArgsForCall)
}

func (fake *FakeCommandRunner) CommandArgsForCall(i int) ([]string, command_registry.Dependency, bool) {
	fake.commandMutex.RLock()
	defer fake.commandMutex.RUnlock()
	return fake.commandArgsForCall[i].arg1, fake.commandArgsForCall[i].arg2, fake.commandArgsForCall[i].arg3
}

func (fake *FakeCommandRunner) CommandReturns(result1 error) {
	fake.CommandStub = nil
	fake.commandReturns = struct {
		result1 error
	}{result1}
}

var _ rpc.CommandRunner = new(FakeCommandRunner)
