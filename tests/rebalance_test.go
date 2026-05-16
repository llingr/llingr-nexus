// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"testing"

	"github.com/llingr/llingr-nexus/nexus"
)

// Test_Assign_Value ensures any dependent code using the int
// value is not broken by a re-ordering of iota constants
func Test_Assign_Value(t *testing.T) {
	if assignValue := int(nexus.Assign); assignValue != 1 {
		t.Error("nexus.Assign should remain set to 1")
	}
}

// Test_Revoke_Value ensures any dependent code using the int
// value is not broken by a re-ordering of iota constants
func Test_Revoke_Value(t *testing.T) {
	if revokeValue := int(nexus.Revoke); revokeValue != 2 {
		t.Error("nexus.Revoke should remain set to 2")
	}
}
