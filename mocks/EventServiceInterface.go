// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	eventservice "github.com/maetad/baroness-api/internal/services/eventservice"
	mock "github.com/stretchr/testify/mock"

	model "github.com/maetad/baroness-api/internal/model"
)

// EventServiceInterface is an autogenerated mock type for the EventServiceInterface type
type EventServiceInterface struct {
	mock.Mock
}

// Create provides a mock function with given fields: r, creator
func (_m *EventServiceInterface) Create(r eventservice.EventCreateRequest, creator *model.User) (*model.Event, error) {
	ret := _m.Called(r, creator)

	var r0 *model.Event
	if rf, ok := ret.Get(0).(func(eventservice.EventCreateRequest, *model.User) *model.Event); ok {
		r0 = rf(r, creator)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Event)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(eventservice.EventCreateRequest, *model.User) error); ok {
		r1 = rf(r, creator)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: event, deletor
func (_m *EventServiceInterface) Delete(event *model.Event, deletor *model.User) error {
	ret := _m.Called(event, deletor)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Event, *model.User) error); ok {
		r0 = rf(event, deletor)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: id
func (_m *EventServiceInterface) Get(id uint) (*model.Event, error) {
	ret := _m.Called(id)

	var r0 *model.Event
	if rf, ok := ret.Get(0).(func(uint) *model.Event); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Event)
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

// List provides a mock function with given fields:
func (_m *EventServiceInterface) List() ([]model.Event, error) {
	ret := _m.Called()

	var r0 []model.Event
	if rf, ok := ret.Get(0).(func() []model.Event); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Event)
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

// Update provides a mock function with given fields: event, r, updator
func (_m *EventServiceInterface) Update(event *model.Event, r eventservice.EventUpdateRequest, updator *model.User) (*model.Event, error) {
	ret := _m.Called(event, r, updator)

	var r0 *model.Event
	if rf, ok := ret.Get(0).(func(*model.Event, eventservice.EventUpdateRequest, *model.User) *model.Event); ok {
		r0 = rf(event, r, updator)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Event)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Event, eventservice.EventUpdateRequest, *model.User) error); ok {
		r1 = rf(event, r, updator)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewEventServiceInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewEventServiceInterface creates a new instance of EventServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEventServiceInterface(t mockConstructorTestingTNewEventServiceInterface) *EventServiceInterface {
	mock := &EventServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
