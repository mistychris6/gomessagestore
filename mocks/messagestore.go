// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/blackhatbrigade/gomessagestore (interfaces: MessageStore)

// Package mock_gomessagestore is a generated GoMock package.
package mock_gomessagestore

import (
	context "context"
	gomessagestore "github.com/blackhatbrigade/gomessagestore"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockMessageStore is a mock of MessageStore interface
type MockMessageStore struct {
	ctrl     *gomock.Controller
	recorder *MockMessageStoreMockRecorder
}

// MockMessageStoreMockRecorder is the mock recorder for MockMessageStore
type MockMessageStoreMockRecorder struct {
	mock *MockMessageStore
}

// NewMockMessageStore creates a new mock instance
func NewMockMessageStore(ctrl *gomock.Controller) *MockMessageStore {
	mock := &MockMessageStore{ctrl: ctrl}
	mock.recorder = &MockMessageStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMessageStore) EXPECT() *MockMessageStoreMockRecorder {
	return m.recorder
}

// CreateProjector mocks base method
func (m *MockMessageStore) CreateProjector(arg0 ...gomessagestore.ProjectorOption) (gomessagestore.Projector, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateProjector", varargs...)
	ret0, _ := ret[0].(gomessagestore.Projector)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProjector indicates an expected call of CreateProjector
func (mr *MockMessageStoreMockRecorder) CreateProjector(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProjector", reflect.TypeOf((*MockMessageStore)(nil).CreateProjector), arg0...)
}

// CreateSubscriber mocks base method
func (m *MockMessageStore) CreateSubscriber(arg0 string, arg1 []gomessagestore.MessageHandler, arg2 ...gomessagestore.SubscriberOption) (gomessagestore.Subscriber, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateSubscriber", varargs...)
	ret0, _ := ret[0].(gomessagestore.Subscriber)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSubscriber indicates an expected call of CreateSubscriber
func (mr *MockMessageStoreMockRecorder) CreateSubscriber(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSubscriber", reflect.TypeOf((*MockMessageStore)(nil).CreateSubscriber), varargs...)
}

// Get mocks base method
func (m *MockMessageStore) Get(arg0 context.Context, arg1 ...gomessagestore.GetOption) ([]gomessagestore.Message, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].([]gomessagestore.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockMessageStoreMockRecorder) Get(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMessageStore)(nil).Get), varargs...)
}

// Write mocks base method
func (m *MockMessageStore) Write(arg0 context.Context, arg1 gomessagestore.Message, arg2 ...gomessagestore.WriteOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Write", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write
func (mr *MockMessageStoreMockRecorder) Write(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockMessageStore)(nil).Write), varargs...)
}
