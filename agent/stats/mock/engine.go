// Copyright 2015-2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aws/amazon-ecs-agent/agent/stats (interfaces: Engine)

// Package mock_stats is a generated GoMock package.
package mock_stats

import (
	reflect "reflect"

	ecstcs "github.com/aws/amazon-ecs-agent/agent/tcs/model/ecstcs"
	gomock "github.com/golang/mock/gomock"
	"github.com/docker/docker/api/types"
)

// MockEngine is a mock of Engine interface
type MockEngine struct {
	ctrl     *gomock.Controller
	recorder *MockEngineMockRecorder
}

// MockEngineMockRecorder is the mock recorder for MockEngine
type MockEngineMockRecorder struct {
	mock *MockEngine
}

// NewMockEngine creates a new mock instance
func NewMockEngine(ctrl *gomock.Controller) *MockEngine {
	mock := &MockEngine{ctrl: ctrl}
	mock.recorder = &MockEngineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEngine) EXPECT() *MockEngineMockRecorder {
	return m.recorder
}

// ContainerDockerStats mocks base method
func (m *MockEngine) ContainerDockerStats(arg0, arg1 string) (*types.Stats, error) {
	ret := m.ctrl.Call(m, "ContainerDockerStats", arg0, arg1)
	ret0, _ := ret[0].(*types.Stats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainerDockerStats indicates an expected call of ContainerDockerStats
func (mr *MockEngineMockRecorder) ContainerDockerStats(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerDockerStats", reflect.TypeOf((*MockEngine)(nil).ContainerDockerStats), arg0, arg1)
}

// GetInstanceMetrics mocks base method
func (m *MockEngine) GetInstanceMetrics() (*ecstcs.MetricsMetadata, []*ecstcs.TaskMetric, error) {
	ret := m.ctrl.Call(m, "GetInstanceMetrics")
	ret0, _ := ret[0].(*ecstcs.MetricsMetadata)
	ret1, _ := ret[1].([]*ecstcs.TaskMetric)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetInstanceMetrics indicates an expected call of GetInstanceMetrics
func (mr *MockEngineMockRecorder) GetInstanceMetrics() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstanceMetrics", reflect.TypeOf((*MockEngine)(nil).GetInstanceMetrics))
}

// GetTaskHealthMetrics mocks base method
func (m *MockEngine) GetTaskHealthMetrics() (*ecstcs.HealthMetadata, []*ecstcs.TaskHealth, error) {
	ret := m.ctrl.Call(m, "GetTaskHealthMetrics")
	ret0, _ := ret[0].(*ecstcs.HealthMetadata)
	ret1, _ := ret[1].([]*ecstcs.TaskHealth)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetTaskHealthMetrics indicates an expected call of GetTaskHealthMetrics
func (mr *MockEngineMockRecorder) GetTaskHealthMetrics() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskHealthMetrics", reflect.TypeOf((*MockEngine)(nil).GetTaskHealthMetrics))
}
