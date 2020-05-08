// Code generated by counterfeiter. DO NOT EDIT.
package mock

import (
	"sync"

	"github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/common/chaincode"
)

type MetadataUpdateListener struct {
	HandleMetadataUpdateStub        func(string, chaincode.MetadataSet)
	handleMetadataUpdateMutex       sync.RWMutex
	handleMetadataUpdateArgsForCall []struct {
		arg1 string
		arg2 chaincode.MetadataSet
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *MetadataUpdateListener) HandleMetadataUpdate(arg1 string, arg2 chaincode.MetadataSet) {
	fake.handleMetadataUpdateMutex.Lock()
	fake.handleMetadataUpdateArgsForCall = append(fake.handleMetadataUpdateArgsForCall, struct {
		arg1 string
		arg2 chaincode.MetadataSet
	}{arg1, arg2})
	fake.recordInvocation("HandleMetadataUpdate", []interface{}{arg1, arg2})
	fake.handleMetadataUpdateMutex.Unlock()
	if fake.HandleMetadataUpdateStub != nil {
		fake.HandleMetadataUpdateStub(arg1, arg2)
	}
}

func (fake *MetadataUpdateListener) HandleMetadataUpdateCallCount() int {
	fake.handleMetadataUpdateMutex.RLock()
	defer fake.handleMetadataUpdateMutex.RUnlock()
	return len(fake.handleMetadataUpdateArgsForCall)
}

func (fake *MetadataUpdateListener) HandleMetadataUpdateCalls(stub func(string, chaincode.MetadataSet)) {
	fake.handleMetadataUpdateMutex.Lock()
	defer fake.handleMetadataUpdateMutex.Unlock()
	fake.HandleMetadataUpdateStub = stub
}

func (fake *MetadataUpdateListener) HandleMetadataUpdateArgsForCall(i int) (string, chaincode.MetadataSet) {
	fake.handleMetadataUpdateMutex.RLock()
	defer fake.handleMetadataUpdateMutex.RUnlock()
	argsForCall := fake.handleMetadataUpdateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *MetadataUpdateListener) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.handleMetadataUpdateMutex.RLock()
	defer fake.handleMetadataUpdateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *MetadataUpdateListener) recordInvocation(key string, args []interface{}) {
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
