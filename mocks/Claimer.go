// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Claimer is an autogenerated mock type for the Claimer type
type Claimer struct {
	mock.Mock
}

// GetClaims provides a mock function with given fields:
func (_m *Claimer) GetClaims() map[string]interface{} {
	ret := _m.Called()

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func() map[string]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

type mockConstructorTestingTNewClaimer interface {
	mock.TestingT
	Cleanup(func())
}

// NewClaimer creates a new instance of Claimer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClaimer(t mockConstructorTestingTNewClaimer) *Claimer {
	mock := &Claimer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
