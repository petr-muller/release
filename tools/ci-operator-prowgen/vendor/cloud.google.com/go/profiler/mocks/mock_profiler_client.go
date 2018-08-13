// Copyright 2017 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by MockGen. DO NOT EDIT.
// Source: google.golang.org/genproto/googleapis/devtools/cloudprofiler/v2 (interfaces: ProfilerServiceClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	context "golang.org/x/net/context"
	v2 "google.golang.org/genproto/googleapis/devtools/cloudprofiler/v2"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockProfilerServiceClient is a mock of ProfilerServiceClient interface
type MockProfilerServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockProfilerServiceClientMockRecorder
}

// MockProfilerServiceClientMockRecorder is the mock recorder for MockProfilerServiceClient
type MockProfilerServiceClientMockRecorder struct {
	mock *MockProfilerServiceClient
}

// NewMockProfilerServiceClient creates a new mock instance
func NewMockProfilerServiceClient(ctrl *gomock.Controller) *MockProfilerServiceClient {
	mock := &MockProfilerServiceClient{ctrl: ctrl}
	mock.recorder = &MockProfilerServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProfilerServiceClient) EXPECT() *MockProfilerServiceClientMockRecorder {
	return m.recorder
}

// CreateOfflineProfile mocks base method
func (m *MockProfilerServiceClient) CreateOfflineProfile(arg0 context.Context, arg1 *v2.CreateOfflineProfileRequest, arg2 ...grpc.CallOption) (*v2.Profile, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateOfflineProfile", varargs...)
	ret0, _ := ret[0].(*v2.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOfflineProfile indicates an expected call of CreateOfflineProfile
func (mr *MockProfilerServiceClientMockRecorder) CreateOfflineProfile(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOfflineProfile", reflect.TypeOf((*MockProfilerServiceClient)(nil).CreateOfflineProfile), varargs...)
}

// CreateProfile mocks base method
func (m *MockProfilerServiceClient) CreateProfile(arg0 context.Context, arg1 *v2.CreateProfileRequest, arg2 ...grpc.CallOption) (*v2.Profile, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateProfile", varargs...)
	ret0, _ := ret[0].(*v2.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProfile indicates an expected call of CreateProfile
func (mr *MockProfilerServiceClientMockRecorder) CreateProfile(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProfile", reflect.TypeOf((*MockProfilerServiceClient)(nil).CreateProfile), varargs...)
}

// UpdateProfile mocks base method
func (m *MockProfilerServiceClient) UpdateProfile(arg0 context.Context, arg1 *v2.UpdateProfileRequest, arg2 ...grpc.CallOption) (*v2.Profile, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateProfile", varargs...)
	ret0, _ := ret[0].(*v2.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProfile indicates an expected call of UpdateProfile
func (mr *MockProfilerServiceClientMockRecorder) UpdateProfile(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfile", reflect.TypeOf((*MockProfilerServiceClient)(nil).UpdateProfile), varargs...)
}
