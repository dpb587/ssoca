// Code generated by counterfeiter. DO NOT EDIT.
package httpclientfakes

import (
	"sync"

	"github.com/dpb587/ssoca/service/env/api"
	"github.com/dpb587/ssoca/service/env/httpclient"
)

type FakeClient struct {
	GetAuthStub        func() (api.AuthResponse, error)
	getAuthMutex       sync.RWMutex
	getAuthArgsForCall []struct {
	}
	getAuthReturns struct {
		result1 api.AuthResponse
		result2 error
	}
	getAuthReturnsOnCall map[int]struct {
		result1 api.AuthResponse
		result2 error
	}
	GetInfoStub        func() (api.InfoResponse, error)
	getInfoMutex       sync.RWMutex
	getInfoArgsForCall []struct {
	}
	getInfoReturns struct {
		result1 api.InfoResponse
		result2 error
	}
	getInfoReturnsOnCall map[int]struct {
		result1 api.InfoResponse
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) GetAuth() (api.AuthResponse, error) {
	fake.getAuthMutex.Lock()
	ret, specificReturn := fake.getAuthReturnsOnCall[len(fake.getAuthArgsForCall)]
	fake.getAuthArgsForCall = append(fake.getAuthArgsForCall, struct {
	}{})
	fake.recordInvocation("GetAuth", []interface{}{})
	fake.getAuthMutex.Unlock()
	if fake.GetAuthStub != nil {
		return fake.GetAuthStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getAuthReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetAuthCallCount() int {
	fake.getAuthMutex.RLock()
	defer fake.getAuthMutex.RUnlock()
	return len(fake.getAuthArgsForCall)
}

func (fake *FakeClient) GetAuthCalls(stub func() (api.AuthResponse, error)) {
	fake.getAuthMutex.Lock()
	defer fake.getAuthMutex.Unlock()
	fake.GetAuthStub = stub
}

func (fake *FakeClient) GetAuthReturns(result1 api.AuthResponse, result2 error) {
	fake.getAuthMutex.Lock()
	defer fake.getAuthMutex.Unlock()
	fake.GetAuthStub = nil
	fake.getAuthReturns = struct {
		result1 api.AuthResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetAuthReturnsOnCall(i int, result1 api.AuthResponse, result2 error) {
	fake.getAuthMutex.Lock()
	defer fake.getAuthMutex.Unlock()
	fake.GetAuthStub = nil
	if fake.getAuthReturnsOnCall == nil {
		fake.getAuthReturnsOnCall = make(map[int]struct {
			result1 api.AuthResponse
			result2 error
		})
	}
	fake.getAuthReturnsOnCall[i] = struct {
		result1 api.AuthResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetInfo() (api.InfoResponse, error) {
	fake.getInfoMutex.Lock()
	ret, specificReturn := fake.getInfoReturnsOnCall[len(fake.getInfoArgsForCall)]
	fake.getInfoArgsForCall = append(fake.getInfoArgsForCall, struct {
	}{})
	fake.recordInvocation("GetInfo", []interface{}{})
	fake.getInfoMutex.Unlock()
	if fake.GetInfoStub != nil {
		return fake.GetInfoStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getInfoReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetInfoCallCount() int {
	fake.getInfoMutex.RLock()
	defer fake.getInfoMutex.RUnlock()
	return len(fake.getInfoArgsForCall)
}

func (fake *FakeClient) GetInfoCalls(stub func() (api.InfoResponse, error)) {
	fake.getInfoMutex.Lock()
	defer fake.getInfoMutex.Unlock()
	fake.GetInfoStub = stub
}

func (fake *FakeClient) GetInfoReturns(result1 api.InfoResponse, result2 error) {
	fake.getInfoMutex.Lock()
	defer fake.getInfoMutex.Unlock()
	fake.GetInfoStub = nil
	fake.getInfoReturns = struct {
		result1 api.InfoResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetInfoReturnsOnCall(i int, result1 api.InfoResponse, result2 error) {
	fake.getInfoMutex.Lock()
	defer fake.getInfoMutex.Unlock()
	fake.GetInfoStub = nil
	if fake.getInfoReturnsOnCall == nil {
		fake.getInfoReturnsOnCall = make(map[int]struct {
			result1 api.InfoResponse
			result2 error
		})
	}
	fake.getInfoReturnsOnCall[i] = struct {
		result1 api.InfoResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getAuthMutex.RLock()
	defer fake.getAuthMutex.RUnlock()
	fake.getInfoMutex.RLock()
	defer fake.getInfoMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ httpclient.Client = new(FakeClient)
