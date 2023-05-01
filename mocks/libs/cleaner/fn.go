// Code generated by mockery v2.20.0. DO NOT EDIT.

package cleaner

import mock "github.com/stretchr/testify/mock"

// Fn is an autogenerated mock type for the Fn type
type Fn struct {
	mock.Mock
}

type Fn_Expecter struct {
	mock *mock.Mock
}

func (_m *Fn) EXPECT() *Fn_Expecter {
	return &Fn_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields:
func (_m *Fn) Execute() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Fn_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type Fn_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
func (_e *Fn_Expecter) Execute() *Fn_Execute_Call {
	return &Fn_Execute_Call{Call: _e.mock.On("Execute")}
}

func (_c *Fn_Execute_Call) Run(run func()) *Fn_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Fn_Execute_Call) Return(_a0 error) *Fn_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Fn_Execute_Call) RunAndReturn(run func() error) *Fn_Execute_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewFn interface {
	mock.TestingT
	Cleanup(func())
}

// NewFn creates a new instance of Fn. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFn(t mockConstructorTestingTNewFn) *Fn {
	mock := &Fn{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
