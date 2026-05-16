// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

import "testing"

// TestFrameworkReservedConstants provides mutation testing anchor points.
// These explicit value checks ensure mutations to frameworkReservedBitCount
// are detected by gremlins (which runs tests in the same directory as mutated code).
func TestFrameworkReservedConstants(t *testing.T) {
	if frameworkReservedBitCount != 10 {
		t.Errorf("frameworkReservedBitCount must be 10, got %d", frameworkReservedBitCount)
	}

	// 0x3FF = (1 << 10) - 1 = 1023 = bits 0-9 set
	if FrameworkReserved != 0x3FF {
		t.Errorf("FrameworkReserved must be 0x3FF (1023), got 0x%X (%d)",
			FrameworkReserved, FrameworkReserved)
	}
}
