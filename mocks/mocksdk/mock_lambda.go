// Code generated by MockGen. DO NOT EDIT.
// Source: ./lambda.go

// Package mocksdk is a generated GoMock package.
package mocksdk

import (
	context "context"
	reflect "reflect"

	lambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	gomock "github.com/golang/mock/gomock"
)

// MockLambda is a mock of Lambda interface.
type MockLambda struct {
	ctrl     *gomock.Controller
	recorder *MockLambdaMockRecorder
}

// MockLambdaMockRecorder is the mock recorder for MockLambda.
type MockLambdaMockRecorder struct {
	mock *MockLambda
}

// NewMockLambda creates a new mock instance.
func NewMockLambda(ctrl *gomock.Controller) *MockLambda {
	mock := &MockLambda{ctrl: ctrl}
	mock.recorder = &MockLambdaMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLambda) EXPECT() *MockLambdaMockRecorder {
	return m.recorder
}

// Invoke mocks base method.
func (m *MockLambda) Invoke(ctx context.Context, params *lambda.InvokeInput, opts ...func(*lambda.Options)) (*lambda.InvokeOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Invoke", varargs...)
	ret0, _ := ret[0].(*lambda.InvokeOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Invoke indicates an expected call of Invoke.
func (mr *MockLambdaMockRecorder) Invoke(ctx, params interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Invoke", reflect.TypeOf((*MockLambda)(nil).Invoke), varargs...)
}
