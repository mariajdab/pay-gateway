// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	entity "github.com/mariajdab/pay-gateway/internal/entity"

	mock "github.com/stretchr/testify/mock"
)

// UseCaseMerchant is an autogenerated mock type for the UseCaseMerchant type
type UseCaseMerchant struct {
	mock.Mock
}

// CreateMerchant provides a mock function with given fields: _a0
func (_m *UseCaseMerchant) CreateMerchant(_a0 entity.Merchant) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for CreateMerchant")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(entity.Merchant) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUseCaseMerchant creates a new instance of UseCaseMerchant. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUseCaseMerchant(t interface {
	mock.TestingT
	Cleanup(func())
}) *UseCaseMerchant {
	mock := &UseCaseMerchant{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}