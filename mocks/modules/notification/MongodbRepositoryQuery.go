// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MongodbRepositoryQuery is an autogenerated mock type for the MongodbRepositoryQuery type
type MongodbRepositoryQuery struct {
	mock.Mock
}

// NewMongodbRepositoryQuery creates a new instance of MongodbRepositoryQuery. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMongodbRepositoryQuery(t interface {
	mock.TestingT
	Cleanup(func())
}) *MongodbRepositoryQuery {
	mock := &MongodbRepositoryQuery{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
