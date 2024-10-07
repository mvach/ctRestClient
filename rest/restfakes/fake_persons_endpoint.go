// Code generated by counterfeiter. DO NOT EDIT.
package restfakes

import (
	"ctRestClient/rest"
	"encoding/json"
	"sync"
)

type FakePersonsEndpoint struct {
	GetPersonStub        func(int) ([]json.RawMessage, error)
	getPersonMutex       sync.RWMutex
	getPersonArgsForCall []struct {
		arg1 int
	}
	getPersonReturns struct {
		result1 []json.RawMessage
		result2 error
	}
	getPersonReturnsOnCall map[int]struct {
		result1 []json.RawMessage
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePersonsEndpoint) GetPerson(arg1 int) ([]json.RawMessage, error) {
	fake.getPersonMutex.Lock()
	ret, specificReturn := fake.getPersonReturnsOnCall[len(fake.getPersonArgsForCall)]
	fake.getPersonArgsForCall = append(fake.getPersonArgsForCall, struct {
		arg1 int
	}{arg1})
	stub := fake.GetPersonStub
	fakeReturns := fake.getPersonReturns
	fake.recordInvocation("GetPerson", []interface{}{arg1})
	fake.getPersonMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePersonsEndpoint) GetPersonCallCount() int {
	fake.getPersonMutex.RLock()
	defer fake.getPersonMutex.RUnlock()
	return len(fake.getPersonArgsForCall)
}

func (fake *FakePersonsEndpoint) GetPersonCalls(stub func(int) ([]json.RawMessage, error)) {
	fake.getPersonMutex.Lock()
	defer fake.getPersonMutex.Unlock()
	fake.GetPersonStub = stub
}

func (fake *FakePersonsEndpoint) GetPersonArgsForCall(i int) int {
	fake.getPersonMutex.RLock()
	defer fake.getPersonMutex.RUnlock()
	argsForCall := fake.getPersonArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakePersonsEndpoint) GetPersonReturns(result1 []json.RawMessage, result2 error) {
	fake.getPersonMutex.Lock()
	defer fake.getPersonMutex.Unlock()
	fake.GetPersonStub = nil
	fake.getPersonReturns = struct {
		result1 []json.RawMessage
		result2 error
	}{result1, result2}
}

func (fake *FakePersonsEndpoint) GetPersonReturnsOnCall(i int, result1 []json.RawMessage, result2 error) {
	fake.getPersonMutex.Lock()
	defer fake.getPersonMutex.Unlock()
	fake.GetPersonStub = nil
	if fake.getPersonReturnsOnCall == nil {
		fake.getPersonReturnsOnCall = make(map[int]struct {
			result1 []json.RawMessage
			result2 error
		})
	}
	fake.getPersonReturnsOnCall[i] = struct {
		result1 []json.RawMessage
		result2 error
	}{result1, result2}
}

func (fake *FakePersonsEndpoint) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getPersonMutex.RLock()
	defer fake.getPersonMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakePersonsEndpoint) recordInvocation(key string, args []interface{}) {
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

var _ rest.PersonsEndpoint = new(FakePersonsEndpoint)
