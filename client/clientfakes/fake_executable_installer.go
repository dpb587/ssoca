// Code generated by counterfeiter. DO NOT EDIT.
package clientfakes

import (
	"sync"

	"github.com/dpb587/ssoca/client"
	"github.com/sirupsen/logrus"
)

type FakeExecutableInstaller struct {
	InstallStub        func(logrus.FieldLogger) error
	installMutex       sync.RWMutex
	installArgsForCall []struct {
		arg1 logrus.FieldLogger
	}
	installReturns struct {
		result1 error
	}
	installReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeExecutableInstaller) Install(arg1 logrus.FieldLogger) error {
	fake.installMutex.Lock()
	ret, specificReturn := fake.installReturnsOnCall[len(fake.installArgsForCall)]
	fake.installArgsForCall = append(fake.installArgsForCall, struct {
		arg1 logrus.FieldLogger
	}{arg1})
	fake.recordInvocation("Install", []interface{}{arg1})
	fake.installMutex.Unlock()
	if fake.InstallStub != nil {
		return fake.InstallStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.installReturns
	return fakeReturns.result1
}

func (fake *FakeExecutableInstaller) InstallCallCount() int {
	fake.installMutex.RLock()
	defer fake.installMutex.RUnlock()
	return len(fake.installArgsForCall)
}

func (fake *FakeExecutableInstaller) InstallCalls(stub func(logrus.FieldLogger) error) {
	fake.installMutex.Lock()
	defer fake.installMutex.Unlock()
	fake.InstallStub = stub
}

func (fake *FakeExecutableInstaller) InstallArgsForCall(i int) logrus.FieldLogger {
	fake.installMutex.RLock()
	defer fake.installMutex.RUnlock()
	argsForCall := fake.installArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeExecutableInstaller) InstallReturns(result1 error) {
	fake.installMutex.Lock()
	defer fake.installMutex.Unlock()
	fake.InstallStub = nil
	fake.installReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeExecutableInstaller) InstallReturnsOnCall(i int, result1 error) {
	fake.installMutex.Lock()
	defer fake.installMutex.Unlock()
	fake.InstallStub = nil
	if fake.installReturnsOnCall == nil {
		fake.installReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.installReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeExecutableInstaller) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.installMutex.RLock()
	defer fake.installMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeExecutableInstaller) recordInvocation(key string, args []interface{}) {
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

var _ client.ExecutableInstaller = new(FakeExecutableInstaller)
