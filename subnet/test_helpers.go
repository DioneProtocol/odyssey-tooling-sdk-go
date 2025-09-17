// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package subnet

import "time"

// Helper function to add delay between tests to avoid rate limiting
func addTestDelay() {
	time.Sleep(2 * time.Second)
}
