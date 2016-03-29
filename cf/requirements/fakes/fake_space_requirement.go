// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/requirements"
)

type FakeSpaceRequirement struct {
	ExecuteStub        func() error
	executeMutex       sync.RWMutex
	executeArgsForCall []struct{}
	executeReturns     struct {
		result1 error
	}
	SetSpaceNameStub        func(string)
	setSpaceNameMutex       sync.RWMutex
	setSpaceNameArgsForCall []struct {
		arg1 string
	}
	GetSpaceStub        func() models.Space
	getSpaceMutex       sync.RWMutex
	getSpaceArgsForCall []struct{}
	getSpaceReturns     struct {
		result1 models.Space
	}
}

func (fake *FakeSpaceRequirement) Execute() error {
	fake.executeMutex.Lock()
	fake.executeArgsForCall = append(fake.executeArgsForCall, struct{}{})
	fake.executeMutex.Unlock()
	if fake.ExecuteStub != nil {
		return fake.ExecuteStub()
	} else {
		return fake.executeReturns.result1
	}
}

func (fake *FakeSpaceRequirement) ExecuteCallCount() int {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	return len(fake.executeArgsForCall)
}

func (fake *FakeSpaceRequirement) ExecuteReturns(result1 error) {
	fake.ExecuteStub = nil
	fake.executeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSpaceRequirement) SetSpaceName(arg1 string) {
	fake.setSpaceNameMutex.Lock()
	fake.setSpaceNameArgsForCall = append(fake.setSpaceNameArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.setSpaceNameMutex.Unlock()
	if fake.SetSpaceNameStub != nil {
		fake.SetSpaceNameStub(arg1)
	}
}

func (fake *FakeSpaceRequirement) SetSpaceNameCallCount() int {
	fake.setSpaceNameMutex.RLock()
	defer fake.setSpaceNameMutex.RUnlock()
	return len(fake.setSpaceNameArgsForCall)
}

func (fake *FakeSpaceRequirement) SetSpaceNameArgsForCall(i int) string {
	fake.setSpaceNameMutex.RLock()
	defer fake.setSpaceNameMutex.RUnlock()
	return fake.setSpaceNameArgsForCall[i].arg1
}

func (fake *FakeSpaceRequirement) GetSpace() models.Space {
	fake.getSpaceMutex.Lock()
	fake.getSpaceArgsForCall = append(fake.getSpaceArgsForCall, struct{}{})
	fake.getSpaceMutex.Unlock()
	if fake.GetSpaceStub != nil {
		return fake.GetSpaceStub()
	} else {
		return fake.getSpaceReturns.result1
	}
}

func (fake *FakeSpaceRequirement) GetSpaceCallCount() int {
	fake.getSpaceMutex.RLock()
	defer fake.getSpaceMutex.RUnlock()
	return len(fake.getSpaceArgsForCall)
}

func (fake *FakeSpaceRequirement) GetSpaceReturns(result1 models.Space) {
	fake.GetSpaceStub = nil
	fake.getSpaceReturns = struct {
		result1 models.Space
	}{result1}
}

var _ requirements.SpaceRequirement = new(FakeSpaceRequirement)
