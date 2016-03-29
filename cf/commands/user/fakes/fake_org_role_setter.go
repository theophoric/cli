// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/theophoric/cf-cli/cf/command_registry"
	"github.com/theophoric/cf-cli/cf/commands/user"
	"github.com/theophoric/cf-cli/cf/requirements"
	"github.com/theophoric/cf-cli/flags"
)

type FakeOrgRoleSetter struct {
	MetaDataStub        func() command_registry.CommandMetadata
	metaDataMutex       sync.RWMutex
	metaDataArgsForCall []struct{}
	metaDataReturns     struct {
		result1 command_registry.CommandMetadata
	}
	SetDependencyStub        func(deps command_registry.Dependency, pluginCall bool) command_registry.Command
	setDependencyMutex       sync.RWMutex
	setDependencyArgsForCall []struct {
		deps       command_registry.Dependency
		pluginCall bool
	}
	setDependencyReturns struct {
		result1 command_registry.Command
	}
	RequirementsStub        func(requirementsFactory requirements.Factory, context flags.FlagContext) []requirements.Requirement
	requirementsMutex       sync.RWMutex
	requirementsArgsForCall []struct {
		requirementsFactory requirements.Factory
		context             flags.FlagContext
	}
	requirementsReturns struct {
		result1 []requirements.Requirement
	}
	ExecuteStub        func(context flags.FlagContext)
	executeMutex       sync.RWMutex
	executeArgsForCall []struct {
		context flags.FlagContext
	}
	SetOrgRoleStub        func(orgGuid string, role, userGuid, userName string) error
	setOrgRoleMutex       sync.RWMutex
	setOrgRoleArgsForCall []struct {
		orgGuid  string
		role     string
		userGuid string
		userName string
	}
	setOrgRoleReturns struct {
		result1 error
	}
}

func (fake *FakeOrgRoleSetter) MetaData() command_registry.CommandMetadata {
	fake.metaDataMutex.Lock()
	fake.metaDataArgsForCall = append(fake.metaDataArgsForCall, struct{}{})
	fake.metaDataMutex.Unlock()
	if fake.MetaDataStub != nil {
		return fake.MetaDataStub()
	} else {
		return fake.metaDataReturns.result1
	}
}

func (fake *FakeOrgRoleSetter) MetaDataCallCount() int {
	fake.metaDataMutex.RLock()
	defer fake.metaDataMutex.RUnlock()
	return len(fake.metaDataArgsForCall)
}

func (fake *FakeOrgRoleSetter) MetaDataReturns(result1 command_registry.CommandMetadata) {
	fake.MetaDataStub = nil
	fake.metaDataReturns = struct {
		result1 command_registry.CommandMetadata
	}{result1}
}

func (fake *FakeOrgRoleSetter) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	fake.setDependencyMutex.Lock()
	fake.setDependencyArgsForCall = append(fake.setDependencyArgsForCall, struct {
		deps       command_registry.Dependency
		pluginCall bool
	}{deps, pluginCall})
	fake.setDependencyMutex.Unlock()
	if fake.SetDependencyStub != nil {
		return fake.SetDependencyStub(deps, pluginCall)
	} else {
		return fake.setDependencyReturns.result1
	}
}

func (fake *FakeOrgRoleSetter) SetDependencyCallCount() int {
	fake.setDependencyMutex.RLock()
	defer fake.setDependencyMutex.RUnlock()
	return len(fake.setDependencyArgsForCall)
}

func (fake *FakeOrgRoleSetter) SetDependencyArgsForCall(i int) (command_registry.Dependency, bool) {
	fake.setDependencyMutex.RLock()
	defer fake.setDependencyMutex.RUnlock()
	return fake.setDependencyArgsForCall[i].deps, fake.setDependencyArgsForCall[i].pluginCall
}

func (fake *FakeOrgRoleSetter) SetDependencyReturns(result1 command_registry.Command) {
	fake.SetDependencyStub = nil
	fake.setDependencyReturns = struct {
		result1 command_registry.Command
	}{result1}
}

func (fake *FakeOrgRoleSetter) Requirements(requirementsFactory requirements.Factory, context flags.FlagContext) []requirements.Requirement {
	fake.requirementsMutex.Lock()
	fake.requirementsArgsForCall = append(fake.requirementsArgsForCall, struct {
		requirementsFactory requirements.Factory
		context             flags.FlagContext
	}{requirementsFactory, context})
	fake.requirementsMutex.Unlock()
	if fake.RequirementsStub != nil {
		return fake.RequirementsStub(requirementsFactory, context)
	} else {
		return fake.requirementsReturns.result1
	}
}

func (fake *FakeOrgRoleSetter) RequirementsCallCount() int {
	fake.requirementsMutex.RLock()
	defer fake.requirementsMutex.RUnlock()
	return len(fake.requirementsArgsForCall)
}

func (fake *FakeOrgRoleSetter) RequirementsArgsForCall(i int) (requirements.Factory, flags.FlagContext) {
	fake.requirementsMutex.RLock()
	defer fake.requirementsMutex.RUnlock()
	return fake.requirementsArgsForCall[i].requirementsFactory, fake.requirementsArgsForCall[i].context
}

func (fake *FakeOrgRoleSetter) RequirementsReturns(result1 []requirements.Requirement) {
	fake.RequirementsStub = nil
	fake.requirementsReturns = struct {
		result1 []requirements.Requirement
	}{result1}
}

func (fake *FakeOrgRoleSetter) Execute(context flags.FlagContext) {
	fake.executeMutex.Lock()
	fake.executeArgsForCall = append(fake.executeArgsForCall, struct {
		context flags.FlagContext
	}{context})
	fake.executeMutex.Unlock()
	if fake.ExecuteStub != nil {
		fake.ExecuteStub(context)
	}
}

func (fake *FakeOrgRoleSetter) ExecuteCallCount() int {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	return len(fake.executeArgsForCall)
}

func (fake *FakeOrgRoleSetter) ExecuteArgsForCall(i int) flags.FlagContext {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	return fake.executeArgsForCall[i].context
}

func (fake *FakeOrgRoleSetter) SetOrgRole(orgGuid string, role string, userGuid string, userName string) error {
	fake.setOrgRoleMutex.Lock()
	fake.setOrgRoleArgsForCall = append(fake.setOrgRoleArgsForCall, struct {
		orgGuid  string
		role     string
		userGuid string
		userName string
	}{orgGuid, role, userGuid, userName})
	fake.setOrgRoleMutex.Unlock()
	if fake.SetOrgRoleStub != nil {
		return fake.SetOrgRoleStub(orgGuid, role, userGuid, userName)
	} else {
		return fake.setOrgRoleReturns.result1
	}
}

func (fake *FakeOrgRoleSetter) SetOrgRoleCallCount() int {
	fake.setOrgRoleMutex.RLock()
	defer fake.setOrgRoleMutex.RUnlock()
	return len(fake.setOrgRoleArgsForCall)
}

func (fake *FakeOrgRoleSetter) SetOrgRoleArgsForCall(i int) (string, string, string, string) {
	fake.setOrgRoleMutex.RLock()
	defer fake.setOrgRoleMutex.RUnlock()
	return fake.setOrgRoleArgsForCall[i].orgGuid, fake.setOrgRoleArgsForCall[i].role, fake.setOrgRoleArgsForCall[i].userGuid, fake.setOrgRoleArgsForCall[i].userName
}

func (fake *FakeOrgRoleSetter) SetOrgRoleReturns(result1 error) {
	fake.SetOrgRoleStub = nil
	fake.setOrgRoleReturns = struct {
		result1 error
	}{result1}
}

var _ user.OrgRoleSetter = new(FakeOrgRoleSetter)
