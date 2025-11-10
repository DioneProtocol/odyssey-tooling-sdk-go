// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCloudService is a mock implementation of cloud services
type MockCloudService struct {
	mock.Mock
}

func (m *MockCloudService) CreateInstances(count int) ([]string, error) {
	args := m.Called(count)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCloudService) WaitForInstances(instanceIds []string) error {
	args := m.Called(instanceIds)
	return args.Error(0)
}

func (m *MockCloudService) GetInstancePublicIPs(instanceIds []string) (map[string]string, error) {
	args := m.Called(instanceIds)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockCloudService) CreateEIP(region string) (string, string, error) {
	args := m.Called(region)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockCloudService) AssociateEIP(instanceID, allocationID string) error {
	args := m.Called(instanceID, allocationID)
	return args.Error(0)
}

func (m *MockCloudService) SetupInstances(zone, network, sshKey, imageID, instanceType string, staticIPs []string, count, volumeSize int) ([]interface{}, error) {
	args := m.Called(zone, network, sshKey, imageID, instanceType, staticIPs, count, volumeSize)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockCloudService) SetPublicIP(zone, name string, count int) ([]string, error) {
	args := m.Called(zone, name, count)
	return args.Get(0).([]string), args.Error(1)
}

// MockNode is a mock implementation of Node for testing
type MockNode struct {
	mock.Mock
	Node
}

func (m *MockNode) WaitForSSHShell(timeout time.Duration) error {
	args := m.Called(timeout)
	return args.Error(0)
}

func (m *MockNode) Connect(port uint) error {
	args := m.Called(port)
	return args.Error(0)
}

func (m *MockNode) RunSSHSetupNode() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNode) RunSSHSetupDockerService() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNode) ComposeSSHSetupNode(hrp string, subnetIDs []string, version string, withMonitoring bool) error {
	args := m.Called(hrp, subnetIDs, version, withMonitoring)
	return args.Error(0)
}

func (m *MockNode) StartDockerCompose(timeout time.Duration) error {
	args := m.Called(timeout)
	return args.Error(0)
}

func (m *MockNode) ComposeSSHSetupLoadTest() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNode) RestartDockerCompose(timeout time.Duration) error {
	args := m.Called(timeout)
	return args.Error(0)
}

func (m *MockNode) RunSSHSetupMonitoringFolders() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNode) ComposeSSHSetupMonitoring() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNode) StartDockerComposeService(composeFile, service string, timeout time.Duration) error {
	args := m.Called(composeFile, service, timeout)
	return args.Error(0)
}

func TestNodeResults_Integration(t *testing.T) {
	// Test NodeResults functionality
	results := NodeResults{}

	// Test adding results
	results.AddResult("node-1", &Node{NodeID: "node-1"}, nil)
	results.AddResult("node-2", nil, assert.AnError)

	// Test getting error
	err := results.Error()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "node-2")
}

func TestCreateNodes_Concurrency(t *testing.T) {
	t.Skip("Cloud functionality has been removed from this SDK")
}

func TestCreateNodes_TimeoutHandling(t *testing.T) {
	t.Skip("Cloud functionality has been removed from this SDK")
}
