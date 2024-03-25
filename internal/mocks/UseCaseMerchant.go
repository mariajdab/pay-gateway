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

// ValidateMerchant provides a mock function with given fields: merchantCode
func (_m *UseCaseMerchant) ValidateMerchant(merchantCode string) (string, error) {
	ret := _m.Called(merchantCode)

	if len(ret) == 0 {
		panic("no return value specified for ValidateMerchant")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(merchantCode)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(merchantCode)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(merchantCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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