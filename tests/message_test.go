// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"testing"

	"github.com/llingr/llingr-nexus/nexus"
)

func TestMessage_AddCustomTraits(t *testing.T) {
	t.Run("core use case - adds custom traits", func(t *testing.T) {
		message := &nexus.Message[string]{}

		// Custom traits in user space (bits 10+)
		const CustomFlag1 nexus.Traits = 1 << 10 // bit 10
		const CustomFlag2 nexus.Traits = 1 << 15 // bit 15

		message.AddCustomTraits(CustomFlag1)
		if message.Traits&CustomFlag1 == 0 {
			t.Error("CustomFlag1 should be set")
		}

		message.AddCustomTraits(CustomFlag2)
		if message.Traits&(CustomFlag1|CustomFlag2) != (CustomFlag1 | CustomFlag2) {
			t.Error("Both custom flags should be set")
		}
	})

	t.Run("framework traits are ignored", func(t *testing.T) {
		message := &nexus.Message[string]{}

		// Try to set framework-reserved traits (bits 0-9)
		maliciousTraits := nexus.ProcessError | nexus.DeadLetter | (1 << 20) // mix of framework + custom

		message.AddCustomTraits(maliciousTraits)

		// Framework bits should NOT be set
		if message.Traits&nexus.ProcessError != 0 {
			t.Error("ProcessError should not be set via AddCustomTraits")
		}
		if message.Traits&nexus.DeadLetter != 0 {
			t.Error("DeadLetter should not be set via AddCustomTraits")
		}

		// But custom bit should be set
		if message.Traits&(1<<20) == 0 {
			t.Error("Custom bit (20) should be set")
		}
	})

	t.Run("idempotent behavior", func(t *testing.T) {
		message := &nexus.Message[string]{}

		const CustomFlag nexus.Traits = 1 << 12

		// call multiple times
		message.AddCustomTraits(CustomFlag)
		firstResult := message.Traits

		message.AddCustomTraits(CustomFlag)
		secondResult := message.Traits

		message.AddCustomTraits(CustomFlag)
		thirdResult := message.Traits

		if firstResult != secondResult || secondResult != thirdResult {
			t.Error("AddCustomTraits should be idempotent - multiple calls should not change result")
		}

		if message.Traits&CustomFlag == 0 {
			t.Error("CustomFlag should be set after all calls")
		}
	})

	t.Run("preserves existing framework traits", func(t *testing.T) {
		message := &nexus.Message[string]{
			Traits: nexus.ProcessError | nexus.CommitBuffered, // pre-existing framework traits
		}

		const CustomFlag nexus.Traits = 1 << 11

		message.AddCustomTraits(CustomFlag)

		// framework traits should be preserved
		if message.Traits&nexus.ProcessError == 0 {
			t.Error("Pre-existing ProcessError should be preserved")
		}
		if message.Traits&nexus.CommitBuffered == 0 {
			t.Error("Pre-existing CommitBuffered should be preserved")
		}

		// custom trait should be added
		if message.Traits&CustomFlag == 0 {
			t.Error("CustomFlag should be added")
		}
	})

	t.Run("preserves existing custom traits", func(t *testing.T) {
		message := &nexus.Message[string]{}

		const FirstCustom nexus.Traits = 1 << 13
		const SecondCustom nexus.Traits = 1 << 25

		message.AddCustomTraits(FirstCustom)
		message.AddCustomTraits(SecondCustom)

		// Both custom traits should be present
		if message.Traits&FirstCustom == 0 {
			t.Error("FirstCustom should be preserved")
		}
		if message.Traits&SecondCustom == 0 {
			t.Error("SecondCustom should be added")
		}
	})

	t.Run("handles zero traits safely", func(t *testing.T) {
		message := &nexus.Message[string]{
			Traits: nexus.ProcessError, // existing trait
		}

		originalTraits := message.Traits

		message.AddCustomTraits(0) // no traits to add

		if message.Traits != originalTraits {
			t.Error("Adding zero traits should not change existing traits")
		}
	})

	t.Run("works with different generic types", func(t *testing.T) {
		// Test with different concrete types to ensure generics work
		messageInt := &nexus.Message[int]{
			Payload: new(int),
		}
		messageBytes := &nexus.Message[[]byte]{
			Payload: &[]byte{},
		}

		const CustomFlag nexus.Traits = 1 << 14

		messageInt.AddCustomTraits(CustomFlag)
		messageBytes.AddCustomTraits(CustomFlag)

		if messageInt.Traits&CustomFlag == 0 {
			t.Error("CustomFlag should be set on Message[int]")
		}
		if messageBytes.Traits&CustomFlag == 0 {
			t.Error("CustomFlag should be set on Message[[]byte]")
		}
	})
}
