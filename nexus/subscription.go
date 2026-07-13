// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

import (
	"time"
)

// Poll broker client to retrieve next message or timeout,
// returns true if message received, false on timeout or error
type Poll[T any] func(timeout time.Duration) (T, bool, error)

// CommitOffsets adapter callback for the provided
// high-watermark message on each partition.
//
// Error contract: a nil error means the ENTIRE batch committed. The engine
// then advances its committed baseline for every message in the batch, so a
// partial or total broker failure MUST be reported via a non-nil error. The
// returned messages are advisory only (diagnostics); they do not stop the
// baseline from advancing.
type CommitOffsets[T any] func(hwm []*Message[T]) ([]*Message[T], error)

// Subscribe to topic. Topic name is provided at adapter construction.
type Subscribe func() error

// Unsubscribe from topic.
type Unsubscribe func() error
