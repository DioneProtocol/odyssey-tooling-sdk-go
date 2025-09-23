// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package evm

import (
	"os"
	"time"
)

// getTestnetRPC returns the testnet RPC endpoint, preferring local node if available
func getTestnetRPC() string {
	if os.Getenv("LOCAL_NODE") == "true" {
		return "http://127.0.0.1:9650/ext/bc/D/rpc"
	}
	return "https://testnode.dioneprotocol.com/ext/bc/D/rpc"
}

var testnetDefaultRPC = getTestnetRPC()

// Helper function to add delay between tests to avoid rate limiting
func addTestDelay() {
	if os.Getenv("LOCAL_NODE") == "true" {
		time.Sleep(100 * time.Millisecond) // Minimal delay for local node
	} else {
		time.Sleep(5 * time.Second) // 5 second delay to avoid 429 errors when using external endpoints
	}
}
