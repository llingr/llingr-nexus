// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

// Port-Binding Builder, where an adapter is both the factory and the dependency.
//
// Traditional DI injects dependencies into a component. Here, the adapter injects itself
// into a builder that the application provides, enabling construction across module
// boundaries that don't share imports.
//
// Application -> creates Builder
//             -> passes to Adapter
//             -> Adapter calls Build(self)
// Application <- fully-wired Consumer returned

package nexus

import (
	"context"
	"time"
)

// ConsumerBuilder carries host application dependencies across module boundaries.
// The adapter calls Build() to inject itself as BrokerPort, completing construction.
//
// This implements a Port-Binding Builder pattern:
//   - Builder travels to adapter, returns fully wired consumer
//   - Build(BrokerPort[T]) is the binding site - compiler-verified type safety
//   - TopicName() exposes the configured topic to adapters
//   - Inspired by Hexagonal Architecture (Alistair Cockburn, 2005)
type ConsumerBuilder[T any] interface {
	// Build wires the adapter as BrokerPort and returns the consumer.
	Build(BrokerPort[T]) AdaptedConsumer[T]

	// TopicName returns the configured topic/stream name.
	// Adapters call this in CreateConsumer() to get the topic for subscription.
	TopicName() string
}

// AdaptedConsumer extends Consumer with adapter-internal methods.
// The adapter receives this from the builder and stores it internally,
// but returns only Consumer to the host application.
type AdaptedConsumer[T any] interface {
	// Consumer is embedded, so the adapter
	// can call Subscribe and Shutdown.
	Consumer[T]

	// TriggerRebalance initiates partition
	// reassignment processing triggered from the broker
	TriggerRebalance(RebalanceType, []RebalanceInfo) error

	// Context for control-plane logging and consumer lifecycle
	Context() context.Context

	// Logger for control-plane
	Logger() Logger
}

// BrokerPort is the adapter interface, connecting the consumer to the
// broker via an appropriate adapter for Kafka, RedPanda, Pulsar...
//
// Inspired by: Hexagonal Architecture (Alistair Cockburn, 2005)
type BrokerPort[T any] interface {
	// Subscribe to topic (topic name provided at adapter construction).
	Subscribe() error

	// Unsubscribe from topic.
	Unsubscribe() error

	// Poll broker for the next message, or timeout.
	// Returns:
	//  - SUCCESS: (message, true, nil)
	//  - TIMEOUT: (zero, false, nil)
	//  - ERROR: (zero, false, error)
	Poll(timeout time.Duration) (T, bool, error)

	// ExtractEnvelope maps broker-specific payload to partition/offset addressing info.
	ExtractEnvelope(T) Envelope

	// CommitOffsets commits the high-watermark offsets.
	CommitOffsets(messages []*Message[T]) ([]*Message[T], error)

	// AckRebalance acknowledges rebalance can complete.
	// RebalanceInfo.Meta is adapter-specific; adapter casts internally.
	AckRebalance(RebalanceType, []RebalanceInfo) error

	// BrokerQuery for committed offsets and other broker queries.
	// Request/Response Data fields are adapter-specific.
	BrokerQuery(QueryRequest) (QueryResponse, error)

	// ConsumerGroup configured on this adapter. May return empty string
	// when using - for example - manual partition assignment).
	ConsumerGroup() string
}
