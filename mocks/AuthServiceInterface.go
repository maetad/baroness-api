// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	jwt "github.com/golang-jwt/jwt/v4"
	authservice "github.com/maetad/baroness-api/internal/services/authservice"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// AuthServiceInterface is an autogenerated mock type for the AuthServiceInterface type
type AuthServiceInterface struct {
	mock.Mock
}

// GenerateToken provides a mock function with given fields: c, expiredIn
func (_m *AuthServiceInterface) GenerateToken(c authservice.Claimer, expiredIn time.Duration) (string, error) {
	ret := _m.Called(c, expiredIn)

	var r0 string
	if rf, ok := ret.Get(0).(func(authservice.Claimer, time.Duration) string); ok {
		r0 = rf(c, expiredIn)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(authservice.Claimer, time.Duration) error); ok {
		r1 = rf(c, expiredIn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseToken provides a mock function with given fields: tokenString
func (_m *AuthServiceInterface) ParseToken(tokenString string) (jwt.MapClaims, error) {
	ret := _m.Called(tokenString)

	var r0 jwt.MapClaims
	if rf, ok := ret.Get(0).(func(string) jwt.MapClaims); ok {
		r0 = rf(tokenString)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(jwt.MapClaims)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tokenString)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAuthServiceInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthServiceInterface creates a new instance of AuthServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthServiceInterface(t mockConstructorTestingTNewAuthServiceInterface) *AuthServiceInterface {
	mock := &AuthServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
