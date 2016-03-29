// This file was generated by counterfeiter
package fakes

import (
	"sync"

	. "github.com/theophoric/cf-cli/cf/api/security_groups/defaults"
	. "github.com/theophoric/cf-cli/cf/api/security_groups/defaults/running"
	"github.com/theophoric/cf-cli/cf/models"
)

type FakeRunningSecurityGroupsRepo struct {
	BindToRunningSetStub        func(string) error
	bindToRunningSetMutex       sync.RWMutex
	bindToRunningSetArgsForCall []struct {
		arg1 string
	}
	bindToRunningSetReturns struct {
		result1 error
	}
	ListStub        func() ([]models.SecurityGroupFields, error)
	listMutex       sync.RWMutex
	listArgsForCall []struct{}
	listReturns     struct {
		result1 []models.SecurityGroupFields
		result2 error
	}
	UnbindFromRunningSetStub        func(string) error
	unbindFromRunningSetMutex       sync.RWMutex
	unbindFromRunningSetArgsForCall []struct {
		arg1 string
	}
	unbindFromRunningSetReturns struct {
		result1 error
	}
}

func (fake *FakeRunningSecurityGroupsRepo) BindToRunningSet(arg1 string) error {
	fake.bindToRunningSetMutex.Lock()
	defer fake.bindToRunningSetMutex.Unlock()
	fake.bindToRunningSetArgsForCall = append(fake.bindToRunningSetArgsForCall, struct {
		arg1 string
	}{arg1})
	if fake.BindToRunningSetStub != nil {
		return fake.BindToRunningSetStub(arg1)
	} else {
		return fake.bindToRunningSetReturns.result1
	}
}

func (fake *FakeRunningSecurityGroupsRepo) BindToRunningSetCallCount() int {
	fake.bindToRunningSetMutex.RLock()
	defer fake.bindToRunningSetMutex.RUnlock()
	return len(fake.bindToRunningSetArgsForCall)
}

func (fake *FakeRunningSecurityGroupsRepo) BindToRunningSetArgsForCall(i int) string {
	fake.bindToRunningSetMutex.RLock()
	defer fake.bindToRunningSetMutex.RUnlock()
	return fake.bindToRunningSetArgsForCall[i].arg1
}

func (fake *FakeRunningSecurityGroupsRepo) BindToRunningSetReturns(result1 error) {
	fake.bindToRunningSetReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRunningSecurityGroupsRepo) List() ([]models.SecurityGroupFields, error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.listArgsForCall = append(fake.listArgsForCall, struct{}{})
	if fake.ListStub != nil {
		return fake.ListStub()
	} else {
		return fake.listReturns.result1, fake.listReturns.result2
	}
}

func (fake *FakeRunningSecurityGroupsRepo) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *FakeRunningSecurityGroupsRepo) ListReturns(result1 []models.SecurityGroupFields, result2 error) {
	fake.listReturns = struct {
		result1 []models.SecurityGroupFields
		result2 error
	}{result1, result2}
}

func (fake *FakeRunningSecurityGroupsRepo) UnbindFromRunningSet(arg1 string) error {
	fake.unbindFromRunningSetMutex.Lock()
	defer fake.unbindFromRunningSetMutex.Unlock()
	fake.unbindFromRunningSetArgsForCall = append(fake.unbindFromRunningSetArgsForCall, struct {
		arg1 string
	}{arg1})
	if fake.UnbindFromRunningSetStub != nil {
		return fake.UnbindFromRunningSetStub(arg1)
	} else {
		return fake.unbindFromRunningSetReturns.result1
	}
}

func (fake *FakeRunningSecurityGroupsRepo) UnbindFromRunningSetCallCount() int {
	fake.unbindFromRunningSetMutex.RLock()
	defer fake.unbindFromRunningSetMutex.RUnlock()
	return len(fake.unbindFromRunningSetArgsForCall)
}

func (fake *FakeRunningSecurityGroupsRepo) UnbindFromRunningSetArgsForCall(i int) string {
	fake.unbindFromRunningSetMutex.RLock()
	defer fake.unbindFromRunningSetMutex.RUnlock()
	return fake.unbindFromRunningSetArgsForCall[i].arg1
}

func (fake *FakeRunningSecurityGroupsRepo) UnbindFromRunningSetReturns(result1 error) {
	fake.unbindFromRunningSetReturns = struct {
		result1 error
	}{result1}
}

var _ RunningSecurityGroupsRepo = new(FakeRunningSecurityGroupsRepo)
