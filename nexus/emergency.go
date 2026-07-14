// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

// EmergencyShutdowner stops a compliant consumer immediately, escalating an
// irrecoverable (typically infrastructure) failure: EmergencyShutdown trips
// an engine's protective shutdown with the given reason (as error).
//
// Processing should be canceled immediately, and the registered ShutdownCallback
// notified (once); typically this won't complete processing like normal
// Consumer.Shutdown() although best-effort cleanup is left to implementors.
//
// Engines implementing this must handle this at any point in their lifecycle:
// pre-subscribe, during assign, during normal operations, and during shutdown.
//
// Calls to this function must be safe from any goroutine, and must be idempotent.
//
// Nothing is required to implement this interface. It exists for type
// assertion: an adapter or application holding a Consumer (or AdaptedConsumer)
// asserts against it to escalate without needing to implement shutdown sequencing
// itself:
//
//	if es, ok := consumer.(nexus.EmergencyShutdowner); ok {
//	    es.EmergencyShutdown(reason)
//	}
//
// Engines predating this contract need not satisfy the interface, and callers must
// fall back to their existing circuit-breaker protocol.
type EmergencyShutdowner interface {
	EmergencyShutdown(reason error)
}
