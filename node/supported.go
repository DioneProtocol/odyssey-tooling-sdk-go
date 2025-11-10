// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type SupportedRole int

const (
	Validator SupportedRole = iota
	API
	Loadtest
	Monitor
)

// NewSupportedRole converts a string to a SupportedRole
func NewSupportedRole(s string) SupportedRole {
	switch s {
	case "validator":
		return Validator
	case "api":
		return API
	case "loadtest":
		return Loadtest
	case "monitor":
		return Monitor
	default:
		return Monitor
	}
}

// String returns the string representation of the SupportedRole
func (r *SupportedRole) String() string {
	switch *r {
	case Validator:
		return "validator"
	case API:
		return "api"
	case Loadtest:
		return "loadtest"
	case Monitor:
		return "monitor"
	default:
		return "unknown"
	}
}

// CheckRoles checks if the combination of roles is valid
func CheckRoles(roles []SupportedRole) error {
	if slices.Contains(roles, Validator) && slices.Contains(roles, API) {
		return fmt.Errorf("cannot have both validator and api roles")
	}
	if slices.Contains(roles, Loadtest) && len(roles) > 1 {
		return fmt.Errorf("%v role cannot be combined with other roles", Loadtest)
	}
	if slices.Contains(roles, Monitor) && len(roles) > 1 {
		return fmt.Errorf("%v role cannot be combined with other roles", Monitor)
	}
	return nil
}
