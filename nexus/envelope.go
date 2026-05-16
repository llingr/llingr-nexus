// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

import "context"

// ExtractEnvelope implemented in adapter to map broker-specific
// payload into core addressing information (Key, Partition, Offset).
type ExtractEnvelope[T any] func(payload T) Envelope

// Envelope with message address information: Key, Partition, Offset.
//
// While Ctx may be needed for cross-cutting traces/spans and other
// implementation-specific concerns, avoid making this cancellable since
// this can interfere with internal framework circuit-breakers.
type Envelope struct {
	Partition int32
	Offset    int64           // within the partition
	Key       string          // partition key
	Ctx       context.Context // framework-provided if nil
}
