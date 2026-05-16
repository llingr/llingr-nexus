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
// high-watermark message on each partition
type CommitOffsets[T any] func(hwm []*Message[T]) ([]*Message[T], error)

// Subscribe to topic. Topic name is provided at adapter construction.
type Subscribe func() error

// Unsubscribe from topic.
type Unsubscribe func() error
