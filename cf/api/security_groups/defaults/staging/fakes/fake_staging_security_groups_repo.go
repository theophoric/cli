// This file was generated by counterfeiter
package fakes

import (
	"sync"

	. "github.com/theophoric/cf-cli/cf/api/security_groups/defaults"
	. "github.com/theophoric/cf-cli/cf/api/security_groups/defaults/staging"
	"github.com/theophoric/cf-cli/cf/models"
)

type FakeStagingSecurityGroupsRepo struct {
	BindToStagingSetStub        func(string) error
	bindToStagingSetMutex       sync.RWMutex
	bindToStagingSetArgsForCall []struct {
		arg1 string
	}
	bindToStagingSetReturns struct {
		result1 error
	}
	ListStub        func() ([]models.SecurityGroupFields, error)
	listMutex       sync.RWMutex
	listArgsForCall []struct{}
	listReturns     struct {
		result1 []models.SecurityGroupFields
		result2 error
	}
	UnbindFromStagingSetStub        func(string) error
	unbindFromStagingSetMutex       sync.RWMutex
	unbindFromStagingSetArgsForCall []struct {
		arg1 string
	}
	unbindFromStagingSetReturns struct {
		result1 error
	}
}

func (fake *FakeStagingSecurityGroupsRepo) BindToStagingSet(arg1 string) error {
	fake.bindToStagingSetMutex.Lock()
	defer fake.bindToStagingSetMutex.Unlock()
	fake.bindToStagingSetArgsForCall = append(fake.bindToStagingSetArgsForCall, struct {
		arg1 string
	}{arg1})
	if fake.BindToStagingSetStub != nil {
		return fake.BindToStagingSetStub(arg1)
	} else {
		return fake.bindToStagingSetReturns.result1
	}
}

func (fake *FakeStagingSecurityGroupsRepo) BindToStagingSetCallCount() int {
	fake.bindToStagingSetMutex.RLock()
	defer fake.bindToStagingSetMutex.RUnlock()
	return len(fake.bindToStagingSetArgsForCall)
}

func (fake *FakeStagingSecurityGroupsRepo) BindToStagingSetArgsForCall(i int) string {
	fake.bindToStagingSetMutex.RLock()
	defer fake.bindToStagingSetMutex.RUnlock()
	return fake.bindToStagingSetArgsForCall[i].arg1
}

func (fake *FakeStagingSecurityGroupsRepo) BindToStagingSetReturns(result1 error) {
	fake.bindToStagingSetReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStagingSecurityGroupsRepo) List() ([]models.SecurityGroupFields, error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.listArgsForCall = append(fake.listArgsForCall, struct{}{})
	if fake.ListStub != nil {
		return fake.ListStub()
	} else {
		return fake.listReturns.result1, fake.listReturns.result2
	}
}

func (fake *FakeStagingSecurityGroupsRepo) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *FakeStagingSecurityGroupsRepo) ListReturns(result1 []models.SecurityGroupFields, result2 error) {
	fake.listReturns = struct {
		result1 []models.SecurityGroupFields
		result2 error
	}{result1, result2}
}

func (fake *FakeStagingSecurityGroupsRepo) UnbindFromStagingSet(arg1 string) error {
	fake.unbindFromStagingSetMutex.Lock()
	defer fake.unbindFromStagingSetMutex.Unlock()
	fake.unbindFromStagingSetArgsForCall = append(fake.unbindFromStagingSetArgsForCall, struct {
		arg1 string
	}{arg1})
	if fake.UnbindFromStagingSetStub != nil {
		return fake.UnbindFromStagingSetStub(arg1)
	} else {
		return fake.unbindFromStagingSetReturns.result1
	}
}

func (fake *FakeStagingSecurityGroupsRepo) UnbindFromStagingSetCallCount() int {
	fake.unbindFromStagingSetMutex.RLock()
	defer fake.unbindFromStagingSetMutex.RUnlock()
	return len(fake.unbindFromStagingSetArgsForCall)
}

func (fake *FakeStagingSecurityGroupsRepo) UnbindFromStagingSetArgsForCall(i int) string {
	fake.unbindFromStagingSetMutex.RLock()
	defer fake.unbindFromStagingSetMutex.RUnlock()
	return fake.unbindFromStagingSetArgsForCall[i].arg1
}

func (fake *FakeStagingSecurityGroupsRepo) UnbindFromStagingSetReturns(result1 error) {
	fake.unbindFromStagingSetReturns = struct {
		result1 error
	}{result1}
}

var _ StagingSecurityGroupsRepo = new(FakeStagingSecurityGroupsRepo)
