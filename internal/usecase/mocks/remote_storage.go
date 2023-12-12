// Code generated by mockery v3.0.0-alpha.0. DO NOT EDIT.

package mocks

import (
	remote_storage "github.com/horockey/distributed_kv_storage/internal/adapter/gateway/remote_storage"
	mock "github.com/stretchr/testify/mock"
)

// RemoteStorage is an autogenerated mock type for the Gateway type
type RemoteStorage struct {
	mock.Mock
}

// Get provides a mock function with given fields: key, node
func (_m *RemoteStorage) Get(key string, node remote_storage.AppNode) (map[string]interface{}, error) {
	ret := _m.Called(key, node)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(string, remote_storage.AppNode) map[string]interface{}); ok {
		r0 = rf(key, node)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, remote_storage.AppNode) error); ok {
		r1 = rf(key, node)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Set provides a mock function with given fields: key, val, node
func (_m *RemoteStorage) Set(key string, val map[string]interface{}, node remote_storage.AppNode) error {
	ret := _m.Called(key, val, node)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, map[string]interface{}, remote_storage.AppNode) error); ok {
		r0 = rf(key, val, node)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRemoteStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewRemoteStorage creates a new instance of RemoteStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRemoteStorage(t mockConstructorTestingTNewRemoteStorage) *RemoteStorage {
	mock := &RemoteStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}