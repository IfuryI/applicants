// Code generated by mockery 2.9.4. DO NOT EDIT.

package delivery

import mock "github.com/stretchr/testify/mock"

// serviceMock is an autogenerated mock type for the Service type
type serviceMock struct {
	mock.Mock
}

// CheckCert provides a mock function with given fields: ogrn, kpp
func (_m *serviceMock) CheckCert(ogrn string, kpp string) (bool, error) {
	ret := _m.Called(ogrn, kpp)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(ogrn, kpp)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(ogrn, kpp)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
