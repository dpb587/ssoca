// Code generated by counterfeiter. DO NOT EDIT.
package clientfakes

import (
	"io"
	"sync"

	"github.com/cloudfoundry/bosh-cli/ui"
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/config"
	"github.com/dpb587/ssoca/httpclient"
	"github.com/dpb587/ssoca/version"
	"github.com/sirupsen/logrus"
)

type FakeRuntime struct {
	GetExecStub        func() string
	getExecMutex       sync.RWMutex
	getExecArgsForCall []struct{}
	getExecReturns     struct {
		result1 string
	}
	getExecReturnsOnCall map[int]struct {
		result1 string
	}
	GetVersionStub        func() version.Version
	getVersionMutex       sync.RWMutex
	getVersionArgsForCall []struct{}
	getVersionReturns     struct {
		result1 version.Version
	}
	getVersionReturnsOnCall map[int]struct {
		result1 version.Version
	}
	GetEnvironmentStub        func() (config.EnvironmentState, error)
	getEnvironmentMutex       sync.RWMutex
	getEnvironmentArgsForCall []struct{}
	getEnvironmentReturns     struct {
		result1 config.EnvironmentState
		result2 error
	}
	getEnvironmentReturnsOnCall map[int]struct {
		result1 config.EnvironmentState
		result2 error
	}
	GetEnvironmentNameStub        func() string
	getEnvironmentNameMutex       sync.RWMutex
	getEnvironmentNameArgsForCall []struct{}
	getEnvironmentNameReturns     struct {
		result1 string
	}
	getEnvironmentNameReturnsOnCall map[int]struct {
		result1 string
	}
	GetConfigManagerStub        func() (config.Manager, error)
	getConfigManagerMutex       sync.RWMutex
	getConfigManagerArgsForCall []struct{}
	getConfigManagerReturns     struct {
		result1 config.Manager
		result2 error
	}
	getConfigManagerReturnsOnCall map[int]struct {
		result1 config.Manager
		result2 error
	}
	GetClientStub        func() (httpclient.Client, error)
	getClientMutex       sync.RWMutex
	getClientArgsForCall []struct{}
	getClientReturns     struct {
		result1 httpclient.Client
		result2 error
	}
	getClientReturnsOnCall map[int]struct {
		result1 httpclient.Client
		result2 error
	}
	GetAuthInterceptClientStub        func() (httpclient.Client, error)
	getAuthInterceptClientMutex       sync.RWMutex
	getAuthInterceptClientArgsForCall []struct{}
	getAuthInterceptClientReturns     struct {
		result1 httpclient.Client
		result2 error
	}
	getAuthInterceptClientReturnsOnCall map[int]struct {
		result1 httpclient.Client
		result2 error
	}
	GetUIStub        func() ui.UI
	getUIMutex       sync.RWMutex
	getUIArgsForCall []struct{}
	getUIReturns     struct {
		result1 ui.UI
	}
	getUIReturnsOnCall map[int]struct {
		result1 ui.UI
	}
	GetLoggerStub        func() logrus.FieldLogger
	getLoggerMutex       sync.RWMutex
	getLoggerArgsForCall []struct{}
	getLoggerReturns     struct {
		result1 logrus.FieldLogger
	}
	getLoggerReturnsOnCall map[int]struct {
		result1 logrus.FieldLogger
	}
	GetStderrStub        func() io.Writer
	getStderrMutex       sync.RWMutex
	getStderrArgsForCall []struct{}
	getStderrReturns     struct {
		result1 io.Writer
	}
	getStderrReturnsOnCall map[int]struct {
		result1 io.Writer
	}
	GetStdoutStub        func() io.Writer
	getStdoutMutex       sync.RWMutex
	getStdoutArgsForCall []struct{}
	getStdoutReturns     struct {
		result1 io.Writer
	}
	getStdoutReturnsOnCall map[int]struct {
		result1 io.Writer
	}
	GetStdinStub        func() io.Reader
	getStdinMutex       sync.RWMutex
	getStdinArgsForCall []struct{}
	getStdinReturns     struct {
		result1 io.Reader
	}
	getStdinReturnsOnCall map[int]struct {
		result1 io.Reader
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRuntime) GetExec() string {
	fake.getExecMutex.Lock()
	ret, specificReturn := fake.getExecReturnsOnCall[len(fake.getExecArgsForCall)]
	fake.getExecArgsForCall = append(fake.getExecArgsForCall, struct{}{})
	fake.recordInvocation("GetExec", []interface{}{})
	fake.getExecMutex.Unlock()
	if fake.GetExecStub != nil {
		return fake.GetExecStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getExecReturns.result1
}

func (fake *FakeRuntime) GetExecCallCount() int {
	fake.getExecMutex.RLock()
	defer fake.getExecMutex.RUnlock()
	return len(fake.getExecArgsForCall)
}

func (fake *FakeRuntime) GetExecReturns(result1 string) {
	fake.GetExecStub = nil
	fake.getExecReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeRuntime) GetExecReturnsOnCall(i int, result1 string) {
	fake.GetExecStub = nil
	if fake.getExecReturnsOnCall == nil {
		fake.getExecReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getExecReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeRuntime) GetVersion() version.Version {
	fake.getVersionMutex.Lock()
	ret, specificReturn := fake.getVersionReturnsOnCall[len(fake.getVersionArgsForCall)]
	fake.getVersionArgsForCall = append(fake.getVersionArgsForCall, struct{}{})
	fake.recordInvocation("GetVersion", []interface{}{})
	fake.getVersionMutex.Unlock()
	if fake.GetVersionStub != nil {
		return fake.GetVersionStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getVersionReturns.result1
}

func (fake *FakeRuntime) GetVersionCallCount() int {
	fake.getVersionMutex.RLock()
	defer fake.getVersionMutex.RUnlock()
	return len(fake.getVersionArgsForCall)
}

func (fake *FakeRuntime) GetVersionReturns(result1 version.Version) {
	fake.GetVersionStub = nil
	fake.getVersionReturns = struct {
		result1 version.Version
	}{result1}
}

func (fake *FakeRuntime) GetVersionReturnsOnCall(i int, result1 version.Version) {
	fake.GetVersionStub = nil
	if fake.getVersionReturnsOnCall == nil {
		fake.getVersionReturnsOnCall = make(map[int]struct {
			result1 version.Version
		})
	}
	fake.getVersionReturnsOnCall[i] = struct {
		result1 version.Version
	}{result1}
}

func (fake *FakeRuntime) GetEnvironment() (config.EnvironmentState, error) {
	fake.getEnvironmentMutex.Lock()
	ret, specificReturn := fake.getEnvironmentReturnsOnCall[len(fake.getEnvironmentArgsForCall)]
	fake.getEnvironmentArgsForCall = append(fake.getEnvironmentArgsForCall, struct{}{})
	fake.recordInvocation("GetEnvironment", []interface{}{})
	fake.getEnvironmentMutex.Unlock()
	if fake.GetEnvironmentStub != nil {
		return fake.GetEnvironmentStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getEnvironmentReturns.result1, fake.getEnvironmentReturns.result2
}

func (fake *FakeRuntime) GetEnvironmentCallCount() int {
	fake.getEnvironmentMutex.RLock()
	defer fake.getEnvironmentMutex.RUnlock()
	return len(fake.getEnvironmentArgsForCall)
}

func (fake *FakeRuntime) GetEnvironmentReturns(result1 config.EnvironmentState, result2 error) {
	fake.GetEnvironmentStub = nil
	fake.getEnvironmentReturns = struct {
		result1 config.EnvironmentState
		result2 error
	}{result1, result2}
}

func (fake *FakeRuntime) GetEnvironmentReturnsOnCall(i int, result1 config.EnvironmentState, result2 error) {
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

func (fake *FakeRuntime) GetEnvironmentName() string {
	fake.getEnvironmentNameMutex.Lock()
	ret, specificReturn := fake.getEnvironmentNameReturnsOnCall[len(fake.getEnvironmentNameArgsForCall)]
	fake.getEnvironmentNameArgsForCall = append(fake.getEnvironmentNameArgsForCall, struct{}{})
	fake.recordInvocation("GetEnvironmentName", []interface{}{})
	fake.getEnvironmentNameMutex.Unlock()
	if fake.GetEnvironmentNameStub != nil {
		return fake.GetEnvironmentNameStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getEnvironmentNameReturns.result1
}

func (fake *FakeRuntime) GetEnvironmentNameCallCount() int {
	fake.getEnvironmentNameMutex.RLock()
	defer fake.getEnvironmentNameMutex.RUnlock()
	return len(fake.getEnvironmentNameArgsForCall)
}

func (fake *FakeRuntime) GetEnvironmentNameReturns(result1 string) {
	fake.GetEnvironmentNameStub = nil
	fake.getEnvironmentNameReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeRuntime) GetEnvironmentNameReturnsOnCall(i int, result1 string) {
	fake.GetEnvironmentNameStub = nil
	if fake.getEnvironmentNameReturnsOnCall == nil {
		fake.getEnvironmentNameReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getEnvironmentNameReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeRuntime) GetConfigManager() (config.Manager, error) {
	fake.getConfigManagerMutex.Lock()
	ret, specificReturn := fake.getConfigManagerReturnsOnCall[len(fake.getConfigManagerArgsForCall)]
	fake.getConfigManagerArgsForCall = append(fake.getConfigManagerArgsForCall, struct{}{})
	fake.recordInvocation("GetConfigManager", []interface{}{})
	fake.getConfigManagerMutex.Unlock()
	if fake.GetConfigManagerStub != nil {
		return fake.GetConfigManagerStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getConfigManagerReturns.result1, fake.getConfigManagerReturns.result2
}

func (fake *FakeRuntime) GetConfigManagerCallCount() int {
	fake.getConfigManagerMutex.RLock()
	defer fake.getConfigManagerMutex.RUnlock()
	return len(fake.getConfigManagerArgsForCall)
}

func (fake *FakeRuntime) GetConfigManagerReturns(result1 config.Manager, result2 error) {
	fake.GetConfigManagerStub = nil
	fake.getConfigManagerReturns = struct {
		result1 config.Manager
		result2 error
	}{result1, result2}
}

func (fake *FakeRuntime) GetConfigManagerReturnsOnCall(i int, result1 config.Manager, result2 error) {
	fake.GetConfigManagerStub = nil
	if fake.getConfigManagerReturnsOnCall == nil {
		fake.getConfigManagerReturnsOnCall = make(map[int]struct {
			result1 config.Manager
			result2 error
		})
	}
	fake.getConfigManagerReturnsOnCall[i] = struct {
		result1 config.Manager
		result2 error
	}{result1, result2}
}

func (fake *FakeRuntime) GetClient() (httpclient.Client, error) {
	fake.getClientMutex.Lock()
	ret, specificReturn := fake.getClientReturnsOnCall[len(fake.getClientArgsForCall)]
	fake.getClientArgsForCall = append(fake.getClientArgsForCall, struct{}{})
	fake.recordInvocation("GetClient", []interface{}{})
	fake.getClientMutex.Unlock()
	if fake.GetClientStub != nil {
		return fake.GetClientStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getClientReturns.result1, fake.getClientReturns.result2
}

func (fake *FakeRuntime) GetClientCallCount() int {
	fake.getClientMutex.RLock()
	defer fake.getClientMutex.RUnlock()
	return len(fake.getClientArgsForCall)
}

func (fake *FakeRuntime) GetClientReturns(result1 httpclient.Client, result2 error) {
	fake.GetClientStub = nil
	fake.getClientReturns = struct {
		result1 httpclient.Client
		result2 error
	}{result1, result2}
}

func (fake *FakeRuntime) GetClientReturnsOnCall(i int, result1 httpclient.Client, result2 error) {
	fake.GetClientStub = nil
	if fake.getClientReturnsOnCall == nil {
		fake.getClientReturnsOnCall = make(map[int]struct {
			result1 httpclient.Client
			result2 error
		})
	}
	fake.getClientReturnsOnCall[i] = struct {
		result1 httpclient.Client
		result2 error
	}{result1, result2}
}

func (fake *FakeRuntime) GetAuthInterceptClient() (httpclient.Client, error) {
	fake.getAuthInterceptClientMutex.Lock()
	ret, specificReturn := fake.getAuthInterceptClientReturnsOnCall[len(fake.getAuthInterceptClientArgsForCall)]
	fake.getAuthInterceptClientArgsForCall = append(fake.getAuthInterceptClientArgsForCall, struct{}{})
	fake.recordInvocation("GetAuthInterceptClient", []interface{}{})
	fake.getAuthInterceptClientMutex.Unlock()
	if fake.GetAuthInterceptClientStub != nil {
		return fake.GetAuthInterceptClientStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getAuthInterceptClientReturns.result1, fake.getAuthInterceptClientReturns.result2
}

func (fake *FakeRuntime) GetAuthInterceptClientCallCount() int {
	fake.getAuthInterceptClientMutex.RLock()
	defer fake.getAuthInterceptClientMutex.RUnlock()
	return len(fake.getAuthInterceptClientArgsForCall)
}

func (fake *FakeRuntime) GetAuthInterceptClientReturns(result1 httpclient.Client, result2 error) {
	fake.GetAuthInterceptClientStub = nil
	fake.getAuthInterceptClientReturns = struct {
		result1 httpclient.Client
		result2 error
	}{result1, result2}
}

func (fake *FakeRuntime) GetAuthInterceptClientReturnsOnCall(i int, result1 httpclient.Client, result2 error) {
	fake.GetAuthInterceptClientStub = nil
	if fake.getAuthInterceptClientReturnsOnCall == nil {
		fake.getAuthInterceptClientReturnsOnCall = make(map[int]struct {
			result1 httpclient.Client
			result2 error
		})
	}
	fake.getAuthInterceptClientReturnsOnCall[i] = struct {
		result1 httpclient.Client
		result2 error
	}{result1, result2}
}

func (fake *FakeRuntime) GetUI() ui.UI {
	fake.getUIMutex.Lock()
	ret, specificReturn := fake.getUIReturnsOnCall[len(fake.getUIArgsForCall)]
	fake.getUIArgsForCall = append(fake.getUIArgsForCall, struct{}{})
	fake.recordInvocation("GetUI", []interface{}{})
	fake.getUIMutex.Unlock()
	if fake.GetUIStub != nil {
		return fake.GetUIStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getUIReturns.result1
}

func (fake *FakeRuntime) GetUICallCount() int {
	fake.getUIMutex.RLock()
	defer fake.getUIMutex.RUnlock()
	return len(fake.getUIArgsForCall)
}

func (fake *FakeRuntime) GetUIReturns(result1 ui.UI) {
	fake.GetUIStub = nil
	fake.getUIReturns = struct {
		result1 ui.UI
	}{result1}
}

func (fake *FakeRuntime) GetUIReturnsOnCall(i int, result1 ui.UI) {
	fake.GetUIStub = nil
	if fake.getUIReturnsOnCall == nil {
		fake.getUIReturnsOnCall = make(map[int]struct {
			result1 ui.UI
		})
	}
	fake.getUIReturnsOnCall[i] = struct {
		result1 ui.UI
	}{result1}
}

func (fake *FakeRuntime) GetLogger() logrus.FieldLogger {
	fake.getLoggerMutex.Lock()
	ret, specificReturn := fake.getLoggerReturnsOnCall[len(fake.getLoggerArgsForCall)]
	fake.getLoggerArgsForCall = append(fake.getLoggerArgsForCall, struct{}{})
	fake.recordInvocation("GetLogger", []interface{}{})
	fake.getLoggerMutex.Unlock()
	if fake.GetLoggerStub != nil {
		return fake.GetLoggerStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getLoggerReturns.result1
}

func (fake *FakeRuntime) GetLoggerCallCount() int {
	fake.getLoggerMutex.RLock()
	defer fake.getLoggerMutex.RUnlock()
	return len(fake.getLoggerArgsForCall)
}

func (fake *FakeRuntime) GetLoggerReturns(result1 logrus.FieldLogger) {
	fake.GetLoggerStub = nil
	fake.getLoggerReturns = struct {
		result1 logrus.FieldLogger
	}{result1}
}

func (fake *FakeRuntime) GetLoggerReturnsOnCall(i int, result1 logrus.FieldLogger) {
	fake.GetLoggerStub = nil
	if fake.getLoggerReturnsOnCall == nil {
		fake.getLoggerReturnsOnCall = make(map[int]struct {
			result1 logrus.FieldLogger
		})
	}
	fake.getLoggerReturnsOnCall[i] = struct {
		result1 logrus.FieldLogger
	}{result1}
}

func (fake *FakeRuntime) GetStderr() io.Writer {
	fake.getStderrMutex.Lock()
	ret, specificReturn := fake.getStderrReturnsOnCall[len(fake.getStderrArgsForCall)]
	fake.getStderrArgsForCall = append(fake.getStderrArgsForCall, struct{}{})
	fake.recordInvocation("GetStderr", []interface{}{})
	fake.getStderrMutex.Unlock()
	if fake.GetStderrStub != nil {
		return fake.GetStderrStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getStderrReturns.result1
}

func (fake *FakeRuntime) GetStderrCallCount() int {
	fake.getStderrMutex.RLock()
	defer fake.getStderrMutex.RUnlock()
	return len(fake.getStderrArgsForCall)
}

func (fake *FakeRuntime) GetStderrReturns(result1 io.Writer) {
	fake.GetStderrStub = nil
	fake.getStderrReturns = struct {
		result1 io.Writer
	}{result1}
}

func (fake *FakeRuntime) GetStderrReturnsOnCall(i int, result1 io.Writer) {
	fake.GetStderrStub = nil
	if fake.getStderrReturnsOnCall == nil {
		fake.getStderrReturnsOnCall = make(map[int]struct {
			result1 io.Writer
		})
	}
	fake.getStderrReturnsOnCall[i] = struct {
		result1 io.Writer
	}{result1}
}

func (fake *FakeRuntime) GetStdout() io.Writer {
	fake.getStdoutMutex.Lock()
	ret, specificReturn := fake.getStdoutReturnsOnCall[len(fake.getStdoutArgsForCall)]
	fake.getStdoutArgsForCall = append(fake.getStdoutArgsForCall, struct{}{})
	fake.recordInvocation("GetStdout", []interface{}{})
	fake.getStdoutMutex.Unlock()
	if fake.GetStdoutStub != nil {
		return fake.GetStdoutStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getStdoutReturns.result1
}

func (fake *FakeRuntime) GetStdoutCallCount() int {
	fake.getStdoutMutex.RLock()
	defer fake.getStdoutMutex.RUnlock()
	return len(fake.getStdoutArgsForCall)
}

func (fake *FakeRuntime) GetStdoutReturns(result1 io.Writer) {
	fake.GetStdoutStub = nil
	fake.getStdoutReturns = struct {
		result1 io.Writer
	}{result1}
}

func (fake *FakeRuntime) GetStdoutReturnsOnCall(i int, result1 io.Writer) {
	fake.GetStdoutStub = nil
	if fake.getStdoutReturnsOnCall == nil {
		fake.getStdoutReturnsOnCall = make(map[int]struct {
			result1 io.Writer
		})
	}
	fake.getStdoutReturnsOnCall[i] = struct {
		result1 io.Writer
	}{result1}
}

func (fake *FakeRuntime) GetStdin() io.Reader {
	fake.getStdinMutex.Lock()
	ret, specificReturn := fake.getStdinReturnsOnCall[len(fake.getStdinArgsForCall)]
	fake.getStdinArgsForCall = append(fake.getStdinArgsForCall, struct{}{})
	fake.recordInvocation("GetStdin", []interface{}{})
	fake.getStdinMutex.Unlock()
	if fake.GetStdinStub != nil {
		return fake.GetStdinStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getStdinReturns.result1
}

func (fake *FakeRuntime) GetStdinCallCount() int {
	fake.getStdinMutex.RLock()
	defer fake.getStdinMutex.RUnlock()
	return len(fake.getStdinArgsForCall)
}

func (fake *FakeRuntime) GetStdinReturns(result1 io.Reader) {
	fake.GetStdinStub = nil
	fake.getStdinReturns = struct {
		result1 io.Reader
	}{result1}
}

func (fake *FakeRuntime) GetStdinReturnsOnCall(i int, result1 io.Reader) {
	fake.GetStdinStub = nil
	if fake.getStdinReturnsOnCall == nil {
		fake.getStdinReturnsOnCall = make(map[int]struct {
			result1 io.Reader
		})
	}
	fake.getStdinReturnsOnCall[i] = struct {
		result1 io.Reader
	}{result1}
}

func (fake *FakeRuntime) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getExecMutex.RLock()
	defer fake.getExecMutex.RUnlock()
	fake.getVersionMutex.RLock()
	defer fake.getVersionMutex.RUnlock()
	fake.getEnvironmentMutex.RLock()
	defer fake.getEnvironmentMutex.RUnlock()
	fake.getEnvironmentNameMutex.RLock()
	defer fake.getEnvironmentNameMutex.RUnlock()
	fake.getConfigManagerMutex.RLock()
	defer fake.getConfigManagerMutex.RUnlock()
	fake.getClientMutex.RLock()
	defer fake.getClientMutex.RUnlock()
	fake.getAuthInterceptClientMutex.RLock()
	defer fake.getAuthInterceptClientMutex.RUnlock()
	fake.getUIMutex.RLock()
	defer fake.getUIMutex.RUnlock()
	fake.getLoggerMutex.RLock()
	defer fake.getLoggerMutex.RUnlock()
	fake.getStderrMutex.RLock()
	defer fake.getStderrMutex.RUnlock()
	fake.getStdoutMutex.RLock()
	defer fake.getStdoutMutex.RUnlock()
	fake.getStdinMutex.RLock()
	defer fake.getStdinMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeRuntime) recordInvocation(key string, args []interface{}) {
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

var _ client.Runtime = new(FakeRuntime)
