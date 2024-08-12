// Code generated by mockery v2.44.1. DO NOT EDIT.

package mockcore

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// MockHTTPRoute is an autogenerated mock type for the HTTPRoute type
type MockHTTPRoute struct {
	mock.Mock
}

type MockHTTPRoute_Expecter struct {
	mock *mock.Mock
}

func (_m *MockHTTPRoute) EXPECT() *MockHTTPRoute_Expecter {
	return &MockHTTPRoute_Expecter{mock: &_m.Mock}
}

// IsPrivateRoute provides a mock function with given fields:
func (_m *MockHTTPRoute) IsPrivateRoute() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsPrivateRoute")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockHTTPRoute_IsPrivateRoute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsPrivateRoute'
type MockHTTPRoute_IsPrivateRoute_Call struct {
	*mock.Call
}

// IsPrivateRoute is a helper method to define mock.On call
func (_e *MockHTTPRoute_Expecter) IsPrivateRoute() *MockHTTPRoute_IsPrivateRoute_Call {
	return &MockHTTPRoute_IsPrivateRoute_Call{Call: _e.mock.On("IsPrivateRoute")}
}

func (_c *MockHTTPRoute_IsPrivateRoute_Call) Run(run func()) *MockHTTPRoute_IsPrivateRoute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockHTTPRoute_IsPrivateRoute_Call) Return(_a0 bool) *MockHTTPRoute_IsPrivateRoute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHTTPRoute_IsPrivateRoute_Call) RunAndReturn(run func() bool) *MockHTTPRoute_IsPrivateRoute_Call {
	_c.Call.Return(run)
	return _c
}

// Pattern provides a mock function with given fields:
func (_m *MockHTTPRoute) Pattern() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Pattern")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockHTTPRoute_Pattern_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Pattern'
type MockHTTPRoute_Pattern_Call struct {
	*mock.Call
}

// Pattern is a helper method to define mock.On call
func (_e *MockHTTPRoute_Expecter) Pattern() *MockHTTPRoute_Pattern_Call {
	return &MockHTTPRoute_Pattern_Call{Call: _e.mock.On("Pattern")}
}

func (_c *MockHTTPRoute_Pattern_Call) Run(run func()) *MockHTTPRoute_Pattern_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockHTTPRoute_Pattern_Call) Return(_a0 string) *MockHTTPRoute_Pattern_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHTTPRoute_Pattern_Call) RunAndReturn(run func() string) *MockHTTPRoute_Pattern_Call {
	_c.Call.Return(run)
	return _c
}

// ServeHTTP provides a mock function with given fields: _a0, _a1
func (_m *MockHTTPRoute) ServeHTTP(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// MockHTTPRoute_ServeHTTP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ServeHTTP'
type MockHTTPRoute_ServeHTTP_Call struct {
	*mock.Call
}

// ServeHTTP is a helper method to define mock.On call
//   - _a0 http.ResponseWriter
//   - _a1 *http.Request
func (_e *MockHTTPRoute_Expecter) ServeHTTP(_a0 interface{}, _a1 interface{}) *MockHTTPRoute_ServeHTTP_Call {
	return &MockHTTPRoute_ServeHTTP_Call{Call: _e.mock.On("ServeHTTP", _a0, _a1)}
}

func (_c *MockHTTPRoute_ServeHTTP_Call) Run(run func(_a0 http.ResponseWriter, _a1 *http.Request)) *MockHTTPRoute_ServeHTTP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(*http.Request))
	})
	return _c
}

func (_c *MockHTTPRoute_ServeHTTP_Call) Return() *MockHTTPRoute_ServeHTTP_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockHTTPRoute_ServeHTTP_Call) RunAndReturn(run func(http.ResponseWriter, *http.Request)) *MockHTTPRoute_ServeHTTP_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockHTTPRoute creates a new instance of MockHTTPRoute. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHTTPRoute(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockHTTPRoute {
	mock := &MockHTTPRoute{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
