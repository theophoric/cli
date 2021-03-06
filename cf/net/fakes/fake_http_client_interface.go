// This file was generated by counterfeiter
package fakes

import (
	_ "crypto/sha512"
	"net/http"
	"sync"

	"github.com/theophoric/cf-cli/cf/net"
)

type FakeHttpClientInterface struct {
	DoStub        func(*http.Request) (*http.Response, error)
	doMutex       sync.RWMutex
	doArgsForCall []struct {
		arg1 *http.Request
	}
	doReturns struct {
		result1 *http.Response
		result2 error
	}
	DumpRequestStub        func(*http.Request)
	dumpRequestMutex       sync.RWMutex
	dumpRequestArgsForCall []struct {
		arg1 *http.Request
	}
	DumpResponseStub        func(*http.Response)
	dumpResponseMutex       sync.RWMutex
	dumpResponseArgsForCall []struct {
		arg1 *http.Response
	}
	ExecuteCheckRedirectStub        func(req *http.Request, via []*http.Request) error
	executeCheckRedirectMutex       sync.RWMutex
	executeCheckRedirectArgsForCall []struct {
		req *http.Request
		via []*http.Request
	}
	executeCheckRedirectReturns struct {
		result1 error
	}
}

func (fake *FakeHttpClientInterface) Do(arg1 *http.Request) (*http.Response, error) {
	fake.doMutex.Lock()
	fake.doArgsForCall = append(fake.doArgsForCall, struct {
		arg1 *http.Request
	}{arg1})
	fake.doMutex.Unlock()
	if fake.DoStub != nil {
		return fake.DoStub(arg1)
	} else {
		return fake.doReturns.result1, fake.doReturns.result2
	}
}

func (fake *FakeHttpClientInterface) DoCallCount() int {
	fake.doMutex.RLock()
	defer fake.doMutex.RUnlock()
	return len(fake.doArgsForCall)
}

func (fake *FakeHttpClientInterface) DoArgsForCall(i int) *http.Request {
	fake.doMutex.RLock()
	defer fake.doMutex.RUnlock()
	return fake.doArgsForCall[i].arg1
}

func (fake *FakeHttpClientInterface) DoReturns(result1 *http.Response, result2 error) {
	fake.DoStub = nil
	fake.doReturns = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeHttpClientInterface) DumpRequest(arg1 *http.Request) {
	fake.dumpRequestMutex.Lock()
	fake.dumpRequestArgsForCall = append(fake.dumpRequestArgsForCall, struct {
		arg1 *http.Request
	}{arg1})
	fake.dumpRequestMutex.Unlock()
	if fake.DumpRequestStub != nil {
		fake.DumpRequestStub(arg1)
	}
}

func (fake *FakeHttpClientInterface) DumpRequestCallCount() int {
	fake.dumpRequestMutex.RLock()
	defer fake.dumpRequestMutex.RUnlock()
	return len(fake.dumpRequestArgsForCall)
}

func (fake *FakeHttpClientInterface) DumpRequestArgsForCall(i int) *http.Request {
	fake.dumpRequestMutex.RLock()
	defer fake.dumpRequestMutex.RUnlock()
	return fake.dumpRequestArgsForCall[i].arg1
}

func (fake *FakeHttpClientInterface) DumpResponse(arg1 *http.Response) {
	fake.dumpResponseMutex.Lock()
	fake.dumpResponseArgsForCall = append(fake.dumpResponseArgsForCall, struct {
		arg1 *http.Response
	}{arg1})
	fake.dumpResponseMutex.Unlock()
	if fake.DumpResponseStub != nil {
		fake.DumpResponseStub(arg1)
	}
}

func (fake *FakeHttpClientInterface) DumpResponseCallCount() int {
	fake.dumpResponseMutex.RLock()
	defer fake.dumpResponseMutex.RUnlock()
	return len(fake.dumpResponseArgsForCall)
}

func (fake *FakeHttpClientInterface) DumpResponseArgsForCall(i int) *http.Response {
	fake.dumpResponseMutex.RLock()
	defer fake.dumpResponseMutex.RUnlock()
	return fake.dumpResponseArgsForCall[i].arg1
}

func (fake *FakeHttpClientInterface) ExecuteCheckRedirect(req *http.Request, via []*http.Request) error {
	fake.executeCheckRedirectMutex.Lock()
	fake.executeCheckRedirectArgsForCall = append(fake.executeCheckRedirectArgsForCall, struct {
		req *http.Request
		via []*http.Request
	}{req, via})
	fake.executeCheckRedirectMutex.Unlock()
	if fake.ExecuteCheckRedirectStub != nil {
		return fake.ExecuteCheckRedirectStub(req, via)
	} else {
		return fake.executeCheckRedirectReturns.result1
	}
}

func (fake *FakeHttpClientInterface) ExecuteCheckRedirectCallCount() int {
	fake.executeCheckRedirectMutex.RLock()
	defer fake.executeCheckRedirectMutex.RUnlock()
	return len(fake.executeCheckRedirectArgsForCall)
}

func (fake *FakeHttpClientInterface) ExecuteCheckRedirectArgsForCall(i int) (*http.Request, []*http.Request) {
	fake.executeCheckRedirectMutex.RLock()
	defer fake.executeCheckRedirectMutex.RUnlock()
	return fake.executeCheckRedirectArgsForCall[i].req, fake.executeCheckRedirectArgsForCall[i].via
}

func (fake *FakeHttpClientInterface) ExecuteCheckRedirectReturns(result1 error) {
	fake.ExecuteCheckRedirectStub = nil
	fake.executeCheckRedirectReturns = struct {
		result1 error
	}{result1}
}

var _ net.HttpClientInterface = new(FakeHttpClientInterface)
