// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

// Traits bit flags capturing message processing context, actionable
// business intelligence and other indicators which can be attached
// during message processing
//
// Examples
//
//   - Set multiple traits in one call:
//     SetTraits(&traits, NewCustomer | AuditRequired)
//
//   - Alert on actionable concern:
//     const AlertCondition = DeadLetter | FraudDetected
//     if metrics.Traits & AlertCondition == AlertCondition  {
//     ... send alert ...
//     }
//
// Features:
//   - Functions to avoid bit flag mental gymnastics
//   - Silent masking protects framework-reserved flags from accidental updates
type Traits uint64

// frameworkReservedBitCount defines how many low-order bits are reserved for framework use.
// Bits 0-9 are reserved; custom traits start at bit 10.
const frameworkReservedBitCount = 10

// FrameworkReserved masks bits 0-9 for system use (0x3FF = first 10 bits set).
const FrameworkReserved Traits = (1 << frameworkReservedBitCount) - 1

// Framework-reserved trait flags (bits 0-9). Set by the framework during message
// processing; protected from modification when using SetTraits/UnsetTraits mutators.
const (
	ProcessError        Traits = 1 << iota // position 0 (1 << 0)
	ProcessPanic                           // position 1 (1 << 1)
	DeadLetter                             // position 2 (1 << 2)
	CommitBuffered                         // position 3 (1 << 3)
	Duplicate                              // position 4 (1 << 4)
	UsedOverflow                           // position 5 (1 << 5)
	Orphaned                               // position 6 (1 << 6) - completed after partition reassignment
	FirstAfterRebalance                    // position 7 (1 << 7)
)

// SetProcessError reserved trait flag
func SetProcessError(traits *Traits) {
	*traits |= ProcessError
}

// SetProcessPanic reserved trait flag
func SetProcessPanic(traits *Traits) {
	*traits |= ProcessPanic
}

// SetDeadLetter reserved trait flag
func SetDeadLetter(traits *Traits) {
	*traits |= DeadLetter
}

// SetCommitBuffered reserved trait flag
func SetCommitBuffered(traits *Traits) {
	*traits |= CommitBuffered
}

// SetDuplicate reserved trait flag
func SetDuplicate(traits *Traits) {
	*traits |= Duplicate
}

// SetUsedOverflow reserved trait flag
func SetUsedOverflow(traits *Traits) {
	*traits |= UsedOverflow
}

// SetOrphaned reserved trait flag - marks WorkItem that
// completed after rebalance / partition reassignment
func SetOrphaned(traits *Traits) {
	*traits |= Orphaned
}

// SetFirstAfterRebalance reserved trait flag - marks first message
// processed on a partition after rebalance assignment
func SetFirstAfterRebalance(traits *Traits) {
	*traits |= FirstAfterRebalance
}

// SetTraits safe, idempotent in-place custom trait(s) setter,
// won't affect reserved flags
func SetTraits(traits *Traits, traitFlags Traits) {
	*traits |= traitFlags &^ FrameworkReserved
}

// UnsetTraits safe, idempotent in-place custom trait(s) removal,
// won't affect reserved flags
func UnsetTraits(traits *Traits, traitFlags Traits) {
	*traits &= ^(traitFlags &^ FrameworkReserved)
}

// ClearTraits clears user-defined traits,
// won't affect reserved flags
func ClearTraits(traits *Traits) {
	*traits &= FrameworkReserved
}
