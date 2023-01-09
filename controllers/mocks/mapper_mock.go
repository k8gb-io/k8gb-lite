// Code generated by MockGen. DO NOT EDIT.
// Source: controllers/mapper/mapper.go

// Package mocks is a generated GoMock package.
package mocks

/*
Copyright 2022 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	reflect "reflect"

	mapper "cloud.example.com/annotation-operator/controllers/mapper"
	gomock "github.com/golang/mock/gomock"
)

// MockMapper is a mock of Mapper interface.
type MockMapper struct {
	ctrl     *gomock.Controller
	recorder *MockMapperMockRecorder
}

// MockMapperMockRecorder is the mock recorder for MockMapper.
type MockMapperMockRecorder struct {
	mock *MockMapper
}

// NewMockMapper creates a new mock instance.
func NewMockMapper(ctrl *gomock.Controller) *MockMapper {
	mock := &MockMapper{ctrl: ctrl}
	mock.recorder = &MockMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMapper) EXPECT() *MockMapperMockRecorder {
	return m.recorder
}

// Equal mocks base method.
func (m *MockMapper) Equal(arg0 *mapper.LoopState) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Equal", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Equal indicates an expected call of Equal.
func (mr *MockMapperMockRecorder) Equal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Equal", reflect.TypeOf((*MockMapper)(nil).Equal), arg0)
}

// GetExposedIPs mocks base method.
func (m *MockMapper) GetExposedIPs() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExposedIPs")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExposedIPs indicates an expected call of GetExposedIPs.
func (mr *MockMapperMockRecorder) GetExposedIPs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExposedIPs", reflect.TypeOf((*MockMapper)(nil).GetExposedIPs))
}

// GetStatus mocks base method.
func (m *MockMapper) GetStatus() mapper.Status {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus")
	ret0, _ := ret[0].(mapper.Status)
	return ret0
}

// GetStatus indicates an expected call of GetStatus.
func (mr *MockMapperMockRecorder) GetStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockMapper)(nil).GetStatus))
}

// SetReference mocks base method.
func (m *MockMapper) SetReference(arg0 *mapper.LoopState) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetReference", arg0)
}

// SetReference indicates an expected call of SetReference.
func (mr *MockMapperMockRecorder) SetReference(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReference", reflect.TypeOf((*MockMapper)(nil).SetReference), arg0)
}

// TryInjectFinalizer mocks base method.
func (m *MockMapper) TryInjectFinalizer() (mapper.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TryInjectFinalizer")
	ret0, _ := ret[0].(mapper.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TryInjectFinalizer indicates an expected call of TryInjectFinalizer.
func (mr *MockMapperMockRecorder) TryInjectFinalizer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryInjectFinalizer", reflect.TypeOf((*MockMapper)(nil).TryInjectFinalizer))
}

// TryRemoveDNSEndpoint mocks base method.
func (m *MockMapper) TryRemoveDNSEndpoint() (mapper.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TryRemoveDNSEndpoint")
	ret0, _ := ret[0].(mapper.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TryRemoveDNSEndpoint indicates an expected call of TryRemoveDNSEndpoint.
func (mr *MockMapperMockRecorder) TryRemoveDNSEndpoint() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryRemoveDNSEndpoint", reflect.TypeOf((*MockMapper)(nil).TryRemoveDNSEndpoint))
}

// TryRemoveFinalizer mocks base method.
func (m *MockMapper) TryRemoveFinalizer(arg0 func(*mapper.LoopState) error) (mapper.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TryRemoveFinalizer", arg0)
	ret0, _ := ret[0].(mapper.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TryRemoveFinalizer indicates an expected call of TryRemoveFinalizer.
func (mr *MockMapperMockRecorder) TryRemoveFinalizer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryRemoveFinalizer", reflect.TypeOf((*MockMapper)(nil).TryRemoveFinalizer), arg0)
}

// UpdateStatusAnnotation mocks base method.
func (m *MockMapper) UpdateStatusAnnotation() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatusAnnotation")
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatusAnnotation indicates an expected call of UpdateStatusAnnotation.
func (mr *MockMapperMockRecorder) UpdateStatusAnnotation() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatusAnnotation", reflect.TypeOf((*MockMapper)(nil).UpdateStatusAnnotation))
}
