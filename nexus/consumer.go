// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

import (
	"context"
)

// Consumer returned to the host application, providing
// high-level lifecycle control to Subscribe and Shutdown
type Consumer[T any] interface {
	// Subscribe and start consuming messages
	Subscribe() error

	// Shutdown stops polling, finishes outstanding
	// message processing, commits, and unsubscribes.
	Shutdown() error
}

// ShutdownCallback is invoked when the consumer exits.
//   - reason is nil for graceful shutdown, non-nil for emergency (e.g., circuit breaker)
type ShutdownCallback func(ctx context.Context, reason error)
