// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

import "context"

// ExtractEnvelope implemented in adapter to map broker-specific
// payload into core addressing information (Key, Partition, Offset).
type ExtractEnvelope[T any] func(payload T) Envelope

// Envelope with message address information: Key, Partition, Offset.
//
// ProcessMessage receives Ctx as the adapter set it, and a nil Ctx is not filled
// in. Ctx may hold values for traces and spans, but a cancellation or deadline
// added to it can interfere with the framework's circuit breakers.
type Envelope struct {
	Partition int32
	Offset    int64           // within the partition
	Key       string          // partition key
	Ctx       context.Context // AdaptedConsumer.Context(), or derived from it
}
