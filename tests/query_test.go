// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"testing"

	"github.com/llingr/llingr-nexus/nexus"
)

// Test_CommittedOffsets_Value ensures any dependent code using
// the int value is not broken by a re-ordering of iota constants
func Test_CommittedOffsets_Value(t *testing.T) {
	if committedOffsetsValue := int(nexus.CommittedOffsets); committedOffsetsValue != 1 {
		t.Error("CommittedOffsets should be set to 1")
	}
}
