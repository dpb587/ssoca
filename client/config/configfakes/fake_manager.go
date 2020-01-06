// Code generated by counterfeiter. DO NOT EDIT.
package configfakes

import (
	"sync"

	"github.com/dpb587/ssoca/client/config"
)

type FakeManager struct {
	GetEnvironmentStub        func(string) (config.EnvironmentState, error)
	getEnvironmentMutex       sync.RWMutex
	getEnvironmentArgsForCall []struct {
		arg1 string
	}
	getEnvironmentReturns struct {
		result1 config.EnvironmentState
		result2 error
	}
	getEnvironmentReturnsOnCall map[int]struct {
		result1 config.EnvironmentState
		result2 error
	}
	GetEnvironmentsStub        func() (config.EnvironmentsState, error)
	getEnvironmentsMutex       sync.RWMutex
	getEnvironmentsArgsForCall []struct {
	}
	getEnvironmentsReturns struct {
		result1 config.EnvironmentsState
		result2 error
	}
	getEnvironmentsReturnsOnCall map[int]struct {
		result1 config.EnvironmentsState
		result2 error
	}
	GetSourceStub        func() string
	getSourceMutex       sync.RWMutex
	getSourceArgsForCall []struct {
	}
	getSourceReturns struct {
		result1 string
	}
	getSourceReturnsOnCall map[int]struct {
		result1 string
	}
	SetEnvironmentStub        func(config.EnvironmentState) error
	setEnvironmentMutex       sync.RWMutex
	setEnvironmentArgsForCall []struct {
		arg1 config.EnvironmentState
	}
	setEnvironmentReturns struct {
		result1 error
	}
	setEnvironmentReturnsOnCall map[int]struct {
		result1 error
	}
	UnsetEnvironmentStub        func(string) error
	unsetEnvironmentMutex       sync.RWMutex
	unsetEnvironmentArgsForCall []struct {
		arg1 string
	}
	unsetEnvironmentReturns struct {
		result1 error
	}
	unsetEnvironmentReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeManager) GetEnvironment(arg1 string) (config.EnvironmentState, error) {
	fake.getEnvironmentMutex.Lock()
	ret, specificReturn := fake.getEnvironmentReturnsOnCall[len(fake.getEnvironmentArgsForCall)]
	fake.getEnvironmentArgsForCall = append(fake.getEnvironmentArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetEnvironment", []interface{}{arg1})
	fake.getEnvironmentMutex.Unlock()
	if fake.GetEnvironmentStub != nil {
		return fake.GetEnvironmentStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getEnvironmentReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) GetEnvironmentCallCount() int {
	fake.getEnvironmentMutex.RLock()
	defer fake.getEnvironmentMutex.RUnlock()
	return len(fake.getEnvironmentArgsForCall)
}

func (fake *FakeManager) GetEnvironmentCalls(stub func(string) (config.EnvironmentState, error)) {
	fake.getEnvironmentMutex.Lock()
	defer fake.getEnvironmentMutex.Unlock()
	fake.GetEnvironmentStub = stub
}

func (fake *FakeManager) GetEnvironmentArgsForCall(i int) string {
	fake.getEnvironmentMutex.RLock()
	defer fake.getEnvironmentMutex.RUnlock()
	argsForCall := fake.getEnvironmentArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeManager) GetEnvironmentReturns(result1 config.EnvironmentState, result2 error) {
	fake.getEnvironmentMutex.Lock()
	defer fake.getEnvironmentMutex.Unlock()
	fake.GetEnvironmentStub = nil
	fake.getEnvironmentReturns = struct {
		result1 config.EnvironmentState
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) GetEnvironmentReturnsOnCall(i int, result1 config.EnvironmentState, result2 error) {
	fake.getEnvironmentMutex.Lock()
	defer fake.getEnvironmentMutex.Unlock()
	fake.GetEnvironmentStub = nil
	if fake.getEnvironmentReturnsOnCall == nil {
		fake.getEnvironmentReturnsOnCall = make(map[int]struct {
			result1 config.EnvironmentState
			result2 error
		})
	}
	fake.getEnvironmentReturnsOnCall[i] = struct {
		result1 config.EnvironmentState
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) GetEnvironments() (config.EnvironmentsState, error) {
	fake.getEnvironmentsMutex.Lock()
	ret, specificReturn := fake.getEnvironmentsReturnsOnCall[len(fake.getEnvironmentsArgsForCall)]
	fake.getEnvironmentsArgsForCall = append(fake.getEnvironmentsArgsForCall, struct {
	}{})
	fake.recordInvocation("GetEnvironments", []interface{}{})
	fake.getEnvironmentsMutex.Unlock()
	if fake.GetEnvironmentsStub != nil {
		return fake.GetEnvironmentsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getEnvironmentsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) GetEnvironmentsCallCount() int {
	fake.getEnvironmentsMutex.RLock()
	defer fake.getEnvironmentsMutex.RUnlock()
	return len(fake.getEnvironmentsArgsForCall)
}

func (fake *FakeManager) GetEnvironmentsCalls(stub func() (config.EnvironmentsState, error)) {
	fake.getEnvironmentsMutex.Lock()
	defer fake.getEnvironmentsMutex.Unlock()
	fake.GetEnvironmentsStub = stub
}

func (fake *FakeManager) GetEnvironmentsReturns(result1 config.EnvironmentsState, result2 error) {
	fake.getEnvironmentsMutex.Lock()
	defer fake.getEnvironmentsMutex.Unlock()
	fake.GetEnvironmentsStub = nil
	fake.getEnvironmentsReturns = struct {
		result1 config.EnvironmentsState
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) GetEnvironmentsReturnsOnCall(i int, result1 config.EnvironmentsState, result2 error) {
	fake.getEnvironmentsMutex.Lock()
	defer fake.getEnvironmentsMutex.Unlock()
	fake.GetEnvironmentsStub = nil
	if fake.getEnvironmentsReturnsOnCall == nil {
		fake.getEnvironmentsReturnsOnCall = make(map[int]struct {
			result1 config.EnvironmentsState
			result2 error
		})
	}
	fake.getEnvironmentsReturnsOnCall[i] = struct {
		result1 config.EnvironmentsState
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) GetSource() string {
	fake.getSourceMutex.Lock()
	ret, specificReturn := fake.getSourceReturnsOnCall[len(fake.getSourceArgsForCall)]
	fake.getSourceArgsForCall = append(fake.getSourceArgsForCall, struct {
	}{})
	fake.recordInvocation("GetSource", []interface{}{})
	fake.getSourceMutex.Unlock()
	if fake.GetSourceStub != nil {
		return fake.GetSourceStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getSourceReturns
	return fakeReturns.result1
}

func (fake *FakeManager) GetSourceCallCount() int {
	fake.getSourceMutex.RLock()
	defer fake.getSourceMutex.RUnlock()
	return len(fake.getSourceArgsForCall)
}

func (fake *FakeManager) GetSourceCalls(stub func() string) {
	fake.getSourceMutex.Lock()
	defer fake.getSourceMutex.Unlock()
	fake.GetSourceStub = stub
}

func (fake *FakeManager) GetSourceReturns(result1 string) {
	fake.getSourceMutex.Lock()
	defer fake.getSourceMutex.Unlock()
	fake.GetSourceStub = nil
	fake.getSourceReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeManager) GetSourceReturnsOnCall(i int, result1 string) {
	fake.getSourceMutex.Lock()
	defer fake.getSourceMutex.Unlock()
	fake.GetSourceStub = nil
	if fake.getSourceReturnsOnCall == nil {
		fake.getSourceReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getSourceReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeManager) SetEnvironment(arg1 config.EnvironmentState) error {
	fake.setEnvironmentMutex.Lock()
	ret, specificReturn := fake.setEnvironmentReturnsOnCall[len(fake.setEnvironmentArgsForCall)]
	fake.setEnvironmentArgsForCall = append(fake.setEnvironmentArgsForCall, struct {
		arg1 config.EnvironmentState
	}{arg1})
	fake.recordInvocation("SetEnvironment", []interface{}{arg1})
	fake.setEnvironmentMutex.Unlock()
	if fake.SetEnvironmentStub != nil {
		return fake.SetEnvironmentStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.setEnvironmentReturns
	return fakeReturns.result1
}

func (fake *FakeManager) SetEnvironmentCallCount() int {
	fake.setEnvironmentMutex.RLock()
	defer fake.setEnvironmentMutex.RUnlock()
	return len(fake.setEnvironmentArgsForCall)
}

func (fake *FakeManager) SetEnvironmentCalls(stub func(config.EnvironmentState) error) {
	fake.setEnvironmentMutex.Lock()
	defer fake.setEnvironmentMutex.Unlock()
	fake.SetEnvironmentStub = stub
}

func (fake *FakeManager) SetEnvironmentArgsForCall(i int) config.EnvironmentState {
	fake.setEnvironmentMutex.RLock()
	defer fake.setEnvironmentMutex.RUnlock()
	argsForCall := fake.setEnvironmentArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeManager) SetEnvironmentReturns(result1 error) {
	fake.setEnvironmentMutex.Lock()
	defer fake.setEnvironmentMutex.Unlock()
	fake.SetEnvironmentStub = nil
	fake.setEnvironmentReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) SetEnvironmentReturnsOnCall(i int, result1 error) {
	fake.setEnvironmentMutex.Lock()
	defer fake.setEnvironmentMutex.Unlock()
	fake.SetEnvironmentStub = nil
	if fake.setEnvironmentReturnsOnCall == nil {
		fake.setEnvironmentReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setEnvironmentReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) UnsetEnvironment(arg1 string) error {
	fake.unsetEnvironmentMutex.Lock()
	ret, specificReturn := fake.unsetEnvironmentReturnsOnCall[len(fake.unsetEnvironmentArgsForCall)]
	fake.unsetEnvironmentArgsForCall = append(fake.unsetEnvironmentArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("UnsetEnvironment", []interface{}{arg1})
	fake.unsetEnvironmentMutex.Unlock()
	if fake.UnsetEnvironmentStub != nil {
		return fake.UnsetEnvironmentStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.unsetEnvironmentReturns
	return fakeReturns.result1
}

func (fake *FakeManager) UnsetEnvironmentCallCount() int {
	fake.unsetEnvironmentMutex.RLock()
	defer fake.unsetEnvironmentMutex.RUnlock()
	return len(fake.unsetEnvironmentArgsForCall)
}

func (fake *FakeManager) UnsetEnvironmentCalls(stub func(string) error) {
	fake.unsetEnvironmentMutex.Lock()
	defer fake.unsetEnvironmentMutex.Unlock()
	fake.UnsetEnvironmentStub = stub
}

func (fake *FakeManager) UnsetEnvironmentArgsForCall(i int) string {
	fake.unsetEnvironmentMutex.RLock()
	defer fake.unsetEnvironmentMutex.RUnlock()
	argsForCall := fake.unsetEnvironmentArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeManager) UnsetEnvironmentReturns(result1 error) {
	fake.unsetEnvironmentMutex.Lock()
	defer fake.unsetEnvironmentMutex.Unlock()
	fake.UnsetEnvironmentStub = nil
	fake.unsetEnvironmentReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) UnsetEnvironmentReturnsOnCall(i int, result1 error) {
	fake.unsetEnvironmentMutex.Lock()
	defer fake.unsetEnvironmentMutex.Unlock()
	fake.UnsetEnvironmentStub = nil
	if fake.unsetEnvironmentReturnsOnCall == nil {
		fake.unsetEnvironmentReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.unsetEnvironmentReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getEnvironmentMutex.RLock()
	defer fake.getEnvironmentMutex.RUnlock()
	fake.getEnvironmentsMutex.RLock()
	defer fake.getEnvironmentsMutex.RUnlock()
	fake.getSourceMutex.RLock()
	defer fake.getSourceMutex.RUnlock()
	fake.setEnvironmentMutex.RLock()
	defer fake.setEnvironmentMutex.RUnlock()
	fake.unsetEnvironmentMutex.RLock()
	defer fake.unsetEnvironmentMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeManager) recordInvocation(key string, args []interface{}) {
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

var _ config.Manager = new(FakeManager)
