// Code generated by MockGen. DO NOT EDIT.
// Source: ./assume_clients.go

// Package mocksdk is a generated GoMock package.
package mocksdk

import (
	reflect "reflect"

	sdk "github.com/Tanemahuta/aws-lambda-server/pkg/aws/sdk"
	arn "github.com/aws/aws-sdk-go-v2/aws/arn"
	gomock "github.com/golang/mock/gomock"
)

// MockAssumeClients is a mock of AssumeClients interface.
type MockAssumeClients[C sdk.Client] struct {
	ctrl     *gomock.Controller
	recorder *MockAssumeClientsMockRecorder[C]
}

// MockAssumeClientsMockRecorder is the mock recorder for MockAssumeClients.
type MockAssumeClientsMockRecorder[C sdk.Client] struct {
	mock *MockAssumeClients[C]
}

// NewMockAssumeClients creates a new mock instance.
func NewMockAssumeClients[C sdk.Client](ctrl *gomock.Controller) *MockAssumeClients[C] {
	mock := &MockAssumeClients[C]{ctrl: ctrl}
	mock.recorder = &MockAssumeClientsMockRecorder[C]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAssumeClients[C]) EXPECT() *MockAssumeClientsMockRecorder[C] {
	return m.recorder
}

// Get mocks base method.
func (m *MockAssumeClients[C]) Get(assumeRole *arn.ARN) C {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", assumeRole)
	ret0, _ := ret[0].(C)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockAssumeClientsMockRecorder[C]) Get(assumeRole interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAssumeClients[C])(nil).Get), assumeRole)
}