// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"testing"

	"github.com/llingr/llingr-nexus/nexus"
)

func TestMetrics_AddCustomTraits(t *testing.T) {
	t.Run("core use case - adds custom traits", func(t *testing.T) {
		metrics := &nexus.Metrics{}

		// Custom traits in user space (bits 10+)
		const CustomFlag1 nexus.Traits = 1 << 10 // bit 10
		const CustomFlag2 nexus.Traits = 1 << 15 // bit 15

		metrics.AddCustomTraits(CustomFlag1)
		if metrics.Traits&CustomFlag1 == 0 {
			t.Error("CustomFlag1 should be set")
		}

		metrics.AddCustomTraits(CustomFlag2)
		if metrics.Traits&(CustomFlag1|CustomFlag2) != (CustomFlag1 | CustomFlag2) {
			t.Error("Both custom flags should be set")
		}
	})

	t.Run("framework traits are ignored", func(t *testing.T) {
		metrics := &nexus.Metrics{}

		// Try to set framework-reserved traits (bits 0-9)
		maliciousTraits := nexus.ProcessError | nexus.DeadLetter | (1 << 20) // mix of framework + custom

		metrics.AddCustomTraits(maliciousTraits)

		// Framework bits should NOT be set
		if metrics.Traits&nexus.ProcessError != 0 {
			t.Error("ProcessError should not be set via AddCustomTraits")
		}
		if metrics.Traits&nexus.DeadLetter != 0 {
			t.Error("DeadLetter should not be set via AddCustomTraits")
		}

		// But custom bit should be set
		if metrics.Traits&(1<<20) == 0 {
			t.Error("Custom bit (20) should be set")
		}
	})

	t.Run("idempotent behavior", func(t *testing.T) {
		metrics := &nexus.Metrics{}

		const CustomFlag nexus.Traits = 1 << 12

		// call multiple times
		metrics.AddCustomTraits(CustomFlag)
		firstResult := metrics.Traits

		metrics.AddCustomTraits(CustomFlag)
		secondResult := metrics.Traits

		metrics.AddCustomTraits(CustomFlag)
		thirdResult := metrics.Traits

		if firstResult != secondResult || secondResult != thirdResult {
			t.Error("AddCustomTraits should be idempotent - multiple calls should not change result")
		}

		if metrics.Traits&CustomFlag == 0 {
			t.Error("CustomFlag should be set after all calls")
		}
	})

	t.Run("preserves existing framework traits", func(t *testing.T) {
		metrics := &nexus.Metrics{
			Traits: nexus.ProcessError | nexus.CommitBuffered, // pre-existing framework traits
		}

		const CustomFlag nexus.Traits = 1 << 11

		metrics.AddCustomTraits(CustomFlag)

		// framework traits should be preserved
		if metrics.Traits&nexus.ProcessError == 0 {
			t.Error("Pre-existing ProcessError should be preserved")
		}
		if metrics.Traits&nexus.CommitBuffered == 0 {
			t.Error("Pre-existing CommitBuffered should be preserved")
		}

		// custom trait should be added
		if metrics.Traits&CustomFlag == 0 {
			t.Error("CustomFlag should be added")
		}
	})

	t.Run("preserves existing custom traits", func(t *testing.T) {
		metrics := &nexus.Metrics{}

		const FirstCustom nexus.Traits = 1 << 13
		const SecondCustom nexus.Traits = 1 << 25

		metrics.AddCustomTraits(FirstCustom)
		metrics.AddCustomTraits(SecondCustom)

		// Both custom traits should be present
		if metrics.Traits&FirstCustom == 0 {
			t.Error("FirstCustom should be preserved")
		}
		if metrics.Traits&SecondCustom == 0 {
			t.Error("SecondCustom should be added")
		}
	})

	t.Run("handles zero traits safely", func(t *testing.T) {
		metrics := &nexus.Metrics{
			Traits: nexus.ProcessError, // existing trait
		}

		originalTraits := metrics.Traits

		metrics.AddCustomTraits(0) // no traits to add

		if metrics.Traits != originalTraits {
			t.Error("Adding zero traits should not change existing traits")
		}
	})
}
