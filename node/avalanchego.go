// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	remoteconfig "github.com/DioneProtocol/odyssey-tooling-sdk-go/node/config"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/utils"
	"github.com/DioneProtocol/odysseygo/api/info"
)

func (h *Node) GetOdysseyGoVersion() (string, error) {
	// Craft and send the HTTP POST request
	requestBody := "{\"jsonrpc\":\"2.0\", \"id\":1,\"method\" :\"info.getNodeVersion\"}"
	resp, err := h.Post("", requestBody)
	if err != nil {
		return "", err
	}
	if odysseyGoVersion, _, err := parseOdysseyGoOutput(resp); err != nil {
		return "", err
	} else {
		return odysseyGoVersion, nil
	}
}

func (h *Node) GetOdysseyGoHealth() (bool, error) {
	// Craft and send the HTTP POST request
	requestBody := "{\"jsonrpc\":\"2.0\", \"id\":1,\"method\":\"health.health\",\"params\": {\"tags\": [\"P\"]}}"
	resp, err := h.Post("/ext/health", requestBody)
	if err != nil {
		return false, err
	}
	return parseHealthyOutput(resp)
}

func parseOdysseyGoOutput(byteValue []byte) (string, uint32, error) {
	reply := map[string]interface{}{}
	if err := json.Unmarshal(byteValue, &reply); err != nil {
		return "", 0, err
	}
	resultMap := reply["result"]
	resultJSON, err := json.Marshal(resultMap)
	if err != nil {
		return "", 0, err
	}

	nodeVersionReply := info.GetNodeVersionReply{}
	if err := json.Unmarshal(resultJSON, &nodeVersionReply); err != nil {
		return "", 0, err
	}
	return nodeVersionReply.VMVersions["platform"], uint32(nodeVersionReply.RPCProtocolVersion), nil
}

func parseHealthyOutput(byteValue []byte) (bool, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(byteValue, &result); err != nil {
		return false, err
	}
	isHealthyInterface, ok := result["result"].(map[string]interface{})
	if ok {
		isHealthy, ok := isHealthyInterface["healthy"].(bool)
		if ok {
			return isHealthy, nil
		}
	}
	return false, fmt.Errorf("unable to parse node healthy status")
}

func (h *Node) GetOdysseyGoNetworkName() (string, error) {
	if nodeConfigFileExists(*h) {
		avagoConfig, err := h.GetOdysseyGoConfigData()
		if err != nil {
			return "", err
		}
		return utils.StringValue(avagoConfig, "network-id")
	} else {
		return "", fmt.Errorf("node config file does not exist")
	}
}

func (h *Node) GetOdysseyGoConfigData() (map[string]interface{}, error) {
	// get remote node.json file
	nodeJSON, err := h.ReadFileBytes(remoteconfig.GetRemoteOdysseyNodeConfig(), constants.SSHFileOpsTimeout)
	if err != nil {
		return nil, err
	}
	var avagoConfig map[string]interface{}
	if err := json.Unmarshal(nodeJSON, &avagoConfig); err != nil {
		return nil, err
	}
	return avagoConfig, nil
}

// WaitForSSHShell waits for the SSH shell to be available on the node within the specified timeout.
func (h *Node) WaitForOdysseyGoHealth(timeout time.Duration) error {
	if h.IP == "" {
		return fmt.Errorf("node IP is empty")
	}
	start := time.Now()
	if err := h.WaitForPort(constants.OdysseygoAPIPort, timeout); err != nil {
		return err
	}

	deadline := start.Add(timeout)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout: OdysseyGo health on node %s is not available after %ds", h.IP, int(timeout.Seconds()))
		}
		if isHealthy, err := h.GetOdysseyGoHealth(); err != nil || !isHealthy {
			time.Sleep(constants.SSHSleepBetweenChecks)
			continue
		} else {
			return nil
		}
	}
}
