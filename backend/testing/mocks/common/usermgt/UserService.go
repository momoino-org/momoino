// Code generated by mockery v2.46.2. DO NOT EDIT.

package mockusermgt

import (
	context "context"
	http "net/http"

	mock "github.com/stretchr/testify/mock"

	usermgt "wano-island/common/usermgt"
)

// MockUserService is an autogenerated mock type for the UserService type
type MockUserService struct {
	mock.Mock
}

type MockUserService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserService) EXPECT() *MockUserService_Expecter {
	return &MockUserService_Expecter{mock: &_m.Mock}
}

// ComparePassword provides a mock function with given fields: ctx, password, hasedPassword
func (_m *MockUserService) ComparePassword(ctx context.Context, password []byte, hasedPassword []byte) error {
	ret := _m.Called(ctx, password, hasedPassword)

	if len(ret) == 0 {
		panic("no return value specified for ComparePassword")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte, []byte) error); ok {
		r0 = rf(ctx, password, hasedPassword)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUserService_ComparePassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ComparePassword'
type MockUserService_ComparePassword_Call struct {
	*mock.Call
}

// ComparePassword is a helper method to define mock.On call
//   - ctx context.Context
//   - password []byte
//   - hasedPassword []byte
func (_e *MockUserService_Expecter) ComparePassword(ctx interface{}, password interface{}, hasedPassword interface{}) *MockUserService_ComparePassword_Call {
	return &MockUserService_ComparePassword_Call{Call: _e.mock.On("ComparePassword", ctx, password, hasedPassword)}
}

func (_c *MockUserService_ComparePassword_Call) Run(run func(ctx context.Context, password []byte, hasedPassword []byte)) *MockUserService_ComparePassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]byte), args[2].([]byte))
	})
	return _c
}

func (_c *MockUserService_ComparePassword_Call) Return(_a0 error) *MockUserService_ComparePassword_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUserService_ComparePassword_Call) RunAndReturn(run func(context.Context, []byte, []byte) error) *MockUserService_ComparePassword_Call {
	_c.Call.Return(run)
	return _c
}

// GenerateJWT provides a mock function with given fields: user
func (_m *MockUserService) GenerateJWT(user usermgt.UserModel) (*usermgt.JWT, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for GenerateJWT")
	}

	var r0 *usermgt.JWT
	var r1 error
	if rf, ok := ret.Get(0).(func(usermgt.UserModel) (*usermgt.JWT, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(usermgt.UserModel) *usermgt.JWT); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*usermgt.JWT)
		}
	}

	if rf, ok := ret.Get(1).(func(usermgt.UserModel) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserService_GenerateJWT_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateJWT'
type MockUserService_GenerateJWT_Call struct {
	*mock.Call
}

// GenerateJWT is a helper method to define mock.On call
//   - user usermgt.UserModel
func (_e *MockUserService_Expecter) GenerateJWT(user interface{}) *MockUserService_GenerateJWT_Call {
	return &MockUserService_GenerateJWT_Call{Call: _e.mock.On("GenerateJWT", user)}
}

func (_c *MockUserService_GenerateJWT_Call) Run(run func(user usermgt.UserModel)) *MockUserService_GenerateJWT_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(usermgt.UserModel))
	})
	return _c
}

func (_c *MockUserService_GenerateJWT_Call) Return(_a0 *usermgt.JWT, _a1 error) *MockUserService_GenerateJWT_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserService_GenerateJWT_Call) RunAndReturn(run func(usermgt.UserModel) (*usermgt.JWT, error)) *MockUserService_GenerateJWT_Call {
	_c.Call.Return(run)
	return _c
}

// HashPassword provides a mock function with given fields: ctx, password
func (_m *MockUserService) HashPassword(ctx context.Context, password string) ([]byte, error) {
	ret := _m.Called(ctx, password)

	if len(ret) == 0 {
		panic("no return value specified for HashPassword")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]byte, error)); ok {
		return rf(ctx, password)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []byte); ok {
		r0 = rf(ctx, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserService_HashPassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HashPassword'
type MockUserService_HashPassword_Call struct {
	*mock.Call
}

// HashPassword is a helper method to define mock.On call
//   - ctx context.Context
//   - password string
func (_e *MockUserService_Expecter) HashPassword(ctx interface{}, password interface{}) *MockUserService_HashPassword_Call {
	return &MockUserService_HashPassword_Call{Call: _e.mock.On("HashPassword", ctx, password)}
}

func (_c *MockUserService_HashPassword_Call) Run(run func(ctx context.Context, password string)) *MockUserService_HashPassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockUserService_HashPassword_Call) Return(_a0 []byte, _a1 error) *MockUserService_HashPassword_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserService_HashPassword_Call) RunAndReturn(run func(context.Context, string) ([]byte, error)) *MockUserService_HashPassword_Call {
	_c.Call.Return(run)
	return _c
}

// SetAuthCookies provides a mock function with given fields: w, jwt
func (_m *MockUserService) SetAuthCookies(w http.ResponseWriter, jwt usermgt.JWT) {
	_m.Called(w, jwt)
}

// MockUserService_SetAuthCookies_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetAuthCookies'
type MockUserService_SetAuthCookies_Call struct {
	*mock.Call
}

// SetAuthCookies is a helper method to define mock.On call
//   - w http.ResponseWriter
//   - jwt usermgt.JWT
func (_e *MockUserService_Expecter) SetAuthCookies(w interface{}, jwt interface{}) *MockUserService_SetAuthCookies_Call {
	return &MockUserService_SetAuthCookies_Call{Call: _e.mock.On("SetAuthCookies", w, jwt)}
}

func (_c *MockUserService_SetAuthCookies_Call) Run(run func(w http.ResponseWriter, jwt usermgt.JWT)) *MockUserService_SetAuthCookies_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(usermgt.JWT))
	})
	return _c
}

func (_c *MockUserService_SetAuthCookies_Call) Return() *MockUserService_SetAuthCookies_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockUserService_SetAuthCookies_Call) RunAndReturn(run func(http.ResponseWriter, usermgt.JWT)) *MockUserService_SetAuthCookies_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUserService creates a new instance of MockUserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserService {
	mock := &MockUserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}