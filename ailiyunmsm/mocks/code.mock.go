// Code generated by MockGen. DO NOT EDIT.
// Source: ./ailiyunmsm/msm.go
//
// Generated by this command:
//
//	mockgen.exe -source=./ailiyunmsm/msm.go -package=codemock -destination=./ailiyunmsm/mocks/code.mock.go MSMRetry
//
// Package codemock is a generated GoMock package.
package codemock

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCode is a mock of Code interface.
type MockCode struct {
	ctrl     *gomock.Controller
	recorder *MockCodeMockRecorder
}

// MockCodeMockRecorder is the mock recorder for MockCode.
type MockCodeMockRecorder struct {
	mock *MockCode
}

// NewMockCode creates a new mock instance.
func NewMockCode(ctrl *gomock.Controller) *MockCode {
	mock := &MockCode{ctrl: ctrl}
	mock.recorder = &MockCodeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCode) EXPECT() *MockCodeMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockCode) Send(singerName, code string, phoneNumber ...string) error {
	m.ctrl.T.Helper()
	varargs := []any{singerName, code}
	for _, a := range phoneNumber {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Send", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockCodeMockRecorder) Send(singerName, code any, phoneNumber ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{singerName, code}, phoneNumber...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockCode)(nil).Send), varargs...)
}
