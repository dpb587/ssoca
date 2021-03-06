// Code generated by counterfeiter. DO NOT EDIT.
package servicefakes

import (
	"sync"

	"github.com/dpb587/ssoca/server/service"
	servicea "github.com/dpb587/ssoca/service"
)

type FakeServiceFactory struct {
	CreateStub        func(string, map[string]interface{}) (service.Service, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		arg1 string
		arg2 map[string]interface{}
	}
	createReturns struct {
		result1 service.Service
		result2 error
	}
	createReturnsOnCall map[int]struct {
		result1 service.Service
		result2 error
	}
	TypeStub        func() servicea.Type
	typeMutex       sync.RWMutex
	typeArgsForCall []struct {
	}
	typeReturns struct {
		result1 servicea.Type
	}
	typeReturnsOnCall map[int]struct {
		result1 servicea.Type
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeServiceFactory) Create(arg1 string, arg2 map[string]interface{}) (service.Service, error) {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		arg1 string
		arg2 map[string]interface{}
	}{arg1, arg2})
	fake.recordInvocation("Create", []interface{}{arg1, arg2})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.createReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeServiceFactory) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeServiceFactory) CreateCalls(stub func(string, map[string]interface{}) (service.Service, error)) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = stub
}

func (fake *FakeServiceFactory) CreateArgsForCall(i int) (string, map[string]interface{}) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	argsForCall := fake.createArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeServiceFactory) CreateReturns(result1 service.Service, result2 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 service.Service
		result2 error
	}{result1, result2}
}

func (fake *FakeServiceFactory) CreateReturnsOnCall(i int, result1 service.Service, result2 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	if fake.createReturnsOnCall == nil {
		fake.createReturnsOnCall = make(map[int]struct {
			result1 service.Service
			result2 error
		})
	}
	fake.createReturnsOnCall[i] = struct {
		result1 service.Service
		result2 error
	}{result1, result2}
}

func (fake *FakeServiceFactory) Type() servicea.Type {
	fake.typeMutex.Lock()
	ret, specificReturn := fake.typeReturnsOnCall[len(fake.typeArgsForCall)]
	fake.typeArgsForCall = append(fake.typeArgsForCall, struct {
	}{})
	fake.recordInvocation("Type", []interface{}{})
	fake.typeMutex.Unlock()
	if fake.TypeStub != nil {
		return fake.TypeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.typeReturns
	return fakeReturns.result1
}

func (fake *FakeServiceFactory) TypeCallCount() int {
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	return len(fake.typeArgsForCall)
}

func (fake *FakeServiceFactory) TypeCalls(stub func() servicea.Type) {
	fake.typeMutex.Lock()
	defer fake.typeMutex.Unlock()
	fake.TypeStub = stub
}

func (fake *FakeServiceFactory) TypeReturns(result1 servicea.Type) {
	fake.typeMutex.Lock()
	defer fake.typeMutex.Unlock()
	fake.TypeStub = nil
	fake.typeReturns = struct {
		result1 servicea.Type
	}{result1}
}

func (fake *FakeServiceFactory) TypeReturnsOnCall(i int, result1 servicea.Type) {
	fake.typeMutex.Lock()
	defer fake.typeMutex.Unlock()
	fake.TypeStub = nil
	if fake.typeReturnsOnCall == nil {
		fake.typeReturnsOnCall = make(map[int]struct {
			result1 servicea.Type
		})
	}
	fake.typeReturnsOnCall[i] = struct {
		result1 servicea.Type
	}{result1}
}

func (fake *FakeServiceFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeServiceFactory) recordInvocation(key string, args []interface{}) {
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

var _ service.ServiceFactory = new(FakeServiceFactory)
