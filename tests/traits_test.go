// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package tests_test

import (
	"math/bits"
	"testing"

	"github.com/llingr/llingr-nexus/nexus"
)

func TestTraitConstants(t *testing.T) {

	// Test that ProcessError is at bit position 0
	processErrorPosition := bits.TrailingZeros64(uint64(nexus.ProcessError))
	if processErrorPosition != 0 {
		t.Errorf("ProcessError should be at bit position 0, got position %d", processErrorPosition)
	}

	processPanicPosition := bits.TrailingZeros64(uint64(nexus.ProcessPanic))
	if processPanicPosition != 1 {
		t.Errorf("ProcessPanic should be at bit position 1, got position %d", processPanicPosition)
	}

	deadLetterPosition := bits.TrailingZeros64(uint64(nexus.DeadLetter))
	if deadLetterPosition != 2 {
		t.Errorf("DeadLetter should be at bit position 2, got position %d", deadLetterPosition)
	}

	commitBufferedPosition := bits.TrailingZeros64(uint64(nexus.CommitBuffered))
	if commitBufferedPosition != 3 {
		t.Errorf("DeadLetter should be at bit position 3, got position %d", commitBufferedPosition)
	}

	duplicatePosition := bits.TrailingZeros64(uint64(nexus.Duplicate))
	if duplicatePosition != 4 {
		t.Errorf("Duplicate should be at bit position 4, got position %d", duplicatePosition)
	}

	usedOverflowPosition := bits.TrailingZeros64(uint64(nexus.UsedOverflow))
	if usedOverflowPosition != 5 {
		t.Errorf("UsedOverflow should be at bit position 5, got position %d", usedOverflowPosition)
	}

	orphanedPosition := bits.TrailingZeros64(uint64(nexus.Orphaned))
	if orphanedPosition != 6 {
		t.Errorf("Orphaned should be at bit position 6, got position %d", orphanedPosition)
	}

	firstAfterRebalancePosition := bits.TrailingZeros64(uint64(nexus.FirstAfterRebalance))
	if firstAfterRebalancePosition != 7 {
		t.Errorf("FirstAfterRebalance should be at bit position 7, got position %d",
			firstAfterRebalancePosition)
	}

	// confirm reserved mask covers positions 0-9 (has 10 bits set)
	reservedBits := bits.OnesCount64(uint64(nexus.FrameworkReserved))
	if reservedBits != 10 {
		t.Errorf("FrameworkReserved should have 10 bits set, got %d", reservedBits)
	}

	// confirm reserved mask starts at bit 0
	reservedStartPosition := bits.TrailingZeros64(uint64(nexus.FrameworkReserved))
	if reservedStartPosition != 0 {
		t.Errorf("FrameworkReserved should start at bit 0, got position %d", reservedStartPosition)
	}
}

func TestSetTraitReservedPositions(t *testing.T) {
	var traits nexus.Traits

	// Test all reserved positions (0-9) should be silently masked
	for i := 0; i < 10; i++ {
		traits = 0 // Reset for each test
		traitFlags := nexus.Traits(1 << i)
		nexus.SetTraits(&traits, traitFlags)

		// Reserved flags should be silently ignored
		if traits != 0 {
			t.Errorf("SetTraits should silently mask reserved position %d (flag %d), got traits %d", i, traitFlags, traits)
		}
	}
}

func TestSetTraitValidPositions(t *testing.T) {
	var traits nexus.Traits

	// Test all valid positions (10-63) should be accepted
	for i := 10; i < 64; i++ {
		traits = 0 // Reset for each test
		traitFlags := nexus.Traits(1 << i)
		nexus.SetTraits(&traits, traitFlags)
		if traits != traitFlags {
			t.Errorf("SetTraits should set traits to %d, got %d", traitFlags, traits)
		}
	}
}

func TestSetTraitMultipleBits(t *testing.T) {
	var traits nexus.Traits

	// Test that multiple bits in reserved range are silently masked
	multipleReserved := nexus.Traits(3) // positions 0 and 1
	nexus.SetTraits(&traits, multipleReserved)
	if traits != 0 {
		t.Error("SetTraits should silently mask reserved bits")
	}

	// Test that multiple bits spanning reserved and unreserved keep only unreserved
	traits = 0
	spanningBits := nexus.Traits((1 << 5) | (1 << 15)) // position 5 (reserved) + 15 (unreserved)
	nexus.SetTraits(&traits, spanningBits)
	expected := nexus.Traits(1 << 15) // Only unreserved bit should remain
	if traits != expected {
		t.Errorf("SetTraits should mask reserved bits, expected %d, got %d", expected, traits)
	}

	// Test that multiple bits in unreserved range work correctly
	traits = 0
	multipleUnreserved := nexus.Traits((1 << 15) | (1 << 20))
	nexus.SetTraits(&traits, multipleUnreserved)
	if traits != multipleUnreserved {
		t.Errorf("SetTraits should set multiple unreserved bits, expected %d, got %d", multipleUnreserved, traits)
	}
}

func TestSetTraits_ZeroNoop(t *testing.T) {
	var traits nexus.Traits
	nexus.SetTraits(&traits, 0)
	if traits != 0 {
		t.Error("no-op should not have adjusted any traits")
	}

	// set a reserved trait, then verify no-op doesn't clear it
	nexus.SetDeadLetter(&traits)
	originalTraits := traits

	nexus.SetTraits(&traits, 0)
	if traits != originalTraits {
		t.Error("no-op should preserve existing traits")
	}
	if traits&nexus.DeadLetter == 0 {
		t.Error("no-op should not clear DeadLetter flag")
	}
}

func TestSetDeadLetter(t *testing.T) {
	var traits nexus.Traits
	nexus.SetDeadLetter(&traits)

	if traits&nexus.DeadLetter == 0 {
		t.Error("SetDeadLetter should set DeadLetter flag")
	}
	if traits != nexus.DeadLetter {
		t.Errorf("SetDeadLetter should only set DeadLetter flag, got %d", traits)
	}
}

func TestSetDuplicate(t *testing.T) {
	var traits nexus.Traits
	nexus.SetDuplicate(&traits)

	if traits&nexus.Duplicate == 0 {
		t.Error("SetDuplicate should set Duplicate flag")
	}
	if traits != nexus.Duplicate {
		t.Errorf("SetDuplicate should only set Duplicate flag, got %d", traits)
	}
}

func TestSetUsedOverflow(t *testing.T) {
	var traits nexus.Traits
	nexus.SetUsedOverflow(&traits)

	if traits&nexus.UsedOverflow != nexus.UsedOverflow {
		t.Error("SetUsedOverflow should set UsedOverflow flag")
	}
	if traits != nexus.UsedOverflow {
		t.Errorf("SetUsedOverflow should only set UsedOverflow flag, got %d", traits)
	}
}

func TestSetOrphaned(t *testing.T) {
	var traits nexus.Traits
	nexus.SetOrphaned(&traits)

	if traits&nexus.Orphaned != nexus.Orphaned {
		t.Error("SetOrphaned should set Orphaned flag")
	}
	if traits != nexus.Orphaned {
		t.Errorf("SetOrphaned should only set Orphaned flag, got %d", traits)
	}
}

func TestSetFirstAfterRebalance(t *testing.T) {
	var traits nexus.Traits
	nexus.SetFirstAfterRebalance(&traits)

	if traits&nexus.FirstAfterRebalance != nexus.FirstAfterRebalance {
		t.Error("SetFirstAfterRebalance should set FirstAfterRebalance flag")
	}
	if traits != nexus.FirstAfterRebalance {
		t.Errorf("SetFirstAfterRebalance should only set FirstAfterRebalance flag, got %d", traits)
	}
}

func TestClearTraitsPreservesReserved(t *testing.T) {
	var traits nexus.Traits

	// reserved traits
	nexus.SetDeadLetter(&traits)
	nexus.SetDuplicate(&traits)

	// user-defined traits
	nexus.SetTraits(&traits, 1<<15)
	nexus.SetTraits(&traits, 1<<25)

	// confirm both reserved and user traits are set
	expectedBefore := nexus.DeadLetter | nexus.Duplicate | (1 << 15) | (1 << 25)
	if traits != expectedBefore {
		t.Errorf("Expected traits before clear to be %d, got %d", expectedBefore, traits)
	}

	nexus.ClearTraits(&traits)

	// confirm reserved traits are preserved
	if traits&nexus.DeadLetter == 0 {
		t.Error("ClearTraits should preserve DeadLetter flag")
	}
	if traits&nexus.Duplicate == 0 {
		t.Error("ClearTraits should preserve Duplicate flag")
	}

	// confirm user traits are cleared
	if traits&(1<<15) != 0 {
		t.Error("ClearTraits should clear user-defined trait at position 15")
	}
	if traits&(1<<25) != 0 {
		t.Error("ClearTraits should clear user-defined trait at position 25")
	}

	// confirm only reserved traits remain
	expectedAfter := nexus.DeadLetter | nexus.Duplicate
	if traits != expectedAfter {
		t.Errorf("Expected traits after clear to be %d, got %d", expectedAfter, traits)
	}
}

func TestClearTraitsWithNoReservedTraits(t *testing.T) {
	var traits nexus.Traits

	// set only user-defined traits (no reserved traits)
	nexus.SetTraits(&traits, 1<<15)
	nexus.SetTraits(&traits, 1<<30)

	nexus.ClearTraits(&traits)

	// should be completely empty (no reserved traits to preserve)
	if traits != 0 {
		t.Errorf("ClearTraits should clear all traits when no reserved traits set, got %d", traits)
	}
}

func TestClearTraitsWithOnlyReservedTraits(t *testing.T) {
	var traits nexus.Traits

	// set only reserved traits
	nexus.SetDeadLetter(&traits)
	nexus.SetDuplicate(&traits)
	originalTraits := traits

	nexus.ClearTraits(&traits)

	// should preserve all reserved traits
	if traits != originalTraits {
		t.Errorf("ClearTraits should preserve reserved traits unchanged, expected %d, got %d", originalTraits, traits)
	}
}

func TestSetTraitIdempotent(t *testing.T) {
	var traits nexus.Traits
	flag := nexus.Traits(1 << 15)

	// set twice
	nexus.SetTraits(&traits, flag)
	firstResult := traits
	nexus.SetTraits(&traits, flag)
	secondResult := traits

	if firstResult != secondResult {
		t.Error("SetTraits should be idempotent")
	}
	if traits != flag {
		t.Errorf("Expected traits to be %d, got %d", flag, traits)
	}
}

func TestReservedTraitsIdempotent(t *testing.T) {
	var traits nexus.Traits

	traits = 0
	nexus.SetProcessError(&traits)
	nexus.SetProcessError(&traits)
	if traits != nexus.ProcessError {
		t.Error("SetProcessError should be idempotent")
	}

	traits = 0
	nexus.SetProcessPanic(&traits)
	nexus.SetProcessPanic(&traits)
	if traits != nexus.ProcessPanic {
		t.Error("SetProcessPanic should be idempotent")
	}

	traits = 0
	nexus.SetDeadLetter(&traits)
	nexus.SetDeadLetter(&traits)
	if traits != nexus.DeadLetter {
		t.Error("SetDeadLetter should be idempotent")
	}

	traits = 0
	nexus.SetCommitBuffered(&traits)
	nexus.SetCommitBuffered(&traits)
	if traits != nexus.CommitBuffered {
		t.Error("SetCommitBuffered should be idempotent")
	}

	traits = 0
	nexus.SetDuplicate(&traits)
	nexus.SetDuplicate(&traits)
	if traits != nexus.Duplicate {
		t.Error("SetDuplicate should be idempotent")
	}

	traits = 0
	nexus.SetUsedOverflow(&traits)
	nexus.SetUsedOverflow(&traits)
	if traits != nexus.UsedOverflow {
		t.Error("SetUsedOverflow should be idempotent")
	}

	traits = 0
	nexus.SetOrphaned(&traits)
	nexus.SetOrphaned(&traits)
	if traits != nexus.Orphaned {
		t.Error("SetOrphaned should be idempotent")
	}

	traits = 0
	nexus.SetFirstAfterRebalance(&traits)
	nexus.SetFirstAfterRebalance(&traits)
	if traits != nexus.FirstAfterRebalance {
		t.Error("SetFirstAfterRebalance should be idempotent")
	}
}

func TestCombinedReservedTraits(t *testing.T) {
	var traits nexus.Traits

	nexus.SetProcessError(&traits)
	nexus.SetProcessPanic(&traits)
	nexus.SetDeadLetter(&traits)
	nexus.SetCommitBuffered(&traits)
	nexus.SetDuplicate(&traits)
	nexus.SetUsedOverflow(&traits)
	nexus.SetOrphaned(&traits)
	nexus.SetFirstAfterRebalance(&traits)

	expected := nexus.ProcessError | nexus.ProcessPanic | nexus.DeadLetter |
		nexus.CommitBuffered | nexus.Duplicate | nexus.UsedOverflow | nexus.Orphaned |
		nexus.FirstAfterRebalance
	if traits != expected {
		t.Errorf("Combined reserved traits should be %d, got %d", expected, traits)
	}

	// confirm all flags are set
	if traits&nexus.ProcessError == 0 {
		t.Error("ProcessError flag should be set")
	}
	if traits&nexus.ProcessPanic == 0 {
		t.Error("ProcessPanic flag should be set")
	}
	if traits&nexus.DeadLetter == 0 {
		t.Error("DeadLetter flag should be set")
	}
	if traits&nexus.CommitBuffered == 0 {
		t.Error("CommitBuffered flag should be set")
	}
	if traits&nexus.Duplicate == 0 {
		t.Error("Duplicate flag should be set")
	}
	if traits&nexus.UsedOverflow == 0 {
		t.Error("UsedOverflow flag should be set")
	}
	if traits&nexus.Orphaned == 0 {
		t.Error("Orphaned flag should be set")
	}
	if traits&nexus.FirstAfterRebalance == 0 {
		t.Error("FirstAfterRebalance flag should be set")
	}
}

func TestTraitsWithUnreservedFlags(t *testing.T) {
	var traits nexus.Traits

	// some reserved traits
	nexus.SetDeadLetter(&traits)
	nexus.SetDuplicate(&traits)

	// some unreserved traits
	flag1 := nexus.Traits(1 << 15)
	flag2 := nexus.Traits(1 << 25)
	nexus.SetTraits(&traits, flag1)
	nexus.SetTraits(&traits, flag2)

	expected := nexus.DeadLetter | nexus.Duplicate | flag1 | flag2
	if traits != expected {
		t.Errorf("Combined traits should be %d, got %d", expected, traits)
	}

	// confirm all flags are set
	if traits&nexus.DeadLetter == 0 {
		t.Error("DeadLetter flag should be set")
	}
	if traits&nexus.Duplicate == 0 {
		t.Error("Duplicate flag should be set")
	}
	if traits&flag1 == 0 {
		t.Error("Unreserved flag1 should be set")
	}
	if traits&flag2 == 0 {
		t.Error("Unreserved flag2 should be set")
	}
}

func TestBoundaryConditions(t *testing.T) {
	var traits nexus.Traits

	// position 10 (first valid unreserved position) should be set
	pos10 := nexus.Traits(1 << 10)
	nexus.SetTraits(&traits, pos10)
	if traits&pos10 == 0 {
		t.Error("position 10 (first unreserved) should be set")
	}

	// position 63 (highest valid position) should be set
	traits = 0
	pos63 := nexus.Traits(1 << 63)
	nexus.SetTraits(&traits, pos63)
	if traits&pos63 == 0 {
		t.Error("position 63 (highest valid) should be set")
	}

	// position 9 (last reserved position) should NOT be set
	traits = 0
	pos9 := nexus.Traits(1 << 9)
	nexus.SetTraits(&traits, pos9)
	if traits&pos9 != 0 {
		t.Error("position 9 (reserved) should not be set via SetTraits")
	}
}

func TestSetTraitsMultipleBitsUnreserved(t *testing.T) {
	var traits nexus.Traits

	// single unreserved bit
	singleBit := nexus.Traits(1 << 15)
	nexus.SetTraits(&traits, singleBit)
	if traits != singleBit {
		t.Error("Single bit should be set correctly")
	}

	// multiple unreserved bits
	traits = 0
	multipleBits := nexus.Traits((1 << 15) | (1 << 16) | (1 << 20))
	nexus.SetTraits(&traits, multipleBits)
	if traits != multipleBits {
		t.Errorf("Multiple unreserved bits should be set correctly, expected %d, got %d", multipleBits, traits)
	}

	// confirm individual bits set
	if traits&(1<<15) == 0 {
		t.Error("Bit 15 should be set")
	}
	if traits&(1<<16) == 0 {
		t.Error("Bit 16 should be set")
	}
	if traits&(1<<20) == 0 {
		t.Error("Bit 20 should be set")
	}
}

func TestUnsetTraitsReservedPositions(t *testing.T) {
	var traits nexus.Traits

	// set unreserved trait flag
	nexus.SetTraits(&traits, 1<<15)
	originalTraits := traits

	// all reserved positions (0-9) should be silently masked
	for i := 0; i < 10; i++ {
		traitFlag := nexus.Traits(1 << i)
		nexus.UnsetTraits(&traits, traitFlag)

		if traits != originalTraits {
			t.Errorf("UnsetTraits should silently mask reserved position %d, traits changed from %d to %d", i, originalTraits, traits)
		}
	}
}

func TestUnsetTraitsValidPositions(t *testing.T) {
	var traits nexus.Traits

	// all unreserved positions (10-63) should be accepted
	for i := 10; i < 64; i++ {
		// First set the trait
		traitFlag := nexus.Traits(1 << i)
		nexus.SetTraits(&traits, traitFlag)

		// should leave nothing set
		nexus.UnsetTraits(&traits, traitFlag)
		if traits != 0 {
			t.Errorf("UnsetTraits should clear traits to 0, got %d", traits)
		}
	}
}

func TestUnsetTraitsMultipleBits(t *testing.T) {
	var traits nexus.Traits

	// set dead letter bit first
	nexus.SetDeadLetter(&traits)

	// set multiple unreserved bits
	flag1 := nexus.Traits(1 << 15)
	flag2 := nexus.Traits(1 << 20)
	flag3 := nexus.Traits(1 << 25)
	combinedFlags := flag1 | flag2 | flag3

	nexus.SetTraits(&traits, combinedFlags)

	// verify all bits are set as expected
	expected := nexus.DeadLetter | combinedFlags
	if traits != expected {
		t.Fatalf("Expected traits to be %d, got %d", expected, traits)
	}

	// attempt to unset multiple bits including dead letter bit
	toUnset := nexus.DeadLetter | flag1 | flag2 // includes reserved bit
	nexus.UnsetTraits(&traits, toUnset)

	// verify dead letter bit remains but unreserved bits are cleared
	expectedAfterUnset := nexus.DeadLetter | flag3 // only flag3 should remain
	if traits != expectedAfterUnset {
		t.Fatalf("Expected traits to be %d after unset, got %d", expectedAfterUnset, traits)
	}

	// verify dead letter bit specifically is still set
	if traits&nexus.DeadLetter == 0 {
		t.Error("Dead letter bit should remain set after unset operation")
	}
}

func TestUnsetTraitsIdempotent(t *testing.T) {
	var traits nexus.Traits
	flag := nexus.Traits(1 << 15)

	nexus.SetTraits(&traits, flag)

	// unset twice
	nexus.UnsetTraits(&traits, flag)
	firstResult := traits
	nexus.UnsetTraits(&traits, flag)
	secondResult := traits

	if firstResult != secondResult {
		t.Error("UnsetTraits should be idempotent")
	}
	if traits != 0 {
		t.Errorf("Expected traits to be 0, got %d", traits)
	}
}

func TestUnsetTraitsReservedHasNoEffect(t *testing.T) {
	var traits nexus.Traits

	// set reserved traits using proper functions
	nexus.SetDeadLetter(&traits)
	nexus.SetDuplicate(&traits)

	// add some custom traits too
	nexus.SetTraits(&traits, 1<<15)
	nexus.SetTraits(&traits, 1<<25)

	originalTraits := traits

	// try to unset reserved traits - should be no-op
	nexus.UnsetTraits(&traits, nexus.DeadLetter)
	nexus.UnsetTraits(&traits, nexus.Duplicate)
	nexus.UnsetTraits(&traits, 1<<9) // position 9 is also reserved

	// traits should be unchanged
	if traits != originalTraits {
		t.Errorf("UnsetTraits on reserved flags should have no effect")
	}

	// verify reserved traits are still set
	if traits&nexus.DeadLetter == 0 {
		t.Error("DeadLetter should still be set after UnsetTraits attempt")
	}
	if traits&nexus.Duplicate == 0 {
		t.Error("Duplicate should still be set after UnsetTraits attempt")
	}
}

func TestUnsetTraits_ZeroNoop(t *testing.T) {
	var traits nexus.Traits
	nexus.UnsetTraits(&traits, 0)
	if traits != 0 {
		t.Error("no-op should not have adjusted any traits")
	}

	// set some traits, then verify zero no-op doesn't clear them
	nexus.SetDeadLetter(&traits)
	nexus.SetTraits(&traits, 1<<15)
	originalTraits := traits

	nexus.UnsetTraits(&traits, 0)
	if traits != originalTraits {
		t.Error("no-op should preserve existing traits")
	}
}

func TestUnsetTraitsMixedReservedUnreserved(t *testing.T) {
	var traits nexus.Traits

	// set up initial state with both reserved and user traits
	nexus.SetDeadLetter(&traits)
	nexus.SetDuplicate(&traits)
	nexus.SetTraits(&traits, 1<<15)
	nexus.SetTraits(&traits, 1<<20)

	originalReserved := traits & nexus.FrameworkReserved

	// try to unset mix of reserved + unreserved flags
	mixedFlags := nexus.DeadLetter | (1 << 15) | (1 << 9) // reserved + unreserved
	nexus.UnsetTraits(&traits, mixedFlags)

	// reserved traits should be unchanged
	if traits&nexus.FrameworkReserved != originalReserved {
		t.Error("reserved traits should be unchanged when mixed with unreserved")
	}

	// unreserved flag should be cleared
	if traits&(1<<15) != 0 {
		t.Error("unreserved flag should be cleared from mixed operation")
	}

	// other unreserved flags should remain
	if traits&(1<<20) == 0 {
		t.Error("other unreserved flags should remain untouched")
	}
}

func TestClearTraitsIdempotent(t *testing.T) {
	var traits nexus.Traits

	// set up initial state
	nexus.SetDeadLetter(&traits)
	nexus.SetDuplicate(&traits)
	nexus.SetTraits(&traits, 1<<15)
	nexus.SetTraits(&traits, 1<<25)

	// clear multiple times
	nexus.ClearTraits(&traits)
	firstResult := traits
	nexus.ClearTraits(&traits)
	secondResult := traits
	nexus.ClearTraits(&traits)
	thirdResult := traits

	if firstResult != secondResult || secondResult != thirdResult {
		t.Error("ClearTraits should be idempotent")
	}

	// should preserve only reserved traits
	expected := nexus.DeadLetter | nexus.Duplicate
	if traits != expected {
		t.Errorf("expected %d after multiple clears, got %d", expected, traits)
	}
}

func TestUnsetTraitsBoundaryConditions(t *testing.T) {
	var traits nexus.Traits

	// test position 10 (first valid unreserved position)
	nexus.SetTraits(&traits, 1<<10)
	nexus.UnsetTraits(&traits, 1<<10)
	if traits != 0 {
		t.Error("position 10 should be unset correctly")
	}

	// test position 63 (highest valid position)
	nexus.SetTraits(&traits, 1<<63)
	nexus.UnsetTraits(&traits, 1<<63)
	if traits != 0 {
		t.Error("position 63 should be unset correctly")
	}

	// test position 9 (last reserved position) - should have no effect
	nexus.SetDeadLetter(&traits) // set some reserved trait
	originalTraits := traits
	nexus.UnsetTraits(&traits, 1<<9)
	if traits != originalTraits {
		t.Error("position 9 should be silently masked (reserved)")
	}
}

// should observe ~0.25 ns/op
func BenchmarkSetTrait(b *testing.B) {
	var traits nexus.Traits
	traitFlags := nexus.Traits(1 << 15)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nexus.SetTraits(&traits, traitFlags)
	}
}

// should observe ~0.25 ns/op
func BenchmarkSetReservedTraits(b *testing.B) {
	var traits nexus.Traits

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nexus.SetDeadLetter(&traits)
		nexus.SetDuplicate(&traits)
	}
}
