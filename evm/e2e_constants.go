// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package evm

import "time"

const testnetDefaultRPC = "https://testnode.dioneprotocol.com/ext/bc/D/rpc"

// Helper function to add delay between tests to avoid rate limiting
func addTestDelay() {
	time.Sleep(5 * time.Second) // 5 second delay to avoid 429 errors when running full test suite
}
