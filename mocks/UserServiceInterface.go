// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	userservice "github.com/pakkaparn/no-idea-api/internal/services/userservice"
	mock "github.com/stretchr/testify/mock"
)

// UserServiceInterface is an autogenerated mock type for the UserServiceInterface type
type UserServiceInterface struct {
	mock.Mock
}

// Create provides a mock function with given fields: r
func (_m *UserServiceInterface) Create(r userservice.UserCreateRequest) (userservice.UserInterface, error) {
	ret := _m.Called(r)

	var r0 userservice.UserInterface
	if rf, ok := ret.Get(0).(func(userservice.UserCreateRequest) userservice.UserInterface); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(userservice.UserInterface)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(userservice.UserCreateRequest) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: id
func (_m *UserServiceInterface) Get(id uint) (userservice.UserInterface, error) {
	ret := _m.Called(id)

	var r0 userservice.UserInterface
	if rf, ok := ret.Get(0).(func(uint) userservice.UserInterface); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(userservice.UserInterface)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUsername provides a mock function with given fields: username
func (_m *UserServiceInterface) GetByUsername(username string) (userservice.UserInterface, error) {
	ret := _m.Called(username)

	var r0 userservice.UserInterface
	if rf, ok := ret.Get(0).(func(string) userservice.UserInterface); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(userservice.UserInterface)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields:
func (_m *UserServiceInterface) List() ([]userservice.UserInterface, error) {
	ret := _m.Called()

	var r0 []userservice.UserInterface
	if rf, ok := ret.Get(0).(func() []userservice.UserInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]userservice.UserInterface)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: user, r
func (_m *UserServiceInterface) Update(user userservice.UserInterface, r userservice.UserUpdateRequest) (userservice.UserInterface, error) {
	ret := _m.Called(user, r)

	var r0 userservice.UserInterface
	if rf, ok := ret.Get(0).(func(userservice.UserInterface, userservice.UserUpdateRequest) userservice.UserInterface); ok {
		r0 = rf(user, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(userservice.UserInterface)
		}
	}

	var r1 error
			if rf, ok := ret.Get(1).(func(userservice.UserInterface, userservice.UserUpdateRequest) error); ok {
		r1 = rf(user, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserServiceInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserServiceInterface creates a new instance of UserServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserServiceInterface(t mockConstructorTestingTNewUserServiceInterface) *UserServiceInterface {
	mock := &UserServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
