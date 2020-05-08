// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import protoutil "github.com/wingbaas/platformsrv/settings/fabric/txgeneratev2/protoutil"

// PolicyEvaluator is an autogenerated mock type for the PolicyEvaluator type
type PolicyEvaluator struct {
	mock.Mock
}

// Evaluate provides a mock function with given fields: policyBytes, signatureSet
func (_m *PolicyEvaluator) Evaluate(policyBytes []byte, signatureSet []*protoutil.SignedData) error {
	ret := _m.Called(policyBytes, signatureSet)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, []*protoutil.SignedData) error); ok {
		r0 = rf(policyBytes, signatureSet)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
