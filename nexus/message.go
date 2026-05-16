// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

import "context"

// ProcessMessage BLOCKING call implemented by application to process Message[T].
//
// DO NOT SPAWN goroutines acting on Message[T] that outlive the function call, otherwise:
//   - can quickly overwhelm an application: framework will accelerate delivery of next Message
//   - risks message loss: framework will commit offset because return is considered 'completed'
//   - high chance of data corruption: for performance, Message pointers are recycled using sync.Pool
type ProcessMessage[T any] func(ctx context.Context, msg *Message[T]) error

// WriteDeadLetter BLOCKING call routes failed messages to an application-specific handler;
// invoked if ProcessMessage returns an error.
//
// DO NOT SPAWN goroutines acting on Message[T] that outlive the function call (see above)
// Note: failed messages will still be committed so the pipeline can keep processing.
type WriteDeadLetter[T any] func(ctx context.Context, msg *Message[T], reason error) error

// Message is broker-agnostic with cached primitives to optimize processing.
// Uses *T for sync.Pool compatibility; T can be any type including a pointer.
//
// Non-zero Traits can be provided by the host application when calling
// ProcessMessage. This is useful for creating custom metrics and injecting
// business intelligence flags.
type Message[T any] struct {
	Traits    Traits // bit flags, positions 0-9 reserved, zero is ok
	Partition int32  //
	Offset    int64  // within the partition
	Key       string // partition key
	Payload   *T     // the underlying broker-specific message type, adapter-provided
}

// AddCustomTraits which are merged with Metrics traits
// prior to metrics collection - idempotent
func (m *Message[T]) AddCustomTraits(customTraits Traits) {
	m.Traits |= customTraits &^ FrameworkReserved
}
