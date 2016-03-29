// This file was generated by counterfeiter
package fakes

import (
	. "github.com/theophoric/cf-cli/cf/api/environment_variable_groups"
	"github.com/theophoric/cf-cli/cf/models"

	"sync"
)

type FakeEnvironmentVariableGroupsRepository struct {
	ListRunningStub        func() (variables []models.EnvironmentVariable, apiErr error)
	listRunningMutex       sync.RWMutex
	listRunningArgsForCall []struct{}
	listRunningReturns     struct {
		result1 []models.EnvironmentVariable
		result2 error
	}
	ListStagingStub        func() (variables []models.EnvironmentVariable, apiErr error)
	listStagingMutex       sync.RWMutex
	listStagingArgsForCall []struct{}
	listStagingReturns     struct {
		result1 []models.EnvironmentVariable
		result2 error
	}
	SetStagingStub        func(string) error
	setStagingMutex       sync.RWMutex
	setStagingArgsForCall []struct {
		arg1 string
	}
	setStagingReturns struct {
		result1 error
	}
	SetRunningStub        func(string) error
	setRunningMutex       sync.RWMutex
	setRunningArgsForCall []struct {
		arg1 string
	}
	setRunningReturns struct {
		result1 error
	}
}

func (fake *FakeEnvironmentVariableGroupsRepository) ListRunning() (variables []models.EnvironmentVariable, apiErr error) {
	fake.listRunningMutex.Lock()
	defer fake.listRunningMutex.Unlock()
	fake.listRunningArgsForCall = append(fake.listRunningArgsForCall, struct{}{})
	if fake.ListRunningStub != nil {
		return fake.ListRunningStub()
	} else {
		return fake.listRunningReturns.result1, fake.listRunningReturns.result2
	}
}

func (fake *FakeEnvironmentVariableGroupsRepository) ListRunningCallCount() int {
	fake.listRunningMutex.RLock()
	defer fake.listRunningMutex.RUnlock()
	return len(fake.listRunningArgsForCall)
}

func (fake *FakeEnvironmentVariableGroupsRepository) ListRunningReturns(result1 []models.EnvironmentVariable, result2 error) {
	fake.listRunningReturns = struct {
		result1 []models.EnvironmentVariable
		result2 error
	}{result1, result2}
}

func (fake *FakeEnvironmentVariableGroupsRepository) ListStaging() (variables []models.EnvironmentVariable, apiErr error) {
	fake.listStagingMutex.Lock()
	defer fake.listStagingMutex.Unlock()
	fake.listStagingArgsForCall = append(fake.listStagingArgsForCall, struct{}{})
	if fake.ListStagingStub != nil {
		return fake.ListStagingStub()
	} else {
		return fake.listStagingReturns.result1, fake.listStagingReturns.result2
	}
}

func (fake *FakeEnvironmentVariableGroupsRepository) ListStagingCallCount() int {
	fake.listStagingMutex.RLock()
	defer fake.listStagingMutex.RUnlock()
	return len(fake.listStagingArgsForCall)
}

func (fake *FakeEnvironmentVariableGroupsRepository) ListStagingReturns(result1 []models.EnvironmentVariable, result2 error) {
	fake.listStagingReturns = struct {
		result1 []models.EnvironmentVariable
		result2 error
	}{result1, result2}
}

func (fake *FakeEnvironmentVariableGroupsRepository) SetStaging(arg1 string) error {
	fake.setStagingMutex.Lock()
	defer fake.setStagingMutex.Unlock()
	fake.setStagingArgsForCall = append(fake.setStagingArgsForCall, struct {
		arg1 string
	}{arg1})
	if fake.SetStagingStub != nil {
		return fake.SetStagingStub(arg1)
	} else {
		return fake.setStagingReturns.result1
	}
}

func (fake *FakeEnvironmentVariableGroupsRepository) SetStagingCallCount() int {
	fake.setStagingMutex.RLock()
	defer fake.setStagingMutex.RUnlock()
	return len(fake.setStagingArgsForCall)
}

func (fake *FakeEnvironmentVariableGroupsRepository) SetStagingArgsForCall(i int) string {
	fake.setStagingMutex.RLock()
	defer fake.setStagingMutex.RUnlock()
	return fake.setStagingArgsForCall[i].arg1
}

func (fake *FakeEnvironmentVariableGroupsRepository) SetStagingReturns(result1 error) {
	fake.setStagingReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeEnvironmentVariableGroupsRepository) SetRunning(arg1 string) error {
	fake.setRunningMutex.Lock()
	defer fake.setRunningMutex.Unlock()
	fake.setRunningArgsForCall = append(fake.setRunningArgsForCall, struct {
		arg1 string
	}{arg1})
	if fake.SetRunningStub != nil {
		return fake.SetRunningStub(arg1)
	} else {
		return fake.setRunningReturns.result1
	}
}

func (fake *FakeEnvironmentVariableGroupsRepository) SetRunningCallCount() int {
	fake.setRunningMutex.RLock()
	defer fake.setRunningMutex.RUnlock()
	return len(fake.setRunningArgsForCall)
}

func (fake *FakeEnvironmentVariableGroupsRepository) SetRunningArgsForCall(i int) string {
	fake.setRunningMutex.RLock()
	defer fake.setRunningMutex.RUnlock()
	return fake.setRunningArgsForCall[i].arg1
}

func (fake *FakeEnvironmentVariableGroupsRepository) SetRunningReturns(result1 error) {
	fake.setRunningReturns = struct {
		result1 error
	}{result1}
}

var _ EnvironmentVariableGroupsRepository = new(FakeEnvironmentVariableGroupsRepository)
